package main

import (
	"log"
	"net"
	"os"

	pbcollector "github.com/Percona-Lab/qan-api/api/collector"
	pbversion "github.com/Percona-Lab/qan-api/api/version"
	"github.com/Percona-Lab/qan-api/models"
	rservice "github.com/Percona-Lab/qan-api/services/receiver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

// HandleVersion implements version.VersionServer
func (s *server) HandleVersion(ctx context.Context, in *pbversion.VersionRequest) (*pbversion.VersionReply, error) {
	log.Println("Version is requested by:", in.Name)
	return &pbversion.VersionReply{Version: "2.0.0-alpha"}, nil
}

func main() {
	bind, ok := os.LookupEnv("QANAPI_BIND")
	if !ok {
		bind = "127.0.0.1:9911"
	}
	dsn, ok := os.LookupEnv("QANAPI_DSN")
	if !ok {
		dsn = "clickhouse://127.0.0.1:9000?database=pmm"
	}

	db, err := NewDB(dsn)
	if err != nil {
		log.Fatal("DB error", err)
	}

	lis, err := net.Listen("tcp", bind)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	qcm := models.NewQueryClass(db)
	grpcServer := grpc.NewServer()
	pbversion.RegisterVersionServer(grpcServer, &server{})
	pbcollector.RegisterAgentServer(grpcServer, rservice.NewService(qcm))
	reflection.Register(grpcServer)
	log.Printf("QAN-API serve: %v\n", bind)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
