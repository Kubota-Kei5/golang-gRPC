package main

import (
	"awsomeProject/pb"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

const (
	filePath = "db/album.json" // JSONファイルに保存されたアルバムデータのパス
	port     = "50051"
)

type AlbumServer struct {
	pb.UnimplementedAlbumServiceServer

	savedAlbums []*pb.Album // サーバーに保存されたアルバムのリスト
}

// Unary RPC
// クライアントから送信されたアルバムのタイトルに基づいて、アルバム情報を返すメソッド
func (s *AlbumServer) GetAlbum(ctx context.Context, req *pb.GetAlbumRequest) (*pb.GetAlbumResponse, error) {
	for _, album := range s.savedAlbums {
		if album.Title == req.Title {
			log.Printf("album found: %s", req.Title)
			return &pb.GetAlbumResponse{Album: album}, nil
		}
	}

	log.Printf("album not found: %s", req.Title)
	return &pb.GetAlbumResponse{Album: &pb.Album{}}, nil
}

// サーバーの初期化時にアルバムデータをロードするメソッド
func (s *AlbumServer) loadAlbums(filePath string) error {
	var (
		data []byte
		err  error
	)
	data, err = os.ReadFile(filePath) // ファイルからアルバム情報を読み取る
	if err != nil {
		return err
	}

    // JSONデータをGoの構造体に変換しsavedAlbumsに保存する
	if err := json.Unmarshal(data, &s.savedAlbums); err != nil {
		return err
	}

	return nil
}

func newServer() *AlbumServer {
	s := &AlbumServer{}
	if err := s.loadAlbums(filePath); err != nil { // サーバー起動時にアルバムデータをロード
		log.Fatalf("failed to load albums: %v", err)
	}

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