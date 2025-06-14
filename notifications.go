package paynow_sdk

type Notification struct {
	PaymentId  string `json:"paymentId"`            // Unique identifier for the payment
	ExternalId string `json:"externalId,omitempty"` // Unique identifier for the payment in the merchant's system
	Status     string `json:"status"`               // Status of the payment, Possible values: [NEW, PENDING, ERROR, COMPLETED, CANCELED]
	ModifiedAt string `json:"modifiedAt"`           // Timestamp of the last modification in ISO 8601 format
}
