package services

import (
	"testing"
)

func TestGenerateBearerToken(t *testing.T) {
	token, _ := GenerateBearerToken(GetEnv("SWITCHER_API_JWT_SECRET"), GetEnv("SWITCHER_API_DOMAIN_ID"))
	AssertNotNil(t, token)
	// println(token)
}

func TestFetchPaylaodFromGraphQLEndpoint(t *testing.T) {
	domainId := GetEnv("SWITCHER_API_DOMAIN_ID")
	token, _ := GenerateBearerToken(GetEnv("SWITCHER_API_JWT_SECRET"), domainId)

	payload, _ := FetchPayloadFromGraphQLEndpoint(GetEnv("SWITCHER_API_URL"), token, domainId, "default")
	AssertNotNil(t, payload)
	// println(FormatJSON(payload))
}
