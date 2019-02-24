package mtlshttp

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Payload struct {
	Body []byte `json:"body,omitempty"`
}

type DiagResponse struct {
	Code    int32  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type DiagRecorderData struct {
	CanId int32 `json:"canId,omitempty"`
	//Payload *Payload `json:"payload,omitempty"`
}

func Start(addr string, cert string, key string, ca string) error {
	http.HandleFunc("/", CreateDiagRecorderData)
	caCert, err := ioutil.ReadFile(ca)
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
		Addr:      addr,
		TLSConfig: tlsConfig,
	}

	log.Println("Listening at", addr)
	err = server.ListenAndServeTLS(cert, key)
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

	json.NewEncoder(w).Encode(DiagResponse{
		Code:    200,
		Message: "OK",
	})
}
