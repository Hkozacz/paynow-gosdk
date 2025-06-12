package main

import (
	"encoding/json"
	"fmt"
	"resty.dev/v3"
)

type PayNowApiClient struct {
	apiKey  string
	secret  string
	baseUrl string
}

func NewPayNowApiClient(apiKey, secret, baseUrl string) *PayNowApiClient {
	return &PayNowApiClient{
		apiKey:  apiKey,
		secret:  secret,
		baseUrl: baseUrl,
	}
}

func (c *PayNowApiClient) CreatePayment(orderId string, amount int64, description, buyer string) (string, error) {
	bodyObj := map[string]interface{}{
		"amount":      amount,
		"externalId":  orderId,
		"description": description,
		"buyer": map[string]string{
			"email": buyer,
		},
	}
	body, err := json.Marshal(bodyObj)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}
	parameters := map[string]interface{}{}
	signature, _ := GenerateV3(c.apiKey, c.secret, orderId, string(body), parameters)
	fmt.Println(signature)
	client := resty.New()
	defer client.Close()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Api-Key", c.apiKey).
		SetHeader("Idempotency-Key", orderId).
		SetHeader("Signature", signature).
		SetBody(string(body)).
		Post("https://api.sandbox.paynow.pl/v3/payments")
	if err != nil {
		return "", fmt.Errorf("failed to create payment: %w", err)
	}
	if resp.IsError() {
		return "", fmt.Errorf("error response from server: %s", resp.String())
	}
	fmt.Println(resp.String())
	return resp.String(), nil
}
