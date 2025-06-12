package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// SignatureBody represents the structure for the signature calculation
type SignatureBody struct {
	Headers    Headers                `json:"headers"`
	Parameters map[string]interface{} `json:"parameters"`
	Body       string                 `json:"body"`
}

type Headers struct {
	ApiKey         string `json:"Api-Key"`
	IdempotencyKey string `json:"Idempotency-Key"`
}

func GenerateV3(apiKey, signatureKey, idempotencyKey, data string, parameters map[string]interface{}) (string, error) {
	// Process parameters: convert single values to slices
	parsedParameters := make(map[string]interface{})
	for key, value := range parameters {
		switch v := value.(type) {
		case []interface{}:
			parsedParameters[key] = v
		default:
			parsedParameters[key] = []interface{}{v}
		}
	}

	// Create signature body
	signatureBody := SignatureBody{
		Headers: Headers{
			ApiKey:         apiKey,
			IdempotencyKey: idempotencyKey,
		},
		Parameters: parsedParameters,
		Body:       data,
	}
	fmt.Println(signatureBody)

	// Marshal to JSON
	message, err := json.Marshal(signatureBody)
	if err != nil {
		return "", err
	}
	fmt.Println(string(message))

	// Create HMAC signature
	h := hmac.New(sha256.New, []byte(signatureKey))
	h.Write(message)
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signature, nil
}
