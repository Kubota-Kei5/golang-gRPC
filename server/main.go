package main

import (
	"awsomeProject/pb"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
)

const (
	filePath = "db/album.json" // JSONファイルに保存されたアルバムデータのパス
	port     = "50051"

	timeSleep = 1 * time.Second // レスポンス間のスリープ時間
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

// Server Streaming RPC
// クライアントからartistを受け取り、artistが一致するAlbumをすべてAlbum型で返すメソッド
func (s *AlbumServer) ListAlbums(req *pb.ListAlbumsRequest, stream pb.AlbumService_ListAlbumsServer) error {
	log.Printf("request: %s", req.Artist)

	for _, album := range s.savedAlbums {
		if album.Artist == req.Artist {
			// ストリーム形式のレスポンス
			if err := stream.Send(&pb.ListAlbumsResponse{Album: album}); err != nil {
				return err
			}
			time.Sleep(timeSleep)
		}
	}

	return nil
}

// Client Streaming RPC
// クライアントから複数のtitleを受け取り、ファイルに存在するAlbumの総数・合計金額・メッセージを返すメソッド
func (s *AlbumServer) GetTotalAmount(stream pb.AlbumService_GetTotalAmountServer) error {
	var (
		albumCount int32
		totalAmount float32
	)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(
				&pb.GetTotalAmountResponse{
					AlbumCount: albumCount,
					TotalAmount: totalAmount,
					Message: "success to get total amount",
				},
			)
		}
		if err != nil {
			return err
		}

		albumCount++

		log.Printf("request: %s", req.Title)
		// クライアントから受け取ったタイトルに基づいてアルバムを検索
		for _, album := range s.savedAlbums {
			if album.Title == req.Title {
				totalAmount += album.Price
				break
			}
		}
	}
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