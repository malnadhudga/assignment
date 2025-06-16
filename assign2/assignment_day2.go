package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	USD = 1.0
	EUR = 0.92
	INR = 83.12
	JPY = 157.45
)

func init() {
	hour := time.Now().Hour()
	if hour < 12 {
		fmt.Println("Good Morning!")
	} else if hour < 18 {
		fmt.Println("Good Afternoon!")
	} else {
		fmt.Println("Good Evening!")
	}
}
func validateInput(args []string) (float64, string, string, error) {

	if len(args) != 4 {
		return 0, "", "", fmt.Errorf("invalid number of arguments")
	}
	amount, err := strconv.ParseFloat(args[1], 64)

	if err != nil || amount <= 0.0 {
		return 0, "", "", fmt.Errorf("invalid amount, please enter the valid amount")
	}

	from := args[2]
	to := args[3]

	if !isSupport(from) {
		return 0, "", "", fmt.Errorf("unsupported source currency %s", from)
	}

	if !isSupport(to) {
		return 0, "", "", fmt.Errorf("unsupported target currency %s", to)
	}

	return amount, from, to, nil

}

func isSupport(currency string) bool {
	return currency == "USD" || currency == "EUR" || currency == "INR" || currency == "JPY"
}

func convertCurrency(amount float64, from, to string) (float64, error) {
	var usdAmount float64
	switch from {
	case "USD":
		usdAmount = amount
	case "EUR":
		usdAmount = amount / EUR
	case "INR":
		usdAmount = amount / INR
	case "JPY":
		usdAmount = amount / JPY
	default:
		return 0, fmt.Errorf("unsupported source currency %s", from)
	}

	switch to {
	case "USD":

		return usdAmount, nil
	case "EUR":

		return usdAmount * EUR, nil
	case "INR":

		return usdAmount * INR, nil
	case "JPY":
		return usdAmount * JPY, nil
	default:
		return 0, fmt.Errorf("unsupported destination currency %s", to)

	}
}

func main() {

	amount, from, to, err := validateInput(os.Args)

	if err != nil {
		fmt.Println(err)
		return
	}
	convertedAmount, err := convertCurrency(amount, from, to)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%.3f %s is equivalent to %.3f %s \n", amount, from, convertedAmount, to)
}
