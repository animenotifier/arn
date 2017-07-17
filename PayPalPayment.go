package arn

import "strconv"

// PayPalPayment is an approved and exeucted PayPal payment.
type PayPalPayment struct {
	ID       string `json:"id"`
	UserID   string `json:"userId"`
	PayerID  string `json:"payerId"`
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
	Method   string `json:"method"`
	Created  string `json:"created"`
}

// AnimeDollar returns the total amount of anime dollars.
func (payment *PayPalPayment) AnimeDollar() int {
	amount, err := strconv.ParseFloat(payment.Amount, 64)

	if err != nil {
		return 0
	}

	return int(amount * 100)
}

// Save saves the paypal payment in the database.
func (payment *PayPalPayment) Save() error {
	return DB.Set("PayPalPayment", payment.ID, payment)
}
