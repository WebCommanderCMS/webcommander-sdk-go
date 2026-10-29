package email

import (
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/dto"
)

// EmailInput is the send body, bare and camelCase as the store takes it; html_content and other unknown keys are ignored.
type EmailInput struct {
	// Recipients, Cc and Bcc hold at most 100 addresses together.
	Recipients []string
	Cc         []string
	Bcc        []string
	Subject    string
	// Send HTMLContent or a plugin template's Identifier, not both.
	HTMLContent string `wire:"htmlContent"`
	ContentType string `wire:"contentType"`
	Identifier  string
	// Data is the template's variables, an object of any shape.
	Data any
}

// ToMap renders the request body.
func (i EmailInput) ToMap() map[string]any { return dto.ToMap(i) }

// EmailRejection is one address the store could not send to, and why, such as send_failed.
type EmailRejection struct {
	Email  *string
	Reason *string
}

// EmailSendResult is the send's verdict: each list holds the addresses trimmed, a repeat in any case once.
type EmailSendResult struct {
	Status   *string
	Accepted []*string
	Rejected []*EmailRejection
	Cc       []*string
	Bcc      []*string
}

// EmailSendResultShape is the typed view of the send's body.
func EmailSendResultShape(decoded any) any { return dto.Read[EmailSendResult](decoded) }
