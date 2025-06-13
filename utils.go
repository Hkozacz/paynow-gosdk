package paynow_sdk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
)

// SignatureBody represents the structure for the signature calculation
type SignatureBody struct {
	Headers    Headers             `json:"headers"`
	Parameters map[string][]string `json:"parameters"`
	Body       string              `json:"body"`
}

type Headers struct {
	ApiKey         string `json:"Api-Key"`
	IdempotencyKey string `json:"Idempotency-Key,omitempty"`
}

func GenerateV3(apiKey, signatureKey, idempotencyKey, data string, parameters map[string]string) (string, error) {
	// Process parameters: convert single values to slices
	parsedParameters := make(map[string][]string)
	for key, value := range parameters {
		parsedParameters[key] = []string{value}

	}
	headers := Headers{
		ApiKey: apiKey,
	}
	if idempotencyKey != "" {
		headers.IdempotencyKey = idempotencyKey
	}

	// Create signature body
	signatureBody := SignatureBody{
		Headers:    headers,
		Parameters: parsedParameters,
		Body:       data,
	}

	// Marshal to JSON
	message, err := json.Marshal(signatureBody)
	if err != nil {
		return "", err
	}

	// Create HMAC signature
	h := hmac.New(sha256.New, []byte(signatureKey))
	h.Write(message)
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signature, nil
}
