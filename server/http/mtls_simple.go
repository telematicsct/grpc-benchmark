package http

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/telematicsct/grpc-benchmark/pkg/auth"
	"github.com/telematicsct/grpc-benchmark/pkg/payload"
	"github.com/telematicsct/grpc-benchmark/server"
)

type defaultHandler struct {
	http.Handler
}

func ServeMTLS(opts *server.ServerOptions) error {
	tlsConfig, err := NewMTLSConfig(opts)
	if err != nil {
		return err
	}
	return doServe(tlsConfig, opts, opts.HTTPMTLSHostPort, &defaultHandler{}, auth.NoAuth)
}

func NewTLSConfig(opts *server.ServerOptions) (*tls.Config, error) {
	caCert, err := ioutil.ReadFile(opts.CACertPath)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// setup HTTPS client
	tlsConfig := &tls.Config{
		RootCAs:    caCertPool,
		ClientAuth: tls.NoClientCert,
	}
	return tlsConfig, nil
}

func NewMTLSConfig(opts *server.ServerOptions) (*tls.Config, error) {
	caCert, err := ioutil.ReadFile(opts.CACertPath)
	if err != nil {
		return nil, err
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
	return tlsConfig, nil
}

func doServe(tlsConfig *tls.Config, opts *server.ServerOptions, listen string, handler http.Handler, authType auth.AuthType) error {
	server := &http.Server{
		Addr:      listen,
		TLSConfig: tlsConfig,
		Handler:   handler,
	}

	switch authType {
	case auth.JWTAuth:
		log.Println("HTTP MTLS HMAC(JWT) Listening at", listen)
	default:
		log.Println("HTTP MTLS Listening at", listen)
	}

	err := server.ListenAndServeTLS(opts.ServerCertPath, opts.ServerKeyPath)
	if err != nil {
		return err
	}

	return nil
}

func (*defaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var data payload.DiagRecorderData
	decoder.Decode(&data)
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	//time.Sleep(50 * time.Millisecond)
	json.NewEncoder(w).Encode(payload.DiagResponse{
		Code:    200,
		Message: "OK",
	})
}
