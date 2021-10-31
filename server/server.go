package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"

	"github.com/golang/protobuf/proto"

	//pb "google.golang.org/grpc/examples/route_guide/routeguide"
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
	savedFeatures []*pb.Square // read-only after initialized

	mu         sync.Mutex // protects routeNotes
}

// GetFeature returns the feature at the given point.
func (s *matildaServer) GetSquare(ctx context.Context, point *pb.Point) (*pb.Square, error) {
	for _, feature := range s.savedFeatures {
		if proto.Equal(feature.Location, point) {
			return feature, nil
		}
	}
	// No feature was found, return an unnamed feature
	return &pb.Square{Location: point}, nil
}

// loadSquares loads features from a JSON file.
func (s *matildaServer) loadSquares(filePath string) {
	var data []byte
	if filePath != "" {
		var err error
		data, err = ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Failed to load default features: %v", err)
		}
	} else {
		data = exampleData
	}
	if err := json.Unmarshal(data, &s.savedFeatures); err != nil {
		log.Fatalf("Failed to load default features: %v", err)
	}
}

func newServer() *matildaServer {
	s := &matildaServer{}
	s.loadSquares(*jsonDBFile)
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

//Todo - move to floor 3x3_floor.json
var exampleData = []byte(`[
{
    "location": {
        "latitude": 401827388,
        "longitude": -740294537
    },
    "name": "A1"
}, {
    "location": {
        "latitude": 410564152,
        "longitude": -743685054
    },
    "name": "A2"
}, {
    "location": {
        "latitude": 408472324,
        "longitude": -740726046
    },
    "name": "A3"
}, {
    "location": {
        "latitude": 412452168,
        "longitude": -740214052
    },
    "name": "B1"
}, {
    "location": {
        "latitude": 409146138,
        "longitude": -746188906
    },
    "name": "B2"
}, {
    "location": {
        "latitude": 404701380,
        "longitude": -744781745
    },
    "name": "B3"
}, {
    "location": {
        "latitude": 409642566,
        "longitude": -746017679
    },
    "name": "C1"
}, {
    "location": {
        "latitude": 408031728,
        "longitude": -748645385
    },
    "name": "C2"
}, {
    "location": {
        "latitude": 413700272,
        "longitude": -742135189
    },
    "name": "C3"
}
]`)
