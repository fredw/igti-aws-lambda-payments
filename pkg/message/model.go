package message

// Customer represents the customer that purchased the order
type Customer struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Birthday  string `json:"birthday"`
	Gender    string `json:"gender"`
}

// Address represents a contact address
type Address struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Street    string `json:"street"`
	Number    string `json:"number"`
	ZipCode   string `json:"zip_code"`
	City      string `json:"city"`
	State     string `json:"state"`
	Country   string `json:"country"`
	Phone     string `json:"phone"`
}

// OrderItem represents the item of the order
type OrderItem struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	UnitPrice float64 `json:"unit_price"`
}

// Order represents the purchased order
type Order struct {
	Id              string      `json:"id"`
	PaymentMethod   string      `json:"payment_method"`
	ShippingAmount  float64     `json:"shipping_amount"`
	Total           float64     `json:"total"`
	OrderItem       []OrderItem `json:"items"`
	BillingAddress  Address     `json:"billing_address"`
	ShippingAddress Address     `json:"shipping_address"`
}

// Message represents the message
type Message struct {
	Id       *string `json:"id"`
	Provider string  `json:"provider"`
	Order    Order   `json:"order"`
}

// Messages represents a list of messages
type Messages []Message
