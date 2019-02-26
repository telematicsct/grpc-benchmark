package payload

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"strconv"

	pb "github.com/telematicsct/grpc-benchmark/dcm"
	"github.com/telematicsct/grpc-benchmark/pkg/env"
)

type Payload struct {
	Body []byte `json:"body,omitempty"`
}

type DiagResponse struct {
	Code    int32  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type DiagRecorderData struct {
	CanId   int32    `json:"canId,omitempty"`
	Payload *Payload `json:"payload,omitempty"`
}

const (
	PayloadSizeKey = "PAYLOAD_SIZE"
)

var payloadSize = initPayloadSize()

func initPayloadSize() int {
	size, _ := strconv.Atoi(env.GetString(PayloadSizeKey, "1000000"))
	if size == 0 {
		size = 1000 * 1000 // 1mb
	}
	fmt.Println("payload size", size)
	return size
}

//GetPayload returns a payload
func GetPayload() ([]byte, error) {
	payload := make([]byte, payloadSize)
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

func NewDiagRecorderDataForHTTP() (*DiagRecorderData, error) {
	body, err := GetPayload()
	if err != nil {
		return nil, err
	}
	data := &DiagRecorderData{CanId: GetCanID(), Payload: &Payload{Body: body}}
	return data, nil
}

func GetHTTPJsonPayload() (string, error) {
	data, err := NewDiagRecorderDataForHTTP()
	if err != nil {
		return "", nil
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(data)
	return buf.String(), nil
}
