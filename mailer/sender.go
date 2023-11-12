package mailer

// Defaults struct holds the common fields required for sending emails.
type Defaults struct {
	Email string // The recipient's email address
	Name  string // The recipient's name
}

type VerifyDTO struct {
	Defaults
	Code string
}

func (data VerifyDTO) Send() error {
	return Send("Verify your Twibber Account", "user_verify", data)
}
