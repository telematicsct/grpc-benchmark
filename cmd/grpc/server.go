package grpc

import (
	"context"
	"crypto/rsa"
	pb "github.com/telematicsct/grpc-benchmark/dcm"
	"io"
	"log"
	"time"
)

// dcmServer is used to implement dcm.DCMServer.
type dcmServer struct {
	jwtPublicKey *rsa.PublicKey
}

//NewDCMServer creates an returns a new DCM server with JWT token
func NewDCMServer() *dcmServer {
	return &dcmServer{}
}

// DiagnosticDataStream records a route composited of a sequence of points.
// It gets a stream of diagnostic data info, and responds with corresponding data
func (s *dcmServer) DiagnosticDataStream(stream pb.DCMService_DiagnosticDataStreamServer) error {
	start := time.Now()
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			log.Println("finished reading stream. took ", time.Since(start))
			time.Sleep(50 * time.Millisecond)
			return stream.SendAndClose(&pb.DiagResponse{Code: 200, Message: "Done"})
		}
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
}

func (s *dcmServer) DiagnosticData(ctx context.Context, data *pb.DiagRecorderData) (*pb.DiagResponse, error) {
	time.Sleep(50 * time.Millisecond)
	response := &pb.DiagResponse{Code: 200, Message: "Done"}
	return response, nil
}
