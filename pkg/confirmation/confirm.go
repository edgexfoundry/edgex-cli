package confirmation

import (
	"fmt"
	"log"
	"strings"
)

const defaultConfirmMsg = "Are you sure? This cannot be undone: [y/n]"
const defaultAbortMsg = "Aborting"

type UserConfirmation struct {
	confirmMsg string
	abortMsg   string
}

//UserConfirmation creates new UserConfirmation using default values for confirmMsg and abortMsg
func New() UserConfirmation {
	return UserConfirmation{
		confirmMsg: defaultConfirmMsg,
		abortMsg:   defaultAbortMsg,
	}
}

//NewCustom creates new UserConfirmation using custom confirmMsg and abortMsg
func NewCustom(confirm, abort string) UserConfirmation {
	confirmMsg := defaultConfirmMsg
	if confirm != "" {
		confirmMsg = confirm
	}

	abortMsg := defaultAbortMsg
	if abort != "" {
		abortMsg = abort
	}
	return UserConfirmation{
		confirmMsg: confirmMsg,
		abortMsg:   abortMsg,
	}
}

// askForConfirmation uses Scanln to parse user input. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user. Typically, you should use fmt to print out a question
// before calling askForConfirmation. E.g. fmt.Println("WARNING: Are you sure? (yes/no)")
func askForConfirmation() bool {
	var response string

	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}

	switch strings.ToLower(response) {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		fmt.Println("I'm sorry but I didn't get what you meant, please type (y)es or (n)o and then press enter:")
		return askForConfirmation()
	}
}

//Ask user to confirm the execution
func (c UserConfirmation) Confirm() bool {
	fmt.Println(c.confirmMsg)
	if !askForConfirmation() {
		fmt.Println(c.abortMsg)
		return false
	}
	return true
}
