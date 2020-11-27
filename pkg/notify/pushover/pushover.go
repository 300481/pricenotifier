package pushover

import (
	"github.com/gregdel/pushover"
)

// Pushover represents a notification client for Pushover
type Pushover struct {
	Token     string
	User      string
	App       *pushover.Pushover
	Recipient *pushover.Recipient
}

// NewPushover returns a client interface for a pushover client
func NewPushover(token, user string) *Pushover {
	app := pushover.New(token)
	recipient := pushover.NewRecipient(user)
	return &Pushover{
		Token:     token,
		User:      user,
		App:       app,
		Recipient: recipient,
	}
}

// Notify will send a notification via Pushover
// returns if the message was send without failure
func (p *Pushover) Notify(message string) bool {
	msg := pushover.NewMessage(message)
	_, err := p.App.SendMessage(msg, p.Recipient)
	return err == nil
}
