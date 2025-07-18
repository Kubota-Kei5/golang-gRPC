package main

import (
	"awsomeProject/pb"
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverAddr = "localhost:50051"

	timeoutDuration = 10 * time.Second
	timeSleep       = 1 * time.Second
)

// Unary RPC
// サーバーにtitleを送り、ファイルに存在するかの確認結果をAlbum型で受け取る関数
func callGetAlbum(client pb.AlbumServiceClient, title string) {
	// リクエストのタイムアウトを設定
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	// GetAlbumメソッドを呼び出してサーバーにリクエストを送り、レスポンスを受け取る
	resp, err := client.GetAlbum(ctx, &pb.GetAlbumRequest{Title: title})
	if err != nil {
		log.Fatalf("client.GetAlbum failed: %v", err)
	}

	log.Printf("response: %v", resp.Album)
}

func main() {
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewAlbumServiceClient(conn)

	// callGetAlbumを実行
	callGetAlbum(client, "Blue Train")
	callGetAlbum(client, "Not Exist Title")
}