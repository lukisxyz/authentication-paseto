package email

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailVerificationBodyRequest struct {
	SUBJECT           string
	EMAIL             string
	VERIFICATION_CODE string
}

type EmailLinkForgotPasswordBodyRequest struct {
	SUBJECT string
	EMAIL   string
	CODE    string
}

// SendVerificationCode implements Mail.
func SendLinkForgotPassword(dest string, data EmailLinkForgotPasswordBodyRequest) {
	cwd, _ := os.Getwd()
	templateFile := filepath.Join(cwd, "/templates/email/change-password-link.html")

	res, err := parseTemplate(templateFile, data)
	if err != nil {
		log.Debug().Err(err).Msg("failed to parse email template")
	} else {
		sendMail(dest, res, data.SUBJECT)
	}
}

// SendVerificationCode implements Mail.
func SendVerificationCode(dest string, data EmailVerificationBodyRequest) {
	cwd, _ := os.Getwd()
	templateFile := filepath.Join(cwd, "/templates/email/verification_email.html")

	res, err := parseTemplate(templateFile, data)
	if err != nil {
		log.Debug().Err(err).Msg("failed to parse email template")
	} else {
		sendMail(dest, res, data.SUBJECT)
	}
}

func sendMail(dest, res, sbj string) {
	from := mail.NewEmail(os.Getenv("MAIL_SENDER_NAME"), os.Getenv("MAIL_SENDER_NAME"))
	to := mail.NewEmail(dest, dest)

	message := mail.NewSingleEmail(from, sbj, to, "", res)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	resp, err := client.Send(message)

	if err != nil {
		log.Debug().Err(err).Msg("failed to send email")
	} else if resp.StatusCode != 200 {
		log.Debug().Err(err).Msg("success send email")
	} else {
		log.Info().Msg("success send email")
	}
}

func parseTemplate(filePath string, data interface{}) (string, error) {
	t, err := template.ParseFiles(filePath)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), err
}
