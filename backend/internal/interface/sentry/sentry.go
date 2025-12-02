package sentry

import (
	"context"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/cockroachdb/errors"
	"github.com/getsentry/sentry-go"
	domainerrors "github.com/shibayama-club/keyhub/internal/domain/errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type IErrorInterceptor interface {
	// エラーをSentryに送信し、イベントIDを返す
	captureError(ctx context.Context, err error, code connect.Code, procedure string, peer connect.Peer) string
	// ドメインエラーをConnect RPCエラーに変換する
	connectError(err error) *connect.Error
	// エラーコードに基づいて適切なレベルでエラーをログ出力する
	logError(ctx context.Context, err error, code connect.Code, eventID string)
	// 適切な詳細を含む最終的なエラーレスポンスを構築する
	buildErrorResponse(ctx context.Context, err error, code connect.Code, eventID string) *connect.Error
	WrapUnary(next connect.UnaryFunc) connect.UnaryFunc
	WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc
	WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc
}

type ErrorInterceptor struct {
	enableDetailedErrors bool
}

func NewErrorInterceptor(enableDetailedErrors bool) IErrorInterceptor {
	return &ErrorInterceptor{
		enableDetailedErrors: enableDetailedErrors,
	}
}

// captureErrorInSentry はエラーをSentryに送信し、イベントIDを返す
func (i *ErrorInterceptor) captureError(ctx context.Context, err error, code connect.Code, procedure string, peer connect.Peer) string {
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("error_code", code.String())
		scope.SetTag("rpc_procedure", procedure)
		scope.SetContext("request", map[string]interface{}{
			"procedure": procedure,
			"peer":      peer.Addr,
			"protocol":  peer.Protocol,
		})

		// スタックトレースが利用可能な場合は追加
		if stackTrace := errors.GetReportableStackTrace(err); stackTrace != nil {
			scope.SetContext("stack_trace", map[string]interface{}{
				"frames": stackTrace,
			})
		}

		// ctxからユーザーコンテキストを追加（将来の拡張）
		// if userID, ok := domain.Value[model.UserID](ctx); ok {
		//     scope.SetUser(sentry.User{ID: userID.String()})
		// }
	})

	eventID := hub.CaptureException(err)
	if eventID != nil {
		return string(*eventID)
	}
	return ""
}

// toConnectError はドメインエラーをConnect RPCエラーに変換する
func (i *ErrorInterceptor) connectError(err error) *connect.Error {
	switch {
	case errors.Is(err, domainerrors.ErrValidation):
		return connect.NewError(connect.CodeInvalidArgument, err)
	case errors.Is(err, domainerrors.ErrNotFound):
		return connect.NewError(connect.CodeNotFound, err)
	case errors.Is(err, domainerrors.ErrUnAuthorized):
		return connect.NewError(connect.CodeUnauthenticated, err)
	case errors.Is(err, domainerrors.ErrAlreadyExists):
		return connect.NewError(connect.CodeAlreadyExists, err)
	case errors.Is(err, domainerrors.ErrInternal):
		return connect.NewError(connect.CodeInternal, err)
	default:
		// 不明なエラーは内部エラーとして扱う
		return connect.NewError(connect.CodeUnknown, err)
	}
}

// logError はエラーコードに基づいて適切なレベルでエラーをログ出力する
func (i *ErrorInterceptor) logError(ctx context.Context, err error, code connect.Code, eventID string) {
	baseAttrs := []any{
		slog.Any("error", err),
		slog.String("code", code.String()),
		slog.String("event_id", eventID),
	}

	// スタックトレースをデバッグレベルで追加
	if stackTrace := errors.GetReportableStackTrace(err); stackTrace != nil {
		slog.DebugContext(ctx, "error stack trace",
			slog.Any("stack_trace", stackTrace),
		)
	}

	// エラーコードに基づいて適切なレベルでログ出力
	switch code {
	case connect.CodeInvalidArgument,
		connect.CodeNotFound,
		connect.CodeAlreadyExists,
		connect.CodeUnauthenticated,
		connect.CodePermissionDenied:
		// クライアント側エラー - 警告として記録
		slog.WarnContext(ctx, "request failed due to client error", baseAttrs...)
	default:
		// サーバー側エラーやその他 - エラーとして記録
		slog.ErrorContext(ctx, "request failed due to server error", baseAttrs...)
	}
}

