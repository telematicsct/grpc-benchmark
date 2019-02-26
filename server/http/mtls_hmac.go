package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/telematicsct/grpc-benchmark/pkg/auth"
	"github.com/telematicsct/grpc-benchmark/pkg/payload"
	"github.com/telematicsct/grpc-benchmark/server"
)

var jwtToken *auth.JWT

type hmacHandler struct {
	http.Handler
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

func ServeMTLSHMAC(opts *server.ServerOptions) error {
	err := newJWT(opts)
	if err != nil {
		return err
	}
	tlsConfig, err := NewMTLSConfig(opts)
	if err != nil {
		return err
	}
	return doServe(tlsConfig, opts, opts.HTTPMTLSHmacHostPort, &hmacHandler{}, auth.JWTAuth)
}

func (*hmacHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
