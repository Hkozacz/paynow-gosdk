package paynow_sdk

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
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

func (c *PayNowApiClient) SendPostRequest(endpoint, idempotencyKey string, bodyObj RequestType, responseObj, responseErrorObj interface{}) error {
	body, err := json.Marshal(bodyObj)
	parsedBody := string(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}
	signature, err := GenerateV3(c.apiKey, c.secret, idempotencyKey, parsedBody, nil)
	if err != nil {
		return fmt.Errorf("failed to generate signature: %w", err)
	}
	client := resty.New()
	defer client.Close()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Api-Key", c.apiKey).
		SetHeader("Idempotency-Key", idempotencyKey).
		SetHeader("Signature", signature).
		SetBody(parsedBody).
		SetResult(responseObj).
		SetError(responseErrorObj).
		Post(c.baseUrl + endpoint)
	if err != nil {
		return fmt.Errorf("failed to send POST request: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("error response from server: %s Status: %s", resp.String(), resp.Status())
	}
	return nil
}

func (c *PayNowApiClient) SendGetRequest(endpoint, idempotencyKey string, queryParams RequestType, responseObj, responseErrorObj interface{}) error {
	queryParamsMap := make(map[string]string)
	if queryParams != nil {
		queryParamsBytes, err := json.Marshal(queryParams)
		if err != nil {
			return fmt.Errorf("failed to marshal query parameters: %w", err)
		}
		if err := json.Unmarshal(queryParamsBytes, &queryParamsMap); err != nil {
			return fmt.Errorf("failed to unmarshal query parameters: %w", err)
		}
	}
	signature, err := GenerateV3(c.apiKey, c.secret, idempotencyKey, "", queryParamsMap)
	if err != nil {
		return fmt.Errorf("failed to generate signature: %w", err)
	}
	client := resty.New()
	defer client.Close()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Api-Key", c.apiKey).
		SetHeader("Idempotency-Key", idempotencyKey).
		SetHeader("Signature", signature).
		SetQueryParams(queryParamsMap).
		SetResult(responseObj).
		SetError(responseErrorObj).
		Get(c.baseUrl + endpoint)
	if err != nil {
		return fmt.Errorf("failed to send GET request: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("error response from server: %s Status: %s", resp.String(), resp.Status())
	}
	return nil
}

func (c *PayNowApiClient) CreatePayment(body *CreatePaymentRequest, idempotencyKey string) (*CreatePaymentResponse, error) {
	responseObj := &CreatePaymentResponse{}
	responseErrorObj := &ErrorResponse{}
	err := c.SendPostRequest("payments", idempotencyKey, body, responseObj, responseErrorObj)
	if err != nil {
		if responseErrorObj != nil && responseErrorObj.StatusCode != 0 {
			parsedErrorResponse, parsingErr := json.Marshal(responseErrorObj)
			if parsingErr == nil {
				return nil, fmt.Errorf("error creating payment: %s", string(parsedErrorResponse))
			}
		}
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}
	return responseObj, nil
}

func (c *PayNowApiClient) GetPaymentStatus(paymentId string) (*GetPaymentStatusResponse, error) {
	responseObj := &GetPaymentStatusResponse{}
	responseErrorObj := &ErrorResponse{}
	err := c.SendGetRequest("payments/"+paymentId+"/status", uuid.New().String(), nil, responseObj, responseErrorObj)
	if err != nil {
		if responseErrorObj != nil && responseErrorObj.StatusCode != 0 {
			parsedErrorResponse, parsingErr := json.Marshal(responseErrorObj)
			if parsingErr == nil {
				return nil, fmt.Errorf("error getting payment status: %s", string(parsedErrorResponse))
			}
		}
		return nil, fmt.Errorf("failed to get payment status: %w", err)
	}
	return responseObj, nil
}

func (c *PayNowApiClient) GetPaymentMethods(queryParameters *GetPaymentMethodsQuery) (*[]GetPaymentMethodsResponse, error) {
	responseObj := &[]GetPaymentMethodsResponse{}
	responseErrorObj := &ErrorResponse{}
	err := c.SendGetRequest("payments/paymentmethods", uuid.New().String(), queryParameters, responseObj, responseErrorObj)
	if err != nil {
		if responseErrorObj != nil && responseErrorObj.StatusCode != 0 {
			parsedErrorResponse, parsingErr := json.Marshal(responseErrorObj)
			if parsingErr == nil {
				return nil, fmt.Errorf("error getting payment methods: %s", string(parsedErrorResponse))
			}
		}
		return nil, fmt.Errorf("failed to get payment methods: %w", err)
	}
	return responseObj, nil
}

func (c *PayNowApiClient) GetGDPRClauses() (*[]GetGDPRClausesResponseItem, error) {
	responseObj := &[]GetGDPRClausesResponseItem{}
	responseErrorObj := &ErrorResponse{}
	err := c.SendGetRequest("payments/dataprocessing/notices", uuid.New().String(), nil, responseObj, responseErrorObj)
	if err != nil {
		if responseErrorObj != nil && responseErrorObj.StatusCode != 0 {
			parsedErrorResponse, parsingErr := json.Marshal(responseErrorObj)
			if parsingErr == nil {
				return nil, fmt.Errorf("error getting GDPR clauses: %s", string(parsedErrorResponse))
			}
		}
		return nil, fmt.Errorf("failed to get GDPR clauses: %w", err)
	}
	return responseObj, nil
}

func (c *PayNowApiClient) CreateRefund(paymentId string, body *CreateRefundRequest, idempotencyKey string) (*CreateRefundResponse, error) {
	responseObj := &CreateRefundResponse{}
	responseErrorObj := &ErrorResponse{}
	err := c.SendPostRequest("payments/"+paymentId+"/refunds", idempotencyKey, body, responseObj, responseErrorObj)
	if err != nil {
		if responseErrorObj != nil && responseErrorObj.StatusCode != 0 {
			parsedErrorResponse, parsingErr := json.Marshal(responseErrorObj)
			if parsingErr == nil {
				return nil, fmt.Errorf("error creating refund: %s", string(parsedErrorResponse))
			}
		}
		return nil, fmt.Errorf("failed to create refund: %w", err)
	}
	return responseObj, nil
}

func (c *PayNowApiClient) GetRefundStatus(refundId string) (*GetRefundStatusResponse, error) {
	responseObj := &GetRefundStatusResponse{}
	responseErrorObj := &ErrorResponse{}
	err := c.SendGetRequest("refunds/"+refundId+"/status", uuid.New().String(), nil, responseObj, responseErrorObj)
	if err != nil {
		if responseErrorObj != nil && responseErrorObj.StatusCode != 0 {
			parsedErrorResponse, parsingErr := json.Marshal(responseErrorObj)
			if parsingErr == nil {
				return nil, fmt.Errorf("error getting refund status: %s", string(parsedErrorResponse))
			}
		}
		return nil, fmt.Errorf("failed to get refund status: %w", err)
	}
	return responseObj, nil
}

func (c *PayNowApiClient) CancelRefund(refundId string, idempotencyKey string) (*GetRefundStatusResponse, error) {
	responseObj := &GetRefundStatusResponse{}
	responseErrorObj := &ErrorResponse{}
	err := c.SendPostRequest("refunds/"+refundId+"/cancel", idempotencyKey, nil, responseObj, responseErrorObj)
	if err != nil {
		if responseErrorObj != nil && responseErrorObj.StatusCode != 0 {
			parsedErrorResponse, parsingErr := json.Marshal(responseErrorObj)
			if parsingErr == nil {
				return nil, fmt.Errorf("error canceling refund: %s", string(parsedErrorResponse))
			}
		}
		return nil, fmt.Errorf("failed to cancel refund: %w", err)
	}
	return responseObj, nil
}

func (c *PayNowApiClient) PatchShopURLs(bodyObj *PatchShopURLsRequest, idempotencyKey string) error {
	body, err := json.Marshal(bodyObj)
	parsedBody := string(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}
	signature, err := GenerateV3(c.apiKey, c.secret, idempotencyKey, parsedBody, nil)
	if err != nil {
		return fmt.Errorf("failed to generate signature: %w", err)
	}
	client := resty.New()
	defer client.Close()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Api-Key", c.apiKey).
		SetHeader("Idempotency-Key", idempotencyKey).
		SetHeader("Signature", signature).
		SetBody(parsedBody).
		Patch(c.baseUrl + "configuration/shop/urls")
	if err != nil {
		return fmt.Errorf("failed to send POST request: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("error response from server: %s Status: %s", resp.String(), resp.Status())
	}
	return nil
}
