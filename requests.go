package paynow_sdk

import (
	"fmt"
	"regexp"
)

type RequestType interface {
	Validate() error
}

type GetPaymentMethodsQuery struct {
	Amount          int64  `json:"amount,omitempty"`
	Currency        string `json:"currency,omitempty"` // ISO 4217 currency code
	ApplePayEnabled bool   `json:"applePayEnabled,omitempty"`
}

func (g *GetPaymentMethodsQuery) Validate() error {
	if g.Currency != "PLN" && g.Currency != "EUR" && g.Currency != "USD" && g.Currency != "GBP" && g.Currency != "CZK" {
		return fmt.Errorf("invalid currency: %s, must be one of [PLN, EUR, USD, GBP, CZK]", g.Currency)
	}
	return nil
}

type Phone struct {
	Prefix string `json:"prefix"` // Phone number prefix, e.g., "+48" for Poland.
	Number int64  `json:"number"` // Phone number without the prefix.
}

func (p *Phone) Validate() error {
	if p.Prefix == "" {
		return fmt.Errorf("phone prefix cannot be empty")
	}
	if p.Number <= 0 {
		return fmt.Errorf("phone number must be a positive integer")
	}
	if p.Number > 9999999999 {
		return fmt.Errorf("phone number is too long, must be less than or equal to 10 digits")
	}
	return nil
}

type AddressType struct {
	Street          string `json:"street,omitempty"`          // Street address, e.g., "123 Main St".
	HouseNumber     string `json:"houseNumber,omitempty"`     // House number, e.g., "12A".
	ApartmentNumber string `json:"apartmentNumber,omitempty"` // Apartment number, e.g., "3B".
	Zipcode         string `json:"zipcode,omitempty"`         // Postal code, e.g., "00-123".
	City            string `json:"city,omitempty"`            // City name, e.g., "Warsaw".
	County          string `json:"county,omitempty"`          // County name, e.g., "Warsaw County".
	Country         string `json:"country,omitempty"`         // Country code, e.g., "PL" for Poland.
}

func (a *AddressType) Validate() error {
	if a.Street != "" && len(a.Street) >= 100 {
		return fmt.Errorf("street address is too long, must be less than or equal to 100 characters")
	}
	if a.HouseNumber != "" && len(a.HouseNumber) >= 16 {
		return fmt.Errorf("house number is too long, must be less than or equal to 16 characters")
	}
	if a.ApartmentNumber != "" && len(a.ApartmentNumber) >= 16 {
		return fmt.Errorf("apartment number is too long, must be less than or equal to 16 characters")
	}
	if a.Zipcode != "" && len(a.Zipcode) >= 16 {
		return fmt.Errorf("zipcode is too long, must be less than or equal to 16 characters")
	}
	matched, err := regexp.MatchString(`^\d{2}-\d{3}$`, a.Zipcode)
	if err != nil {
		return fmt.Errorf("error validating zipcode: %v", err)
	}
	if a.Zipcode != "" && !matched {
		return fmt.Errorf("zipcode must be in the format 'XX-XXX', e.g., '00-123'")
	}
	if a.City != "" && len(a.City) >= 100 && len(a.City) <= 2 {
		return fmt.Errorf("city name is too long, must be between 2 and 100 characters")
	}
	if a.County != "" && len(a.County) >= 100 {
		return fmt.Errorf("county name is too long, must be less than or equal to 100 characters")
	}
	if a.Country != "" && len(a.Country) != 2 {
		return fmt.Errorf("country code must be exactly 2 characters, e.g., 'PL' for Poland")
	}
	return nil
}

type Address struct {
	Billing  *AddressType `json:"billing,omitempty"`  // Billing address, e.g., "123 Main St, Warsaw, PL".
	Shipping *AddressType `json:"shipping,omitempty"` // Shipping address, e.g., "456 Elm St, Warsaw, PL".
}

func (a *Address) Validate() error {
	if a.Billing != nil {
		if err := a.Billing.Validate(); err != nil {
			return fmt.Errorf("billing address validation failed: %v", err)
		}
	}
	if a.Shipping != nil {
		if err := a.Shipping.Validate(); err != nil {
			return fmt.Errorf("shipping address validation failed: %v", err)
		}
	}
	return nil
}

type BuyerInfo struct {
	Email             string `json:"email"`
	FirstName         string `json:"firstName,omitempty"`
	LastName          string `json:"lastName,omitempty"`
	*Phone            `json:"phone,omitempty"`
	*Address          `json:"address,omitempty"`
	Locale            string `json:"locale,omitempty"`            // Buyer's language tag compliant with BCP47/RFC5646.
	ExternalId        string `json:"externalId,omitempty"`        // Unique identifier for the buyer in the merchant's system.
	DeviceFingerprint string `json:"deviceFingerprint,omitempty"` // Fingerprint of Buyer’s device
}

func (b *BuyerInfo) Validate() error {
	if b.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	if len(b.Email) > 50 {
		return fmt.Errorf("email is too long, must be less than or equal to 254 characters")
	}
	if b.FirstName != "" && len(b.FirstName) > 50 {
		return fmt.Errorf("first name is too long, must be less than or equal to 100 characters")
	}
	if b.LastName != "" && len(b.LastName) > 50 {
		return fmt.Errorf("last name is too long, must be less than or equal to 100 characters")
	}
	if b.Phone == nil {
		if err := b.Phone.Validate(); err != nil {
			return fmt.Errorf("phone validation failed: %v", err)
		}
	}
	if b.Address == nil {
		if err := b.Address.Validate(); err != nil {
			return fmt.Errorf("address validation failed: %v", err)
		}
	}
	if b.Locale != "" && len(b.Locale) > 35 {
		return fmt.Errorf("locale is too long, must be less than or equal to 10 characters")
	}
	if b.ExternalId != "" && len(b.ExternalId) > 100 {
		return fmt.Errorf("externalId is too long, must be less than or equal to 50 characters")
	}
	return nil
}