// buildErrorResponse は適切な詳細を含む最終的なエラーレスポンスを構築する
func (i *ErrorInterceptor) buildErrorResponse(ctx context.Context, err error, code connect.Code, eventID string) *connect.Error {
	var newError *connect.Error

	if i.enableDetailedErrors {
		// 開発環境: 詳細なエラーメッセージを含む
		newError = connect.NewError(code, err)
	} else {
		// 本番環境: エラー詳細を隠蔽し、イベントIDのみ表示
		message := fmt.Sprintf("An error occurred. Event ID: %s", eventID)
		newError = connect.NewError(code, errors.New(message))
	}

	// メタデータにイベントIDを追加
	newError.Meta().Set("Error-ID", eventID)

	// エラーヒントが利用可能な場合は日本語メッセージとして追加
	hints := errors.FlattenHints(err)
	if len(hints) > 0 {
		detail, detailErr := connect.NewErrorDetail(&errdetails.LocalizedMessage{
			Locale:  "ja-JP",
			Message: hints,
		})
		if detailErr != nil {
			slog.ErrorContext(ctx, "failed to create error detail",
				slog.String("error", detailErr.Error()),
			)
		} else {
			newError.AddDetail(detail)
		}
	}

	return newError
}

// WrapUnary はエラーハンドリングを伴うUnary RPC呼び出しをラップする
func (i *ErrorInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		res, err := next(ctx, req)
		if err == nil {
			return res, nil
		}

		// 必要に応じてConnectエラーに変換
		var connectErr *connect.Error
		if !errors.As(err, &connectErr) {
			connectErr = i.connectError(err)
		}

		code := connectErr.Code()

		// Sentryにエラーをキャプチャ（サーバー側エラーのみ）
		var eventID string
		if shouldCaptureInSentry(code) {
			eventID = i.captureError(ctx, err, code, req.Spec().Procedure, req.Peer())
		}

		// エラーをログ出力
		i.logError(ctx, err, code, eventID)

		// エラーレスポンスを構築して返却
		return nil, i.buildErrorResponse(ctx, err, code, eventID)
	}
}

// WrapStreamingClient はストリーミングクライアント呼び出しをラップする（現在は何もしない）
func (i *ErrorInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

// WrapStreamingHandler はエラーハンドリングを伴うストリーミングハンドラー呼び出しをラップする
func (i *ErrorInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		err := next(ctx, conn)
		if err == nil {
			return nil
		}

		// 必要に応じてConnectエラーに変換
		var connectErr *connect.Error
		if !errors.As(err, &connectErr) {
			connectErr = i.connectError(err)
		}

		code := connectErr.Code()

		// Sentryにエラーをキャプチャ（サーバー側エラーのみ）
		var eventID string
		if shouldCaptureInSentry(code) {
			eventID = i.captureError(ctx, err, code, conn.Spec().Procedure, conn.Peer())
		}

		// エラーをログ出力
		i.logError(ctx, err, code, eventID)

		// ストリーミングの場合、エラーを直接返す
		// クライアントはエラー詳細を受信する
		if i.enableDetailedErrors {
			return connectErr
		}

		// 本番環境ではエラー詳細を隠蔽
		message := fmt.Sprintf("An error occurred. Event ID: %s", eventID)
		newError := connect.NewError(code, errors.New(message))
		newError.Meta().Set("Error-ID", eventID)

		return newError
	}
}

// shouldCaptureInSentry はエラーをSentryに送信すべきか判定する
// サーバー側エラー（5xx）のみキャプチャされる
func shouldCaptureInSentry(code connect.Code) bool {
	switch code {
	case connect.CodeInternal,
		connect.CodeUnknown,
		connect.CodeDataLoss,
		connect.CodeUnimplemented,
		connect.CodeUnavailable:
		return true
	default:
		return false
	}
}
