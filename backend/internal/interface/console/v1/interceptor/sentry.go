package interceptor

import (
	"context"

	"connectrpc.com/connect"
	"github.com/getsentry/sentry-go"
)

type sentryInterceptor struct{}

func NewSentryInterceptor() connect.Interceptor {
	return &sentryInterceptor{}
}

func (i *sentryInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		resp, err := next(ctx, req)
		if err != nil {
			// Connect RPCエラーコードに基づいて、Sentryに送信すべきかを判断
			if connectErr, ok := err.(*connect.Error); ok {
				// 5xxエラー（サーバー内部エラー）のみSentryに送信
				if connectErr.Code() == connect.CodeInternal ||
					connectErr.Code() == connect.CodeUnknown ||
					connectErr.Code() == connect.CodeDataLoss {

					hub := sentry.CurrentHub().Clone()
					hub.WithScope(func(scope *sentry.Scope) {
						scope.SetTag("rpc_method", req.Spec().Procedure)
						scope.SetContext("request", map[string]interface{}{
							"procedure": req.Spec().Procedure,
							"peer":      req.Peer(),
						})
						hub.CaptureException(err)
					})
				}
			}
		}
		return resp, err
	}
}

func (i *sentryInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (i *sentryInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		err := next(ctx, conn)
		if err != nil {
			// Connect RPCエラーコードに基づいて、Sentryに送信すべきかを判断
			if connectErr, ok := err.(*connect.Error); ok {
				// 5xxエラー（サーバー内部エラー）のみSentryに送信
				if connectErr.Code() == connect.CodeInternal ||
					connectErr.Code() == connect.CodeUnknown ||
					connectErr.Code() == connect.CodeDataLoss {

					hub := sentry.CurrentHub().Clone()
					hub.WithScope(func(scope *sentry.Scope) {
						scope.SetTag("rpc_method", conn.Spec().Procedure)
						scope.SetContext("request", map[string]interface{}{
							"procedure": conn.Spec().Procedure,
							"peer":      conn.Peer(),
						})
						hub.CaptureException(err)
					})
				}
			}
		}
		return err
	}
}
