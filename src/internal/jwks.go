package internal

import (
	"crypto/rsa"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwk"
)

var privateKey *rsa.PrivateKey

func init() {
	content, err := os.ReadFile("keys/private.pem")
	if err != nil {
		log.Fatal("failed to read the private key")
	}

	if privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(content); err != nil {
		log.Fatal("failed to encode the private key")
	}
}

func HandleJWKS(w http.ResponseWriter, _ *http.Request) {
	// create a new jwk set
	jwkSet := jwk.NewSet()

	// add the public key to the jwk set
	jwkKey, err := jwk.New(privateKey.Public())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_ = jwkKey.Set("alg", "RS256")
	jwkSet.Add(jwkKey)

	// Return the JWKS in the response
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(jwkSet)
}
