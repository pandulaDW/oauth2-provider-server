package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"example.com/src/internal"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/golang-jwt/jwt"
)

var privateKey []byte

func init() {
	var err error
	if privateKey, err = os.ReadFile("keys/private.pem"); err != nil {
		log.Fatal("failed to read the private key")
	}
}

func main() {
	// create a new manage instance
	manager := manage.NewDefaultManager()

	// in-memory token store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// set the access-token generation to be based on JWT.
	g := generates.NewJWTAccessGenerate("", privateKey, jwt.GetSigningMethod("RS256"))
	manager.MapAccessGenerate(g)

	// set up a client memory store and add it to the manager
	clientStore := store.NewClientStore()
	manager.MapClientStorage(clientStore)

	// add a new client to the store
	_ = clientStore.Set("000000", &models.Client{
		ID:     "000000",
		Secret: "999999",
		Domain: "http://localhost:4180/oauth2/callback", // callback url
	})

	// create authorization server
	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)

	// set error handler endpoints
	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})
	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	// Set the user authorization handler. We can check if the user is valid here.
	srv.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		return "my-user-id-1", nil
	})

	srv.SetResponseTokenHandler(func(w http.ResponseWriter, data map[string]interface{}, _ http.Header, _ ...int) error {
		pKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privateKey)

		token := jwt.New(jwt.SigningMethodRS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["sub"] = "my-user-id-1"
		claims["iss"] = "localhost:9096"
		claims["aud"] = "000000"
		claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
		claims["iat"] = time.Now().Unix()
		signedToken, _ := token.SignedString(pKey)

		data["id_token"] = signedToken
		data["access_token"] = signedToken

		w.WriteHeader(http.StatusOK)
		return json.NewEncoder(w).Encode(data)
	})

	// OIDC endpoints
	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		log.Println("authorization endpoint was called")
		err := srv.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})
	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		log.Println("token endpoint was called")
		err := srv.HandleTokenRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	// discovery endpoint and the JWKs endpoint
	http.HandleFunc("/.well-known/openid-configuration", internal.DiscoveryHandler)
	http.HandleFunc("/jwks", internal.HandleJWKS)

	// callback page
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "src/web/index.html")
	})
	http.HandleFunc("/main.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "src/web/main.js")
	})

	// run the authorization server
	log.Println("Authorization server running at port 9096...")
	log.Fatal(http.ListenAndServe(":9096", nil))
}
