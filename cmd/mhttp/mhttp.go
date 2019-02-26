package mhttp

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/telematicsct/grpc-benchmark/cmd"
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

func Serve(cliopts *cmd.CliOptions) error {
	http.HandleFunc("/", CreateDiagRecorderData)
	caCert, err := ioutil.ReadFile(cliopts.CACertPath)
	if err != nil {
		return err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// setup HTTPS client
	tlsConfig := &tls.Config{
		ClientCAs: caCertPool,
		// NoClientCent
		// RequestClientCert
		// RequiredAnyClientCert
		// VerifyClientCartIfGiven
		// RequireAndVerifyClientCert
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()

	server := &http.Server{
		Addr:      cliopts.HTTPHostPort,
		TLSConfig: tlsConfig,
	}

	log.Println("HTTP MTLS Listening at", cliopts.HTTPHostPort)
	err = server.ListenAndServeTLS(cliopts.ServerCertPath, cliopts.ServerKeyPath)
	if err != nil {
		return err
	}
	return nil
}

func CreateDiagRecorderData(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var data DiagRecorderData
	decoder.Decode(&data)
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	//time.Sleep(50 * time.Millisecond)
	json.NewEncoder(w).Encode(DiagResponse{
		Code:    200,
		Message: "OK",
	})
}
