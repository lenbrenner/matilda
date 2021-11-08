package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	pb "takeoff.com/matilda/api"
	"takeoff.com/matilda/applications"
	"takeoff.com/matilda/data"
	"takeoff.com/matilda/model"
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
	application applications.Application
}

func transitionsToPB(transitions []model.Transition) []*pb.Transition {
	pbTransitions := make([]*pb.Transition, len(transitions))
	for i, transition := range transitions {
		pbTransitions[i] = &pb.Transition{
			Direction: int32(transition.Direction),
			Destination: string(transition.Destination),
		}
	}
	return pbTransitions
}
func (s *matildaServer) GetLocation(ctx context.Context, point *pb.Point) (*pb.Location, error) {
	//Todo - add Dao method GetByPoint
	for _, location := range s.application.LocationService.GetAll() {
		locationPoint := &pb.Point{Latitude: location.Latitude, Longitude: location.Longitude}
		if proto.Equal(locationPoint, point) {
			return &pb.Location{
				Label: string(location.Label),
				Location: locationPoint,
				Transitions: transitionsToPB(location.Transitions),
			}, nil
		}
	}
	// No location was found, return an unnamed location
	return &pb.Location{Location: point}, nil
}

func newServer() *matildaServer {
	return &matildaServer{application: applications.Get()}
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