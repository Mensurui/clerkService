package middleware

import (
	"context"
	"errors"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2/session"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/golang-jwt/jwt/v5"
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

type UserClaims struct {
	jwt.RegisteredClaims
	PhoneNumber string `json:"phone_number"`
	UserID      string `json:"id"`
}

type ServInterceptor struct{}

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
	//step 1 get the token from the request/metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, errors.New("no metadata found in context")
	}
	authHeaders := md.Get("authorization")[0]
	if len(authHeaders) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing authorization header")
	}

	tokenString := strings.TrimPrefix(authHeaders, "Bearer ")

	//step 2 verify session using clerk sdk
	sess, err := session.Get(ctx, tokenString)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid session")
	}

	//step 3 check if session is active
	if sess.Status != "active" {
		return nil, status.Error(codes.Unauthenticated, "inactive session")
	}

	//step 4 get user details
	user, err := user.Get(ctx, sess.UserID)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid user")
	}

	phoneNumber := ""
	if len(user.PhoneNumbers) > 0 {
		phoneNumber = user.PhoneNumbers[0].PhoneNumber
	}

	//step 6 add claims to context
	ctx = context.WithValue(ctx, userIDKey, sess.UserID)
	ctx = context.WithValue(ctx, phoneNumberKey, phoneNumber)
	return ctx, nil
}
