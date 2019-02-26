package mgrpc

import (
	"context"
	"crypto/rsa"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	pb "github.com/telematicsct/grpc-benchmark/dcm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type AuthType int

const (
	NoAuth AuthType = iota
	JWTAuth
)

const (
	// AuthHeader defines authorization header.
	AuthHeader = "Authorization"
	// AuthScheme defines authorization scheme.
	AuthScheme = "Bearer"
	// AuthorizationKey is the key used to store authorization token data
	AuthorizationKey = "authorization"
)

var jwtPublicKey *rsa.PublicKey

// dcmServer is used to implement dcm.DCMServer.
type dcmServer struct {
	authType AuthType

	jwtPrivateKey *rsa.PrivateKey
	jwtPublicKey  *rsa.PublicKey

	apiKey string
}

//NewDCMServer creates an returns a new DCM server with JWT token
func NewDCMServer() *dcmServer {
	dcm := &dcmServer{}
	dcm.authType = NoAuth
	return dcm
}

func parsePrivateKey(rsaPrivateKeyFile string) (*rsa.PrivateKey, error) {
	key, err := ioutil.ReadFile(rsaPrivateKeyFile)
	if err != nil {
		return nil, err
	}
	parsedKey, err := jwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil {
		return nil, err
	}
	return parsedKey, nil
}

func parsePublicKey(rsaPublicKeyFile string) (*rsa.PublicKey, error) {
	data, err := ioutil.ReadFile(rsaPublicKeyFile)
	if err != nil {
		return nil, err
	}

	publickey, err := jwt.ParseRSAPublicKeyFromPEM(data)
	if err != nil {
		return nil, err
	}
	return publickey, nil
}

func NewDCMServerWithJWT(rsaPrivateKeyFile string, rsaPublicKeyFile string) (*dcmServer, error) {
	privateKey, err := parsePrivateKey(rsaPrivateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("Error prasing private key: %s", err)
	}
	publicKey, err := parsePublicKey(rsaPublicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("Error parsing public key: %s", err)
	}
	// sample token
	token := jwt.New(jwt.SigningMethodRS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	//claims["admin"] = true
	//claims["iss"] = "auth.service"
	claims["iat"] = time.Now().Unix()
	//claims["email"] = "admin@dcm"
	//claims["sub"] = "admin"
	claims = claims
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("Error signing the token: %s", err)
	}
	log.Println("auth-token:", tokenString)
	dcm := &dcmServer{}
	dcm.authType = JWTAuth
	dcm.jwtPrivateKey = privateKey
	dcm.jwtPublicKey = publicKey
	jwtPublicKey = publicKey
	return dcm, nil
}

// DiagnosticDataStream records a route composited of a sequence of points.
// It gets a stream of diagnostic data info, and responds with corresponding data
func (s *dcmServer) DiagnosticDataStream(stream pb.DCMService_DiagnosticDataStreamServer) error {
	start := time.Now()
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			log.Println("finished reading stream. took ", time.Since(start))
			//time.Sleep(50 * time.Millisecond)
			return stream.SendAndClose(&pb.DiagResponse{Code: 200, Message: "Done"})
		}
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
}

func (s *dcmServer) DiagnosticData(ctx context.Context, data *pb.DiagRecorderData) (*pb.DiagResponse, error) {
	//time.Sleep(50 * time.Millisecond)
	switch s.authType {
	case JWTAuth:
		_, err := jwtAuthFunc(ctx)
		if err != nil {
			return nil, err
		}
	}
	response := &pb.DiagResponse{Code: 200, Message: "Done"}
	return response, nil
}

func jwtAuthFunc(ctx context.Context) (context.Context, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "missing context metadata")
	}

	keys, ok := meta[AuthorizationKey]
	if !ok || len(meta[AuthorizationKey]) == 0 {
		return nil, grpc.Errorf(codes.Unauthenticated, "no key provided")
	}

	_, err := validateToken(keys[0], jwtPublicKey)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid token")
	}

	// val, err := grpc_auth.AuthFromMD(ctx, AuthHeader)
	// if err != nil {
	// 	return nil, err
	// }
	// if val != APIKey {
	// 	log.Fatalln("invalid api key")
	// 	return nil, grpc.Errorf(codes.Unauthenticated, "Invalid API key")
	// }
	return ctx, nil
}

func validateToken(token string, publicKey *rsa.PublicKey) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			log.Printf("Unexpected signing method: %v", t.Header["alg"])
			return nil, fmt.Errorf("invalid token")
		}
		return publicKey, nil
	})
	if err == nil && jwtToken.Valid {
		return jwtToken, nil
	}
	return nil, err
}
