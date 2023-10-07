package mailer

type Defaults struct {
	Email string
	Name  string
}

// --------------------

type VerifyDTO struct {
	Defaults
	Code string
}

func (data VerifyDTO) Send() error {
	return Send("Verify your Twibber Account", "user_verify", data)
}
