package main

import (
	"fmt"
)

// this can be left as unused
// but when var is defined and not used
// it will raise error
const TICKET_COUNT int = 50

func main() {
	var remaningTickets int = TICKET_COUNT
	var name string = "Cruel World"
	fmt.Printf("Welcome to %v, it suckes. But you're gonna love it!\n", name)
	fmt.Printf("Remaining tickets %v. Data Type: %T\n", remaningTickets, remaningTickets)
	fmt.Printf("Enter name:\n")
	var userName string
	fmt.Scan(&userName)
	var email string
	fmt.Printf("Enter email\n")
	fmt.Scan(&email)

	var ticketNumToBook int = 0
	fmt.Printf("Enter number of tickets to book\n")
	fmt.Scan(&ticketNumToBook)

	fmt.Printf("Ticket: %v, booked by %v, confirmation sent at %v\n", ticketNumToBook, name, email)

}
