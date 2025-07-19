package interceptor

import (
	"context"

	"google.golang.org/grpc"
)

// Unary RPCに対してリクエストの前後に処理を挟むためのインターセプター
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	// リクエストを処理するハンドラーを返す
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return unaryServerInterceptorHandler(ctx, req, info, handler)
	}
}

// リクエストの前後にロギング処理を挟むための関数
func unaryServerInterceptorHandler(
	ctx context.Context,           // gRPCリクエストのコンテキスト（リクエストごとの設定や状態管理）
	req interface{},               // クライアントから送信されたリクエストデータ
	info *grpc.UnaryServerInfo,    // gRPCメソッドの情報（メソッド名など）
	handler grpc.UnaryHandler,     // リクエストを実際に処理するハンドラ関数
) (interface{}, error) {
	// 呼び出されたgRPCメソッドをログ出力
	Logger("Unary Request - Method: %s", info.FullMethod)

	// gRPCリクエストを処理する
	m, err := handler(ctx, req)

	// レスポンス時にエラーがあればエラーログを出力
	if err != nil {
		Logger("Unary Response - Method: %s, Error: %v", info.FullMethod, err)
	}

	// レスポンスを返す
	return m, err
}