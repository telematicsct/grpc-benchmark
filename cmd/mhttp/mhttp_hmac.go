package mhttp

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/telematicsct/grpc-benchmark/cmd"
)

const (
	AuthorizationKey = "authorization"
)

type AuthType int

const (
	NoAuth AuthType = iota
	JWTAuth
	APIKeyAuth
	HmacAuth
)

var jwt *cmd.JWT

type hmacHandler struct {
	http.Handler
}

func newJWT(cliopts *cmd.CliOptions) error {
	j, err := cmd.NewJWT(cliopts.JWTPrivateKey, cliopts.JWTPublicKey)
	if err != nil {
		return err
	}
	token, err := j.SampleToken()
	if err != nil {
		return err
	}
	log.Println("mhttp sample jwt token:", token)
	jwt = j
	return nil
}

func ServeHMAC(cliopts *cmd.CliOptions) error {
	err := newJWT(cliopts)
	if err != nil {
		return err
	}
	return doServe(cliopts, cliopts.HTTPHMACHostPort, &hmacHandler{}, JWTAuth)
}

func (*hmacHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get(AuthorizationKey)
	if token == "" {
		http.Error(w, "token missing", 500)
		return
	}
	_, err := jwt.Validate(token)
	if err != nil {
		http.Error(w, "invalid token", 500)
		return
	}

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
