package interceptor

import (
	"time"

	"google.golang.org/grpc"
)

// ストリーム形式のRPCに対して処理を挟むためのインターセプター
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	// リクエストを処理するハンドラーを返す
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return StreamServerInterceptorHandler(srv, ss, info, handler)
	}
}

// grpcのServerStream{}をラップし、送受信時にカスタム処理を追加する
type wrappedServerStream struct {
	grpc.ServerStream
}

// メッセージの受信時に処理を挟むカスタムメソッド
func (w *wrappedServerStream) RecvMsg(m interface{}) error {
	Logger("Receive a message (type: %T) at %v", m, time.Now()) // メッセージの受信をログに記録
	return w.ServerStream.RecvMsg(m)
}

// メッセージの送信時に処理を挟むカスタムメソッド
func (w *wrappedServerStream) SendMsg(m interface{}) error {
	Logger("Send a message (type: %T) at %v", m, time.Now()) // メッセージの送信をログに記録
	return w.ServerStream.SendMsg(m)
}

// ServerStreamをラップするカスタム関数
func newWrappedServerStream(ss grpc.ServerStream) grpc.ServerStream {
	return &wrappedServerStream{ss}
}

// リクエストとレスポンスの前後にカスタム処理を挟む関数
func StreamServerInterceptorHandler(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	// 呼び出されたgRPCメソッドをログ出力
	Logger("Stream Request - Method: %s", info.FullMethod)

	err := handler(srv, newWrappedServerStream(ss))
	// レスポンス時にエラーがあればエラーログを出力
	if err != nil {
		Logger("Stream Response - Method: %s, Error: %v", info.FullMethod, err)
	}

	// エラーを返す
	return err
}