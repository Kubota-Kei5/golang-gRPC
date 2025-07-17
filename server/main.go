package main

import (
	"awsomeProject/pb"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

const (
	filePath = "db/album.json" // アルバム情報が格納されたJSONファイルへのパス(後ほど追加します。)
	port     = "50051"         // サーバーが待ち受けるポート番号
)

// pb.AlbumServiceServerインターフェースを満たすサーバーを定義する
// 未実装のメソッドはpb.UnimplementedAlbumServiceServerのメソッドが使用される
type AlbumServer struct {
	pb.UnimplementedAlbumServiceServer
}

func newServer() *AlbumServer {
	s := &AlbumServer{}
    // 後ほど初期化時に実行する処理を追加します。

	return s
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAlbumServiceServer(grpcServer, newServer()) // 作成したサーバーをgrpcServerに登録

	log.Println("server started")
	if err := grpcServer.Serve(lis); err != nil { // grpcServerを起動
		log.Fatalf("failed to serve: %v", err)
	}
}