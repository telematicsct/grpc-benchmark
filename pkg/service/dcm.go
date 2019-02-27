package service

import (
	"context"
	"io"
	"log"

	"github.com/telematicsct/grpc-benchmark/dcm"
	"github.com/telematicsct/grpc-benchmark/pkg/auth"
)

// DCM is used to implement dcm.DCMServer.
type DCM struct {
	AuthType auth.AuthType
	Token    *auth.JWT
}

//NewDCMService creates an returns a new DCM service
func NewDCMService() *DCM {
	dcm := &DCM{}
	dcm.AuthType = auth.NoAuth
	return dcm
}

//NewDCMServiceWithJWT returns a new DCM service initialized with JWT token
func NewDCMServiceWithJWT(rsaPrivateKeyFile string, rsaPublicKeyFile string) (*DCM, error) {
	jwtToken, err := auth.New(rsaPrivateKeyFile, rsaPublicKeyFile)
	if err != nil {
		return nil, err
	}
	token, err := jwtToken.SampleToken()
	if err != nil {
		return nil, err
	}
	log.Println("mgrpc sample jwt token:", token)
	return &DCM{Token: jwtToken, AuthType: auth.JWTAuth}, nil
}

// DiagnosticDataStream gets a stream of diagnostic data info, and responds with corresponding data
func (s *DCM) DiagnosticDataStream(stream dcm.DCMService_DiagnosticDataStreamServer) error {
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

// DiagnosticDataHMAC gets a diagnostic data info and responds after verifying the JWT token
func (s *DCM) DiagnosticDataHMAC(ctx context.Context, data *dcm.DiagRecorderData) (*dcm.DiagResponse, error) {
	//time.Sleep(50 * time.Millisecond)
	_, err := auth.JWTAuthFunc(s.Token)(ctx)
	if err != nil {
		return nil, err
	}
	response := &dcm.DiagResponse{Code: 200, Message: "Done"}
	return response, nil
}

// DiagnosticData gets a diagnostic data info and responds accordingly
func (s *DCM) DiagnosticData(ctx context.Context, data *dcm.DiagRecorderData) (*dcm.DiagResponse, error) {
	//time.Sleep(50 * time.Millisecond)
	switch s.AuthType {
	case auth.JWTAuth:
		_, err := auth.JWTAuthFunc(s.Token)(ctx)
		if err != nil {
			return nil, err
		}
	}
	response := &dcm.DiagResponse{Code: 200, Message: "Done"}
	return response, nil
}
