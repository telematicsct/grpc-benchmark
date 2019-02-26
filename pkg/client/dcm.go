package client

import (
	pb "github.com/telematicsct/grpc-benchmark/dcm"
	"google.golang.org/grpc"
)

func NewDCMServiceClient(connectAddr string) (pb.DCMServiceClient, error) {
	return NewDCMServiceClientWithAuth(connectAddr, "")
}

//NewDCMServiceClientWithAuth returns a new DCM service client
func NewDCMServiceClientWithAuth(connectAddr string, token string) (pb.DCMServiceClient, error) {
	conn, err := NewGRPCClient(connectAddr, token)
	if err != nil {
		return nil, err
	}
	return pb.NewDCMServiceClient(conn), nil
}

//NewDCMServiceClientFromConn returns a new DCM service client from a given connection
func NewDCMServiceClientFromConn(conn *grpc.ClientConn) (pb.DCMServiceClient, error) {
	return pb.NewDCMServiceClient(conn), nil
}
