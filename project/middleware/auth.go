package middleware

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/session"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const (
	userIDKey      contextKey = "UserID"
	phoneNumberKey contextKey = "PhoneNumber"
)

type ServInterceptor struct{}

// Initialize Clerk SDK once when the interceptor is created
func NewServInterceptor() *ServInterceptor {
	return &ServInterceptor{}
}

func (middleware *ServInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx, err = middleware.authorize(ctx)
		if err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func (middleware *ServInterceptor) authorize(ctx context.Context) (context.Context, error) {
	clerk.SetKey(os.Getenv("CLERK_SECRET_KEY"))

	// Step 1: Extract authorization header
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, errors.New("no metadata found in context")
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing authorization header")
	}

	// Step 2: Extract JWT token
	tokenString := strings.TrimPrefix(authHeaders[0], "Bearer ")
	if tokenString == "" {
		return nil, status.Error(codes.Unauthenticated, "empty authorization header")
	}

	// Step 3: Verify JWT token
	verifyParams := &jwt.VerifyParams{Token: tokenString}
	claims, err := jwt.Verify(ctx, verifyParams)
	if err != nil {
		log.Printf("JWT verification failed: %v", err)
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	// Step 4: Extract session ID from claims
	sessionID := claims.SessionID
	if sessionID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing session ID in token")
	}

	// Step 5: Get session details
	sess, err := session.Get(ctx, sessionID)
	if err != nil {
		log.Printf("Session verification failed: %v", err)
		return nil, status.Error(codes.Unauthenticated, "invalid session")
	}

	// Step 6: Validate session status
	if sess.Status != "active" {
		return nil, status.Error(codes.Unauthenticated, "inactive session")
	}

	// Step 7: Get user details
	user, err := user.Get(ctx, sess.UserID)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid user")
	}

	// Step 8: Extract phone number
	phoneNumber := ""
	if len(user.PhoneNumbers) > 0 {
		phoneNumber = user.PhoneNumbers[0].PhoneNumber
	}

	// Step 9: Add claims to context
	ctx = context.WithValue(ctx, userIDKey, sess.UserID)
	ctx = context.WithValue(ctx, phoneNumberKey, phoneNumber)

	return ctx, nil
}
