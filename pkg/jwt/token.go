package jwt

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtPublicKey *rsa.PublicKey

type JWT struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
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

//New creates a new JWT token
func New(rsaPrivateKeyFile string, rsaPublicKeyFile string) (*JWT, error) {
	privateKey, err := parsePrivateKey(rsaPrivateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("Error prasing private key: %s", err)
	}
	publicKey, err := parsePublicKey(rsaPublicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("Error parsing public key: %s", err)
	}
	return &JWT{privateKey, publicKey}, nil
}

func (j *JWT) SampleToken() (string, error) {
	// sample token
	token := jwt.New(jwt.SigningMethodRS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	//claims["admin"] = true
	//claims["iss"] = "auth.service"
	claims["iat"] = time.Now().Unix()
	//claims["email"] = "admin@dcm"
	//claims["sub"] = "admin"
	token.Claims = claims
	tokenString, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", fmt.Errorf("Error signing the sample token: %s", err)
	}
	return tokenString, nil
}

func (j *JWT) Validate(token string) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			log.Printf("Unexpected signing method: %v", t.Header["alg"])
			return nil, fmt.Errorf("invalid token")
		}
		return j.publicKey, nil
	})
	if err == nil && jwtToken.Valid {
		return jwtToken, nil
	}
	return nil, err
}

// type jwt struct {
// 	token string
// }

// // New holds per-rpc metadata for the gRPC clients
// func New(token string) credentials.PerRPCCredentials {
// 	return jwt{token}
// }

// func (j jwt) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
// 	return map[string]string{
// 		"authorization": j.token,
// 	}, nil
// }

// func (j jwt) RequireTransportSecurity() bool {
// 	return true
// }
