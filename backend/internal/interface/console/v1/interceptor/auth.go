package interceptor

import (
	"context"
	"strings"

	"connectrpc.com/connect"
	"github.com/shibayama-club/keyhub/internal/usecase/console/iface"
)

type contextKey string

const ConsoleSessionKey contextKey = "console_session"

type authInterceptor struct {
	useCase iface.IUseCase
}

func NewAuthInterceptor(useCase iface.IUseCase) connect.Interceptor {
	return &authInterceptor{useCase: useCase}
}

func (i *authInterceptor) authenticate(ctx context.Context, procedure string, authHeader string) (context.Context, error) {
	if strings.Contains(procedure, "LoginWithOrgId") {
		return ctx, nil
	}

	if authHeader == "" {
		return ctx, connect.NewError(connect.CodeUnauthenticated, nil)
	}

	token := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}

	session, err := i.useCase.ValidateSession(ctx, token)
	if err != nil {
		return ctx, connect.NewError(connect.CodeUnauthenticated, err)
	}

	return context.WithValue(ctx, ConsoleSessionKey, session), nil
}

func (i *authInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		ctx, err := i.authenticate(ctx, req.Spec().Procedure, req.Header().Get("Authorization"))
		if err != nil {
			return nil, err
		}
		return next(ctx, req)
	}
}

func (i *authInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

// 将来的に使う可能性あり、現在は使用していない。
func (i *authInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		ctx, err := i.authenticate(ctx, conn.Spec().Procedure, conn.RequestHeader().Get("Authorization"))
		if err != nil {
			return err
		}
		return next(ctx, conn)
	}
}
