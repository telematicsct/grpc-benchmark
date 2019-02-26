package grpc

import (
	"context"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"github.com/telematicsct/grpc-benchmark/dcm"
	"github.com/telematicsct/grpc-benchmark/pkg/auth"
)

var jwtToken *auth.JWT

// dcmServer is used to implement dcm.DCMServer.
type dcmServer struct {
	authType auth.AuthType
}

//NewDCMServer creates an returns a new DCM server with JWT token
func NewDCMServer() *dcmServer {
	dcm := &dcmServer{}
	dcm.authType = auth.NoAuth
	return dcm
}

func NewDCMServerWithJWT(rsaPrivateKeyFile string, rsaPublicKeyFile string) (*dcmServer, error) {
	dcm := &dcmServer{}
	dcm.authType = auth.JWTAuth
	j, err := auth.New(rsaPrivateKeyFile, rsaPublicKeyFile)
	if err != nil {
		return nil, err
	}
	token, err := j.SampleToken()
	if err != nil {
		return nil, err
	}
	log.Println("mgrpc sample jwt token:", token)
	jwtToken = j
	return dcm, nil
}

// DiagnosticDataStream records a route composited of a sequence of points.
// It gets a stream of diagnostic data info, and responds with corresponding data
func (s *dcmServer) DiagnosticDataStream(stream dcm.DCMService_DiagnosticDataStreamServer) error {
	//start := time.Now()
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			//log.Println("finished reading stream. took ", time.Since(start))
			//time.Sleep(50 * time.Millisecond)
			return stream.SendAndClose(&dcm.DiagResponse{Code: 200, Message: "Done"})
		}
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
}

func (s *dcmServer) DiagnosticDataHMAC(ctx context.Context, data *dcm.DiagRecorderData) (*dcm.DiagResponse, error) {
	//time.Sleep(50 * time.Millisecond)
	_, err := jwtAuthFunc(ctx)
	if err != nil {
		return nil, err
	}
	response := &dcm.DiagResponse{Code: 200, Message: "Done"}
	return response, nil
}

func (s *dcmServer) DiagnosticData(ctx context.Context, data *dcm.DiagRecorderData) (*dcm.DiagResponse, error) {
	//time.Sleep(50 * time.Millisecond)
	switch s.authType {
	case auth.JWTAuth:
		_, err := jwtAuthFunc(ctx)
		if err != nil {
			return nil, err
		}
	}
	response := &dcm.DiagResponse{Code: 200, Message: "Done"}
	return response, nil
}

func jwtAuthFunc(ctx context.Context) (context.Context, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "missing context metadata")
	}

	keys, ok := meta[auth.AuthorizationKey]
	if !ok || len(meta[auth.AuthorizationKey]) == 0 {
		return nil, grpc.Errorf(codes.Unauthenticated, "no key provided")
	}

	_, err := jwtToken.Validate(keys[0])
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid token")
	}

	// val, err := grpc_auth.AuthFromMD(ctx, AuthHeader)
	// if err != nil {
	// 	return nil, err
	// }
	// if val != APIKey {
	// 	log.Fatalln("invalid api key")
	// 	return nil, grpc.Errorf(codes.Unauthenticated, "Invalid API key")
	// }
	return ctx, nil
}
