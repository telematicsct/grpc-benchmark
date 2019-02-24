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
	jwtPublicKey *rsa.PublicKey
}

//NewDCMServer creates an returns a new DCM server with JWT token
func NewDCMServer(rsaPublicKey string) (*dcmServer, error) {

	var publicKey *rsa.PublicKey
	if rsaPublicKey != "" {
		data, err := ioutil.ReadFile(rsaPublicKey)
		if err != nil {
			return nil, fmt.Errorf("Error reading the jwt public key: %v", err)
		}

		publicKey, err = jwt.ParseRSAPublicKeyFromPEM(data)
		if err != nil {
			return nil, fmt.Errorf("Error parsing the jwt public key: %s", err)
		}

	}

	return &dcmServer{publicKey}, nil
}

// DiagnosticData records a route composited of a sequence of points.
// It gets a stream of diagnostic data info, and responds with corresponding data
func (s *dcmServer) DiagnosticData(stream pb.DCMService_DiagnosticDataServer) error {
	start := time.Now()
	for {
		diagReq, err := stream.Recv()
		if err == io.EOF {
			log.Println("finished reading stream. took ", time.Since(start))
			return stream.SendAndClose(&pb.DiagResponse{Code: 200, Message: "Done"})
		}
		if err != nil {
			log.Fatal(err)
			return err
		}
		log.Println("receiving diagnostic data: ", diagReq)
	}
}
