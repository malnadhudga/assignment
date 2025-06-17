package main

import "fmt"

type BankAccount struct {
	Owner   string
	Balance float64
}

// DisplayBalance uses a value receiver. This means it operates on a copy
// of the BankAccount. Any modifications to 'b' within this method would
// not affect the original variable. This is suitable for methods that
// only read data.
func (b BankAccount) DisplayBalance() {
	fmt.Printf("%s's current balance: $%.2f\n", b.Owner, b.Balance)
}

// Deposit uses a pointer receiver (*BankAccount). This allows the method
// to modify the original BankAccount variable directly. This is essential
// for methods that need to change the state of the receiver, like updating
// the balance.
func (b *BankAccount) Deposit(amount float64) {
	if amount <= 0 {
		fmt.Printf("Deposit amount must be positive. Skipping deposit of $%.2f.\n", amount)
		return
	}
	b.Balance += amount // Modifies the actual BankAccount's balance
	fmt.Printf("Deposited $%.2f. New balance: $%.2f\n", amount, b.Balance)
}

// Withdraw uses a pointer receiver (*BankAccount). Similar to Deposit,
// it can directly modify the account's balance. It also includes logic
// to prevent withdrawals that would result in a negative balance.
func (b *BankAccount) Withdraw(amount float64) {
	if amount <= 0 {
		fmt.Printf("Withdrawal amount must be positive. Skipping withdrawal of $%.2f.\n", amount)
		return
	}
	if b.Balance >= amount {
		b.Balance -= amount // Modifies the actual BankAccount's balance
		fmt.Printf("Withdrew $%.2f. New balance: $%.2f\n", amount, b.Balance)
	} else {
		fmt.Printf("Insufficient funds for withdrawal of $%.2f. Current balance: $%.2f.\n", amount, b.Balance)
	}
}

// TryToModifyBalance is a function demonstrating pass-by-value. When
// a BankAccount is passed to this function, only a copy is made.
// Therefore, any changes made to 'b' inside this function will not
// impact the original variable outside of it.
func TryToModifyBalance(b BankAccount, amount float64) {
	b.Balance += amount // Only modifies the copy of BankAccount
	fmt.Printf("   (Inside TryToModifyBalance func): Balance changed to $%.2f\n", b.Balance)
}

func main() {
	fmt.Println("--- Bank Account System Demonstration ---")

	myAccount := BankAccount{
		Owner:   "Alice",
		Balance: 100.00, // Initial balance
	}
	fmt.Println("\n--- Initial State ---")
	myAccount.DisplayBalance()
	fmt.Println("\n--- Deposit Funds ---")
	myAccount.Deposit(50.50)
	myAccount.DisplayBalance() // Observe the change
	fmt.Println("\n--- Attempting Withdrawals ---")
	myAccount.Withdraw(25.00)
	myAccount.DisplayBalance()
	myAccount.Withdraw(200.00)
	myAccount.DisplayBalance()

	fmt.Println("\n--- Final State ---")
	myAccount.DisplayBalance()
	fmt.Println("\n--- Demonstrating Value Receiver vs. Pointer Receiver Effect ---")
	tempAccount := BankAccount{Owner: "Bob", Balance: 50.00}
	fmt.Printf("Bob's balance BEFORE TryToModifyBalance: $%.2f\n", tempAccount.Balance)
	TryToModifyBalance(tempAccount, 20.00) // Pass by value
	fmt.Printf("Bob's balance AFTER TryToModifyBalance: $%.2f (Did not change externally)\n", tempAccount.Balance)

	fmt.Printf("Alice's balance BEFORE Deposit (again): $%.2f\n", myAccount.Balance)
	myAccount.Deposit(10.00) // This actually modifies myAccount
	fmt.Printf("Alice's balance AFTER Deposit (again): $%.2f (Changed externally)\n", myAccount.Balance)
}
