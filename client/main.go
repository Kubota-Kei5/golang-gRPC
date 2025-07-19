package main

import (
	"awsomeProject/pb"
	"context"
	"io"
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

// Server Streaming RPC
// サーバーにartistを送り、ファイルに存在するAlbumをすべてAlbum型で受け取る関数
func callListAlbums(client pb.AlbumServiceClient, artist string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	// ListAlbumsメソッドを呼び出してサーバーにリクエストを送り、ストリームを受け取る
	stream, err := client.ListAlbums(ctx, &pb.ListAlbumsRequest{Artist: artist})
	if err != nil {
		log.Fatalf("client.ListAlbums failed: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Printf("stream.Recv failed: %v", err)
			break
		}

		log.Printf("response: %v", resp.Album)
	}
}

// Client Streaming RPC
// サーバーに複数のtitleを送り、ファイルに存在するAlbumの総数・合計金額・メッセージを受け取る関数
func callGetTotalAmount(client pb.AlbumServiceClient) {
	titles := []string{
		"Blue Train",
		"Giant Steps",
		"Speak to Evil",
		"Weather Report",
		"A Portrait in Jazz",
		"Chet Baker Sings",
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	// GetTotalAmountメソッドを呼び出してクライアントストリームを作成
	stream, err := client.GetTotalAmount(ctx)
	if err != nil {
		log.Fatalf("client.GetTotalAmount failed: %v", err)
	}

	// 複数のリクエストをストリームに送信
	for _, title := range titles {
		if err := stream.Send(&pb.GetTotalAmountRequest{Title: title}); err != nil {
			log.Fatalf("client.GetTotalAmount: stream.Send(%s) failed: %v", title, err)
		}

		time.Sleep(timeSleep)
	}

	// サーバーからのレスポンスを受け取る
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("client.GetTotalAmount: stream.CloseAndRecv failed: %v", err)
	}

	log.Printf("response: %v", resp)
}

// Bidirectional Streaming RPC
// サーバーに複数のAlbumを連続で送信し、そのたびにサーバーからのメッセージを受け取る関数
func callUploadAndNotify(client pb.AlbumServiceClient) {
	// アルバムのデータを作成
	albums := []*pb.Album{
		{Title: "New Album", Artist: "New Artist", Price: 10.99},
		{Title: "New Album 2", Artist: "New Artist 2", Price: 20.99},
		{Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
		{Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	// UploadAndNotifyメソッドを呼び出して双方向ストリームを作成
	stream, err := client.UploadAndNotify(ctx)
	if err != nil {
		log.Fatalf("client.UploadAndNotify failed: %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}

			if err != nil {
				log.Fatalf("client.UploadAndNotify: stream.Recv failed: %v", err)
			}

			log.Printf("response: %v", resp)
		}
	}()

	// 複数のリクエストをストリームに送信
	for _, album := range albums {
		req := &pb.UploadAndNotifyRequest{Album: album}
		if err := stream.Send(req); err != nil {
			log.Fatalf("client.UploadAndNotify: stream.Send(%v) failed: %v", album, err)
		}

		time.Sleep(timeSleep)
	}

	// ストリームを閉じる
	if err := stream.CloseSend(); err != nil {
		log.Fatalf("client.UploadAndNotify: stream.CloseSend() failed: %v", err)
	}
	<-waitc

}

func main() {
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewAlbumServiceClient(conn)

	// callGetAlbum(client, "Blue Train")
	// callGetAlbum(client, "Not Exist Title")

	// callListAlbums(client, "Miles Davis")

	// callGetTotalAmount(client)

	// callGetAlbumを実行
	callUploadAndNotify(client)
}