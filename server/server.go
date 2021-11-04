package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"net"

	"takeoff.com/matilda/data"
	pb "takeoff.com/matilda/api"
)

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port       = flag.Int("port", 10000, "The server port")
)

type matildaServer struct {
	pb.UnimplementedMatildaServer
	savedLocations []*pb.Location
}

func (s *matildaServer) GetLocation(ctx context.Context, point *pb.Point) (*pb.Location, error) {
	for _, location := range s.savedLocations {
		if proto.Equal(location.Location, point) {
			return location, nil
		}
	}
	// No location was found, return an unnamed location
	return &pb.Location{Location: point}, nil
}

// loadLocations loads features from a JSON file.
func (s *matildaServer) loadLocations(filePath string) {
	var data []byte
	if filePath != "" {
		var err error
		data, err = ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Failed to load default features: %v", err)
		}
	}
	if err := json.Unmarshal(data, &s.savedLocations); err != nil {
		log.Fatalf("Failed to load default features: %v", err)
	}
	fmt.Println("hello")
}

func newServer() *matildaServer {
	s := &matildaServer{}
	//s.loadLocations(*jsonDBFile)
	resourcePath := data.Path("maps/lat_lon_floor.json")
	s.loadLocations(resourcePath)
	return s
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if *tls {
		if *certFile == "" {
			*certFile = data.Path("x509/server_cert.pem")
		}
		if *keyFile == "" {
			*keyFile = data.Path("x509/server_key.pem")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterMatildaServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}