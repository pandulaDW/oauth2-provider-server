package internal

import (
	"encoding/json"
	"net/http"
)

// DiscoveryResponse represents the response returned by the discovery endpoint
type DiscoveryResponse struct {
	Issuer           string   `json:"issuer"`
	AuthorizationURL string   `json:"authorization_endpoint"`
	TokenURL         string   `json:"token_endpoint"`
	JwksURL          string   `json:"jwks_uri"`
	ResponseTypes    []string `json:"response_types_supported"`
	Scopes           []string `json:"scopes_supported"`
}

func DiscoveryHandler(w http.ResponseWriter, _ *http.Request) {
	// Construct and return OIDC discovery metadata
	discoveryResponse := DiscoveryResponse{
		Issuer:           "http://localhost:9096",
		AuthorizationURL: "http://localhost:9096/authorize",
		TokenURL:         "http://localhost:9096/token",
		JwksURL:          "http://localhost:9096/jwks",
		ResponseTypes:    []string{"code"},
		Scopes:           []string{"openid", "profile", "email"}, // Add supported scopes
	}

	// Convert the response to JSON and write it to the response
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(discoveryResponse)
}
