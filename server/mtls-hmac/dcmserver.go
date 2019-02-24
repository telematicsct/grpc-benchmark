package main

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	pb "github.com/telematicsct/grpc-benchmark/dcm"
	"io"
	"io/ioutil"
	"log"
	"time"
)

// dcmServer is used to implement dcm.DCMServer.
type dcmServer struct {
}

//NewDCMServer creates an returns a new DCM server with JWT token
func NewDCMServer() *dcmServer {
	return &dcmServer{}
}

// DiagnosticData records a route composited of a sequence of points.
// It gets a stream of diagnostic data info, and responds with corresponding data
func (s *dcmServer) DiagnosticData(stream pb.DCMService_DiagnosticDataServer) error {
	start := time.Now()
	log.Println("receiving streaming data: ")
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			log.Println("finished reading stream. took ", time.Since(start))
			return stream.SendAndClose(&pb.DiagResponse{Code: 200, Message: "Done"})
		}
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
}
