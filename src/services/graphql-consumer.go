package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

type GraphQLRequest struct {
	Query string `json:"query"`
}

func GenerateBearerToken(apiKey string, subject string) (string, error) {
	// Define the claims for the JWT token
	claims := jwt.MapClaims{
		"iss":     "GitOps Service",
		"sub":     "/resource",
		"subject": subject,
		"exp":     time.Now().Add(time.Minute).Unix(),
	}

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the API key
	signedToken, err := token.SignedString([]byte(apiKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate bearer token: %v", err)
	}

	return signedToken, nil
}

func FetchPayloadFromGraphQLEndpoint(url string, token string, domainId string, environment string) (string, error) {
	// Define the GraphQL query
	query := createQuery(domainId, environment)

	// Create a new request
	reqBody, _ := json.Marshal(GraphQLRequest{Query: query})
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read and print the response
	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}

func createQuery(domainId string, environment string) string {
	return fmt.Sprintf(`
    {
        domain(_id: "%s", environment: "%s") {
            name
            version
            group {
                name
                description
                activated
                config {
                    key
                    description
                    activated
                    strategies {
                        strategy
                        activated
                        operation
                        values
                    }
                    components
                }
            }
        }
    }`, domainId, environment)
}
