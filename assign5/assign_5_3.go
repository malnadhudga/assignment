package main

import (
	"fmt"
)

type PaymentMethod interface {
	Pay(amount float64) string
}

type OTPGenerator interface {
	GenerateOTP() string
}

type CreditCard struct {
	CardHolder string
	CardNumber string
	ExpiryDate string
}

func (cc *CreditCard) Pay(amount float64) string {
	return fmt.Sprintf("[CreditCard] Paid ₹%.2f using card ending with %s", amount, cc.CardNumber[len(cc.CardNumber)-4:])
}

func (cc *CreditCard) GenerateOTP() string {
	return "[CreditCard] OTP sent to the registered number"
}

type PayPal struct {
	Email string
}

func (pp *PayPal) Pay(amount float64) string {
	return fmt.Sprintf("[PayPal] Paid ₹%.2f using PayPal account: %s", amount, pp.Email)
}

type UPI struct {
	UPIID string
}

func (u *UPI) Pay(amount float64) string {
	return fmt.Sprintf("[UPI] Paid ₹%.2f using UPI: %s", amount, u.UPIID)
}

func (u *UPI) GenerateOTP() string {
	return "[UPI] OTP sent to the registered device"
}

func main() {
	fmt.Println("Payment Processing System")

	creditCard := &CreditCard{
		CardHolder: "Ram",
		CardNumber: "1234-4444-9012-2233",
		ExpiryDate: "1/29",
	}

	payPal := &PayPal{
		Email: "user@gmail.com",
	}

	upi := &UPI{
		UPIID: "example@upi",
	}

	paymentOptions := []PaymentMethod{
		creditCard,
		payPal,
		upi,
	}

	paymentAmount := 100.00

	fmt.Printf("\nProcessing the Payment of ₹%.2f for each method :\n\n", paymentAmount)

	for _, method := range paymentOptions {
		if otpMethod, ok := method.(OTPGenerator); ok {
			fmt.Println(otpMethod.GenerateOTP())
		}

		fmt.Println(method.Pay(paymentAmount))

	}

}
