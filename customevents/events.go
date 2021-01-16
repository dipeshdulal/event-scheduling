package customevents

import "log"

// SendEmail sends the email
func SendEmail(data string) {
	log.Println("📨 Sending email with data: ", data)
}

// PayBills pays the bills
func PayBills(data string) {
	log.Println("💲 Pay me a bill: ", data)
}
