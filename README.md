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

client := paynow_sdk.NewPayNowApiClient("API_KEY", "API_SECRET", "https://api.paynow.pl/v1/")
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

---
