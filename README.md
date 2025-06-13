# Paynow Go SDK

SDK for integrating with the Paynow API in Go. It allows you to handle payments, refunds, retrieve statuses, and available payment methods.

## Installation

Add to your project:

```bash
go get github.com/Hkozacz/paynow-gosdk
```

## Client Initialization

```go
import "github.com/Hkozacz/paynow-gosdk"

client := paynow_sdk.NewPayNowApiClient("API_KEY", "API_SECRET", "https://api.paynow.pl/v3/")
```

## Creating a Payment

```go
paymentReq := &paynow_sdk.CreatePaymentRequest{
    Amount:      1000,
    Currency:    "PLN",
    ExternalId:  "order-123",
    Description: "Product purchase",
    Buyer: &paynow_sdk.BuyerInfo{
        Email: "customer@example.com",
    },
}
resp, err := client.CreatePayment(paymentReq, "unique-idempotency-key")
if err != nil {
    // error handling
}
fmt.Println(resp.RedirectUrl)
```

## Retrieving Payment Status

```go
status, err := client.GetPaymentStatus("paymentId")
if err != nil {
    // error handling
}
fmt.Println(status.Status)
```

## Retrieving Available Payment Methods

```go
methods, err := client.GetPaymentMethods(&paynow_sdk.GetPaymentMethodsQuery{
    Amount:   1000,
    Currency: "PLN",
})
if err != nil {
    // error handling
}
fmt.Println((*methods)[0].Type)
```

## Creating a Refund

```go
refundReq := &paynow_sdk.CreateRefundRequest{
    Amount: 1000,
    Reason: "RMA",
}
refund, err := client.CreateRefund("paymentId", refundReq, "unique-idempotency-key")
if err != nil {
    // error handling
}
fmt.Println(refund.RefundId)
```

## Retrieving Refund Status

```go
refundStatus, err := client.GetRefundStatus("refundId")
if err != nil {
    // error handling
}
fmt.Println(refundStatus.Status)
```

## Retrieving GDPR Clauses

```go
gdpr, err := client.GetGDPRClauses()
if err != nil {
    // error handling
}
fmt.Println((*gdpr)[0].Title)
```

## Structure Validation

Each request structure has a `Validate()` method that is being called before sending a request.

## Request and Response Structures

The SDK provides Go structs for all request and response payloads. These are used to build requests and parse responses from the Paynow API.

### Example Request Structs

#### CreatePaymentRequest
```go
// Required fields are marked with (required)
type CreatePaymentRequest struct {
    Amount        int64        `json:"amount"`                // (required) Amount in the smallest currency unit (e.g., cents)
    Currency      string       `json:"currency,omitempty"`    // (required) ISO 4217 currency code (PLN, EUR, USD, GBP, CZK)
    ExternalId    string       `json:"externalId"`            // (required) Unique order identifier
    Description   string       `json:"description"`           // (required) Payment description (max 255 chars)
    Buyer         *BuyerInfo   `json:"buyer"`                 // (required) Buyer information
    OrderItems    []*OrderItem `json:"orderItems,omitempty"`  // (optional) List of order items
    ContinueUrl   string       `json:"continueUrl,omitempty"` // (optional) Redirect URL after payment
    ValidityTime  int64        `json:"validityTime,omitempty"`// (optional) Payment validity in seconds (60-864000)
    PayoutAccount string       `json:"payoutAccount,omitempty"`// (optional) Account for payout
}
```
- All required fields must be set and pass validation (e.g., Amount > 0, valid currency, non-empty ExternalId, Description, Buyer).
- `OrderItems` is optional but if provided, each item must pass its own validation.

#### BuyerInfo
```go
type BuyerInfo struct {
    Email             string `json:"email"`                    // (required) Buyer's email (max 254 chars)
    FirstName         string `json:"firstName,omitempty"`      // (optional)
    LastName          string `json:"lastName,omitempty"`       // (optional)
    Phone             *Phone `json:"phone,omitempty"`          // (optional, but if set, must be valid)
    Address           *Address `json:"address,omitempty"`      // (optional, but if set, must be valid)
    Locale            string `json:"locale,omitempty"`         // (optional, max 35 chars)
    ExternalId        string `json:"externalId,omitempty"`     // (optional, max 100 chars)
    DeviceFingerprint string `json:"deviceFingerprint,omitempty"` // (optional)
}
```
- Only `Email` is required, but if `Phone` or `Address` are set, they must be valid.

#### CreateRefundRequest
```go
type CreateRefundRequest struct {
    Amount int64  `json:"amount"`           // (required) Amount to refund
    Reason string `json:"reason,omitempty"` // (required) One of: RMA, REFUND_BEFORE_14, REFUND_AFTER_14, OTHER
}
```
- Both fields are required and must pass validation.

#### GetPaymentMethodsQuery
```go
type GetPaymentMethodsQuery struct {
    Amount          int64  `json:"amount,omitempty"`   // (optional)
    Currency        string `json:"currency,omitempty"` // (required) ISO 4217 currency code
    ApplePayEnabled bool   `json:"applePayEnabled,omitempty"` // (optional)
}
```
- `Currency` is required and must be one of the allowed values.

### Example Response Structs

#### CreatePaymentResponse
```go
type CreatePaymentResponse struct {
    RedirectUrl string `json:"redirectUrl"` // URL to redirect the user
    PaymentId   string `json:"paymentId"`   // Payment identifier
    Status      string `json:"status"`      // Payment status
}
```

#### GetPaymentStatusResponse
```go
type GetPaymentStatusResponse struct {
    PaymentId string `json:"paymentId"`
    Status    string `json:"status"`
}
```

#### GetPaymentMethodsResponse
```go
type GetPaymentMethodsResponse struct {
    Type           string          `json:"type"`
    PaymentMethods []PaymentMethod `json:"paymentMethods"`
}
```

#### CreateRefundResponse
```go
type CreateRefundResponse struct {
    RefundId string `json:"refundId"`
    Status   string `json:"status"`
}
```

#### GetRefundStatusResponse
```go
type GetRefundStatusResponse struct {
    RefundId      string `json:"refundId"`
    Status        string `json:"status"`
}
```

All request and response structs are defined in `requests.go` and `responses.go`. For more details, see the source code and the `Validate()` methods for each struct.

---
