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

var jwtToken *auth.JWT

type httpHandler struct {
	opts     *server.ServerOptions
	tlsType  server.TLSType
	authType auth.AuthType
}

func Serve(opts *server.ServerOptions, tlsType server.TLSType, authType auth.AuthType) error {
	return doServe(&httpHandler{opts, tlsType, authType})
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

func newJWT(opts *server.ServerOptions) error {
	j, err := auth.New(opts.JWTPrivateKey, opts.JWTPublicKey)
	if err != nil {
		return err
	}
	token, err := j.SampleToken()
	if err != nil {
		return err
	}
	log.Println("mhttp sample jwt token:", token)
	jwtToken = j
	return nil
}

func doServe(h *httpHandler) error {
	var err error

	switch h.authType {
	case auth.JWTAuth:
		err = newJWT(h.opts)
		if err != nil {
			return err
		}
	}

	var tlsConfig *tls.Config
	switch h.tlsType {
	case server.TLS:
		tlsConfig, err = NewTLSConfig(h.opts)
	case server.MTLS:
		tlsConfig, err = NewMTLSConfig(h.opts)
	}
	if err != nil {
		return err
	}

	listen := h.opts.GetBind(server.HTTP, h.tlsType, h.authType)
	server := &http.Server{
		Addr:      listen,
		TLSConfig: tlsConfig,
		Handler:   h,
	}

	switch h.authType {
	case auth.JWTAuth:
		log.Println("HTTP MTLS HMAC(JWT) Listening at", listen)
	default:
		log.Println("HTTP MTLS Listening at", listen)
	}

	err = server.ListenAndServeTLS(h.opts.ServerCertPath, h.opts.ServerKeyPath)
	if err != nil {
		return err
	}

	return nil
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch h.authType {
	case auth.JWTAuth:
		token := r.Header.Get(auth.AuthorizationKey)
		if token == "" {
			http.Error(w, "token missing", 500)
			return
		}
		_, err := jwtToken.Validate(token)
		if err != nil {
			http.Error(w, "invalid token", 500)
			return
		}
	}

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
