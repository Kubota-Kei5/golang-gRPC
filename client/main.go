package main

import (
	"awsomeProject/pb"
	"log"
	"time"

	"google.golang.org/grpc"
)

var (
	serverAddr = "localhost:50051"     // サーバーのアドレスとポート番号

	timeoutDuration = 10 * time.Second // gRPCリクエストのタイムアウト時間
	timeSleep       = 1 * time.Second  // 通信の挙動確認用のリクエスト間のスリープ時間（待機時間）
)

func main() {
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials())) // gRPCクライアントを作成

	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewAlbumServiceClient(conn) // gRPCクライアントをAlbumServiceClientに変換

	// 各通信方式のリクエストを実行する
}