package client

import (
	pb "github.com/telematicsct/grpc-benchmark/dcm"
	"google.golang.org/grpc"
)

//NewDCMServiceClient returns a new DCM service client
func NewDCMServiceClient(connectAddr string) (pb.DCMServiceClient, error) {
	conn, err := NewGRPCClient(connectAddr)
	if err != nil {
		return nil, err
	}
	return pb.NewDCMServiceClient(conn), nil
}

//NewDCMServiceClientFromConn returns a new DCM service client from a given connection
func NewDCMServiceClientFromConn(conn *grpc.ClientConn) (pb.DCMServiceClient, error) {
	return pb.NewDCMServiceClient(conn), nil
}
