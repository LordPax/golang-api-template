package services

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"math/rand"
	"os"
	"time"

	sib_api_v3_sdk "github.com/getbrevo/brevo-go/lib"
)

type EmailTemplate struct {
	Name string
	Path string
}

var (
	WelcomeEmail       = EmailTemplate{Name: "welcome", Path: "template-html/welcome.html"}
	PasswordResetEmail = EmailTemplate{Name: "password_reset", Path: "template-html/password_reset.html"}
	VerificationEmail  = EmailTemplate{Name: "verification", Path: "template-html/verification.html"}
)

func SendEmail(userEmail string, template EmailTemplate) error {
	apiKey := os.Getenv("BREVO_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("BREVO_API_KEY environment variable is not set")
	}

	senderEmail := os.Getenv("BREVO_SENDER")
	if senderEmail == "" {
		return fmt.Errorf("BREVO_SENDER environment variable is not set")
	}

	cfg := sib_api_v3_sdk.NewConfiguration()
	cfg.AddDefaultHeader("api-key", apiKey)

	client := sib_api_v3_sdk.NewAPIClient(cfg)

	subject, htmlContent := getEmailContent(template, "")

	to := []sib_api_v3_sdk.SendSmtpEmailTo{
		{
			Email: userEmail,
		},
	}

	body := sib_api_v3_sdk.SendSmtpEmail{
		HtmlContent: htmlContent,
		Subject:     subject,
		Sender: &sib_api_v3_sdk.SendSmtpEmailSender{
			Name:  os.Getenv("BREVO_SENDER"),
			Email: senderEmail,
		},
		To: to,
		Params: map[string]interface{}{
			"subject": subject,
		},
	}

	ctx := context.Background()
	response, httpResp, err := client.TransactionalEmailsApi.SendTransacEmail(ctx, body)
	if err != nil {
		return fmt.Errorf("error sending email: %w, response: %v, http response: %v", err, response, httpResp)
	}

	return nil
}

func SendVerificationEmail(userEmail, code string) error {
	return SendEmailWithCode(userEmail, VerificationEmail, code)
}

func SendEmailWithCode(userEmail string, template EmailTemplate, code string) error {
	apiKey := os.Getenv("BREVO_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("BREVO_API_KEY environment variable is not set")
	}

	senderEmail := os.Getenv("BREVO_SENDER")
	if senderEmail == "" {
		return fmt.Errorf("BREVO_SENDER environment variable is not set")
	}

	cfg := sib_api_v3_sdk.NewConfiguration()
	cfg.AddDefaultHeader("api-key", apiKey)

	client := sib_api_v3_sdk.NewAPIClient(cfg)

	subject, htmlContent := getEmailContent(template, code)

	to := []sib_api_v3_sdk.SendSmtpEmailTo{
		{
			Email: userEmail,
		},
	}

	body := sib_api_v3_sdk.SendSmtpEmail{
		HtmlContent: htmlContent,
		Subject:     subject,
		Sender: &sib_api_v3_sdk.SendSmtpEmailSender{
			Name:  os.Getenv("BREVO_SENDER"),
			Email: senderEmail,
		},
		To: to,
		Params: map[string]interface{}{
			"subject": subject,
			"code":    code,
		},
	}

	ctx := context.Background()
	response, httpResp, err := client.TransactionalEmailsApi.SendTransacEmail(ctx, body)
	if err != nil {
		return fmt.Errorf("error sending email: %w, response: %v, http response: %v", err, response, httpResp)
	}

	return nil
}

func getEmailContent(template EmailTemplate, code string) (subject, htmlContent string) {
	params := map[string]interface{}{
		"Code": code,
	}

	htmlContent = loadHTMLTemplate(template.Path, params)
	switch template.Name {
	case WelcomeEmail.Name:
		subject = "Bienvenue sur notre plateforme!"
	case PasswordResetEmail.Name:
		subject = "Réinitialisation de votre mot de passe"
	case VerificationEmail.Name:
		subject = "Vérifiez votre compte"
	default:
		subject = "Hello from Golang"
		htmlContent = "<p>Hello world</p>"
	}
	return subject, htmlContent
}

func loadHTMLTemplate(filePath string, params map[string]interface{}) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error loading HTML template: %v\n", err)
		return "<p>Error loading template</p>"
	}

	tmpl := template.New("email")
	tmpl, err = tmpl.Parse(string(content))
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		return "<p>Error parsing template</p>"
	}

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, params); err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		return "<p>Error executing template</p>"
	}

	return buffer.String()
}

func GenerateVerificationCode() int {
	source := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(source)
	return 10000 + randomGenerator.Intn(90000)
}
