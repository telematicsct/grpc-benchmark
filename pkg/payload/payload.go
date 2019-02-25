package payload

import (
	"crypto/rand"

	pb "github.com/telematicsct/grpc-benchmark/dcm"
)

//GetPayload returns a payload of size 100kb
func GetPayload() ([]byte, error) {
	//100000 - 100kb
	payload := make([]byte, 100000)
	if _, err := rand.Read(payload); err != nil {
		return nil, err
	}
	return payload, nil
}

//GetCanID returns CAN ID 123456789
func GetCanID() int32 {
	return 123456789
}

//NewDiagRecorderData returns random diag recorder data
func NewDiagRecorderData() (*pb.DiagRecorderData, error) {
	payload, err := GetPayload()
	if err != nil {
		return nil, err
	}
	data := &pb.DiagRecorderData{CanId: GetCanID(), Payload: &pb.Payload{Body: payload}}
	return data, nil
}
