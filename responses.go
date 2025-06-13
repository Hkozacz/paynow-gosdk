package paynow_sdk

type Error struct {
	ErrorType string `json:"errorType"`
	Message   string `json:"message"`
}

type ErrorResponse struct {
	StatusCode int     `json:"statusCode"`
	Errors     []Error `json:"errors"` // List of error messages
}

type PaymentMethod struct {
	Id                int64  `json:"id"`                // Unique identifier for the payment method
	Name              string `json:"name"`              // Name of the payment method, e.g., "Visa", "MasterCard"
	Description       string `json:"description"`       // Description of the payment method, e.g., "Credit card payment"
	Image             string `json:"image"`             // URL to the image representing the payment method
	Status            string `json:"status"`            // Status of the payment method, Possible values: [ENABLED, DISABLED]
	AuthorizationType string `json:"authorizationType"` // Type of authorization, Possible values: [REDIRECT, CODE]

}
type GetPaymentMethodsResponse struct {
	Type           string          `json:"type"`           // Possible values: [APPLE_PAY, BLIK, CARD, ECOMMERCE, GOOGLE_PAY, PAYPO, PBL]
	PaymentMethods []PaymentMethod `json:"paymentMethods"` // List of available payment methods
}

type GetGDPRClausesResponseItem struct {
	Title   string `json:"title"`   // Title of the GDPR clause
	Content string `json:"content"` // Content of the GDPR clause
	Locale  string `json:"locale"`  // Language tag compliant with BCP47/RFC5646
}

type CreatePaymentResponse struct {
	RedirectUrl string `json:"redirectUrl"`
	PaymentId   string `json:"paymentId"`
	Status      string `json:"status"` // "NEW" "PENDING" "ERROR"
}

type GetPaymentStatusResponse struct {
	PaymentId string `json:"paymentId"` // Unique identifier for the payment
	Status    string `json:"status"`    // Status of the payment, Possible values: [NEW, PENDING, ERROR, COMPLETED, CANCELED]
}

type CreateRefundResponse struct {
	RefundId string `json:"refundId"` // Unique identifier for the refund
	Status   string `json:"status"`   // Status of the refund, Possible values: [NEW, PENDING, ERROR, COMPLETED, CANCELED]
}

type GetRefundStatusResponse struct {
	RefundId      string `json:"refundId"`                // Unique identifier for the refund
	Status        string `json:"status"`                  // Status of the refund, Possible values: [NEW, PENDING, SUCCESSFUL, FAILED, CANCELLED]
	FailureReason string `json:"failureReason,omitempty"` // Reason for failure, if applicable Possible values: [CARD_BALANCE_ERROR, BUYER_ACCOUNT_CLOSED, OTHER]
}