type OrderItem struct {
	Name     string `json:"name"`               // Name of the item, e.g., "T-shirt".
	Producer string `json:"producer,omitempty"` // Producer of the item, e.g., "BrandX".
	Category string `json:"category"`           // Category of the item, e.g., "Clothing".
	Quantity int64  `json:"quantity"`           // Quantity of the item, e.g., 2.
	Price    int64  `json:"price"`              // Price of the item in the smallest currency unit, e.g., cents.
}

func (o *OrderItem) Validate() error {
	if o.Name == "" {
		return fmt.Errorf("item name cannot be empty")
	}
	if len(o.Name) > 120 {
		return fmt.Errorf("item name is too long, must be less than or equal to 100 characters")
	}
	if o.Producer != "" && len(o.Producer) > 120 {
		return fmt.Errorf("producer name is too long, must be less than or equal to 100 characters")
	}
	if o.Category == "" || len(o.Category) > 1000 {
		return fmt.Errorf("category name is too long, must be less than or equal to 100 characters")
	}
	if o.Quantity > 9999999999 || o.Quantity <= 0 {
		return fmt.Errorf("quantity must be a positive integer and less than or equal to 10 digits")
	}
	if o.Price <= 0 || o.Price > 9999999999 {
		return fmt.Errorf("price must be a positive integer and less than or equal to 10 digits")
	}
	return nil
}

type CreatePaymentRequest struct {
	Amount        int64        `json:"amount"`
	Currency      string       `json:"currency,omitempty"` // ISO 4217 currency code
	ExternalId    string       `json:"externalId"`
	Description   string       `json:"description"`
	Buyer         *BuyerInfo   `json:"buyer"`
	OrderItems    []*OrderItem `json:"orderItems,omitempty"` // List of items in the order
	ContinueUrl   string       `json:"continueUrl,omitempty"`
	ValidityTime  int64        `json:"validityTime,omitempty"`  // Time in seconds until the payment expires
	PayoutAccount string       `json:"payoutAccount,omitempty"` // Account to which the payment will be made, e.g., "PL61109010140000071219812874".
}

func (c *CreatePaymentRequest) Validate() error {
	if c.Amount <= 0 || c.Amount > 9999999999 {
		return fmt.Errorf("amount must be a positive integer and less than or equal to 10 digits")
	}
	if c.Currency != "PLN" && c.Currency != "EUR" && c.Currency != "USD" && c.Currency != "GBP" && c.Currency != "CZK" {
		return fmt.Errorf("invalid currency: %s, must be one of [PLN, EUR, USD, GBP, CZK]", c.Currency)
	}
	if c.ExternalId == "" || len(c.ExternalId) > 100 {
		return fmt.Errorf("externalId cannot be empty and must be less than or equal to 100 characters")
	}
	if c.Description == "" || len(c.Description) > 255 {
		return fmt.Errorf("description cannot be empty and must be less than or equal to 255 characters")
	}
	if c.Buyer != nil {
		if err := c.Buyer.Validate(); err != nil {
			return fmt.Errorf("buyer validation failed: %v", err)
		}
	}
	if len(c.OrderItems) != 0 {
		for _, item := range c.OrderItems {
			if err := item.Validate(); err != nil {
				return fmt.Errorf("order item validation failed: %v", err)
			}
		}
	}
	if c.ContinueUrl != "" && len(c.ContinueUrl) > 1000 {
		return fmt.Errorf("continueUrl is too long, must be less than or equal to 1000 characters")
	}
	if c.ValidityTime < 60 || c.ValidityTime > 864000 {
		return fmt.Errorf("validityTime must be between 60 seconds and 10 days (864000 seconds)")
	}
	return nil
}

type CreateRefundRequest struct {
	Amount int64  `json:"amount"`           // Amount to refund in the smallest currency unit, e.g., cents.
	Reason string `json:"reason,omitempty"` // Reason for the refund, "RMA", "REFUND_BEFORE_14", "REFUND_AFTER_14", "OTHER"
}

func (c *CreateRefundRequest) Validate() error {
	if c.Amount <= 0 || c.Amount > 9999999999 {
		return fmt.Errorf("amount must be a positive integer and less than or equal to 10 digits")
	}
	if c.Reason != "RMA" && c.Reason != "REFUND_BEFORE_14" && c.Reason != "REFUND_AFTER_14" && c.Reason != "OTHER" {
		return fmt.Errorf("invalid reason: %s, must be one of [RMA, REFUND_BEFORE_14, REFUND_AFTER_14, OTHER]", c.Reason)
	}
	return nil
}

type PatchShopURLsRequest struct {
	NotificationUrl string `json:"notificationUrl,omitempty"` // URL for payment notifications
	ContinueUrl     string `json:"continueUrl,omitempty"`     // URL to redirect after payment completion
}

func (p *PatchShopURLsRequest) Validate() error {
	return nil
}
