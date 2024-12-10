package infrastructure

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	emailtemplates "mini-wallet/infrastructure/email_templates"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const CONFIG_SMTP_HOST = "smtp.sendgrid.net"
const CONFIG_SMTP_PORT = 587
const CONFIG_SENDER_NAME = "Namulaki <corporation@namulaki.id>"
const CONFIG_AUTH_EMAIL = "apikey"
const CONFIG_AUTH_PASSWORD = "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"

func SendEmailVerificationLink(email string, userFullName string, token string, domain string) {

	from := mail.NewEmail("Namulaki", "corporation@namulaki.id")
	// subject := "Reset Password"

	to := mail.NewEmail(userFullName, email)

	content := mail.NewContent("text/html", emailtemplates.BuildVerifyEmailTemplate(token))
	m := mail.NewV3MailInit(from, "Verifikasi Akun "+domain, to, content)
	m.SetTemplateID(os.Getenv("SENDGRID_TEMPLATE_ID"))

	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	client := &rest.Client{
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				MaxIdleConns:          2,
				MaxIdleConnsPerHost:   2,
				IdleConnTimeout:       90 * time.Millisecond,
			},
			Timeout: 5 * time.Second,
		},
	}

	response, err := client.Send(request)
	if err != nil {
		log.Fatal("error sending email:", err.Error())
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		log.Fatal("an unexpected error occurred:", response.StatusCode, response.Body, response.Headers)
	}

}

func SendPasswordResetLink(email string, userFullName string, token string, domain string) {

	from := mail.NewEmail("Namulaki", "corporation@namulaki.id")
	// subject := "Reset Password"

	to := mail.NewEmail(userFullName, email)

	content := mail.NewContent("text/html", emailtemplates.BuildResetPasswordEmailTemplate(token))
	m := mail.NewV3MailInit(from, "Atur Ulang Kata Sandi "+domain, to, content)
	m.SetTemplateID(os.Getenv("SENDGRID_TEMPLATE_ID"))

	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	client := &rest.Client{
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				MaxIdleConns:          2,
				MaxIdleConnsPerHost:   2,
				IdleConnTimeout:       90 * time.Millisecond,
			},
			Timeout: 5 * time.Second,
		},
	}

	response, err := client.Send(request)
	if err != nil {
		fmt.Println("error sending email:", err.Error())
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		fmt.Println("an unexpected error occurred:", response.StatusCode, response.Body, response.Headers)
	}

}
