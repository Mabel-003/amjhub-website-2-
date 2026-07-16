package backend

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"os"
	"strings"
	"time"
)

type EmailConfig struct {
	SMTPHost     string // e.g. smtp.gmail.com, smtp.office365.com, smtp.zoho.com
	SMTPPort     string // e.g. 587
	SMTPUsername string // sender account used to authenticate with the SMTP server
	SMTPPassword string // app password / SMTP password for the sender account
	FromAddress  string // "From" address shown on the email (often same as SMTPUsername)
	ToAddress    string // recipient inbox — THE CLIENT'S BUSINESS EMAIL
}

// LoadEmailConfig reads all email settings from environment variables.
//
// Required environment variables:
//
//	SMTP_HOST       - SMTP server hostname (e.g. "smtp.gmail.com")
//	SMTP_PORT       - SMTP server port (e.g. "587")
//	SMTP_USERNAME   - account used to log in to the SMTP server
//	SMTP_PASSWORD   - password or app-specific password for that account
//	EMAIL_FROM      - the "From" address on outgoing mail (optional, falls back to SMTP_USERNAME)
//	CONTACT_RECIPIENT_EMAIL - the inbox that should receive form submissions (the client's email)
func LoadEmailConfig() EmailConfig {
	from := os.Getenv("EMAIL_FROM")
	if from == "" {
		from = os.Getenv("SMTP_USERNAME")
	}

	return EmailConfig{
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		FromAddress:  from,
		// 👇 This is the key line for your client: the recipient inbox is
		// read entirely from the CONTACT_RECIPIENT_EMAIL environment variable.
		// They can change where inquiries get delivered without redeploying code.
		ToAddress: os.Getenv("CONTACT_RECIPIENT_EMAIL"),
	}
}

// IsSMTPConfigured reports whether the SMTP credentials and sender address
// are available. This is the baseline required for any outgoing email.
func (c EmailConfig) IsSMTPConfigured() bool {
	return c.SMTPHost != "" &&
		c.SMTPPort != "" &&
		c.SMTPUsername != "" &&
		c.SMTPPassword != "" &&
		c.FromAddress != ""
}

// IsBusinessConfigured reports whether the business notification recipient
// is also configured.
func (c EmailConfig) IsBusinessConfigured() bool {
	return c.IsSMTPConfigured() && c.ToAddress != ""
}

// SendBusinessNotification emails the contact form submission to the
// configured business recipient using the existing SMTP credentials.
func SendBusinessNotification(cfg EmailConfig, data ContactFormData) error {
	if !cfg.IsBusinessConfigured() {
		return fmt.Errorf("email not sent: SMTP_HOST, SMTP_PORT, SMTP_USERNAME, SMTP_PASSWORD, EMAIL_FROM, or CONTACT_RECIPIENT_EMAIL is not set")
	}

	subject := fmt.Sprintf("New Inquiry from %s — AMJ HUB Website", data.FullName)
	body := buildEmailBody(data)

	// Include Reply-To so the business can hit reply directly to the user.
	headers := map[string]string{"Reply-To": data.Email}
	return sendEmail(cfg, cfg.ToAddress, subject, body, "", headers)
}

// SendConfirmationEmail sends an acknowledgement email to the user after the
// business notification has been successfully delivered.
func SendConfirmationEmail(cfg EmailConfig, data ContactFormData) error {
	if !cfg.IsSMTPConfigured() {
		return fmt.Errorf("email not sent: SMTP_HOST, SMTP_PORT, SMTP_USERNAME, SMTP_PASSWORD, or EMAIL_FROM is not set")
	}

	subject := "We've Received Your Message | AMJ HUB"
	text := fmt.Sprintf("Hi %s,\n\nThank you for contacting AMJ HUB.\n\nWe've successfully received your message and appreciate your interest in our photography, videography, and drone services.\n\nOur team is currently reviewing your enquiry and will get back to you as soon as possible, usually within 24 hours.\n\nIf your request is urgent, feel free to contact us directly.\n\nWe look forward to bringing your vision to life.\n\nBest regards,\n\nAMJ HUB\nPhotography | Videography | Drone Services\n", data.FullName)
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>We've Received Your Message | AMJ HUB</title>
</head>
<body>
<p>Hi %s,</p>
<p>Thank you for contacting <strong>AMJ HUB</strong>.</p>
<p>We've successfully received your message and appreciate your interest in our photography, videography, and drone services.</p>
<p>Our team is currently reviewing your enquiry and will get back to you as soon as possible, usually within 24 hours.</p>
<p>If your request is urgent, feel free to contact us directly.</p>
<p>We look forward to bringing your vision to life.</p>
<p>Best regards,<br>
<strong>AMJ HUB</strong><br>
Photography | Videography | Drone Services</p>
</body>
</html>`, data.FullName)

	return sendEmail(cfg, data.Email, subject, text, html, nil)
}

// sendEmail is a reusable helper that sends an email via SMTP using the
// configured SMTP credentials. It supports text-only or HTML+text messages.
func sendEmail(cfg EmailConfig, to, subject, textBody, htmlBody string, extra map[string]string) error {
	addr := cfg.SMTPHost + ":" + cfg.SMTPPort
	auth := smtp.PlainAuth("", cfg.SMTPUsername, cfg.SMTPPassword, cfg.SMTPHost)
	msg, err := buildMIMEMessage(cfg.FromAddress, to, subject, textBody, htmlBody, extra)
	if err != nil {
		return err
	}

	if err := smtp.SendMail(addr, auth, cfg.FromAddress, []string{to}, msg); err != nil {
		return fmt.Errorf("smtp send failed: %w", err)
	}
	return nil
}

// buildMIMEMessage assembles a valid RFC 822 / MIME email with headers and
// either a plain-text body or a multipart/alternative body when HTML is
// provided.
func buildMIMEMessage(from, to, subject, textBody, htmlBody string, extra map[string]string) ([]byte, error) {
	headers := map[string]string{
		"From":         from,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
	}
	for k, v := range extra {
		headers[k] = v
	}

	var buf bytes.Buffer
	if htmlBody == "" {
		headers["Content-Type"] = `text/plain; charset="UTF-8"`
		for k, v := range headers {
			buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
		}
		buf.WriteString("\r\n")
		buf.WriteString(textBody)
		return buf.Bytes(), nil
	}

	boundary := fmt.Sprintf("boundary_%d", time.Now().UnixNano())
	headers["Content-Type"] = fmt.Sprintf("multipart/alternative; boundary=%s", boundary)
	for k, v := range headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	buf.WriteString("\r\n")

	writer := multipart.NewWriter(&buf)
	if err := writer.SetBoundary(boundary); err != nil {
		return nil, fmt.Errorf("failed to set MIME boundary: %w", err)
	}

	textPartHeader := make(textproto.MIMEHeader)
	textPartHeader.Set("Content-Type", "text/plain; charset=\"UTF-8\"")
	textPart, err := writer.CreatePart(textPartHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to create text part: %w", err)
	}
	textPart.Write([]byte(textBody))

	htmlPartHeader := make(textproto.MIMEHeader)
	htmlPartHeader.Set("Content-Type", "text/html; charset=\"UTF-8\"")
	htmlPart, err := writer.CreatePart(htmlPartHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to create html part: %w", err)
	}
	htmlPart.Write([]byte(htmlBody))

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close MIME writer: %w", err)
	}

	return buf.Bytes(), nil
}

// buildEmailBody formats the contact form data into a readable plain-text
// email body for the recipient.
func buildEmailBody(data ContactFormData) string {
	var sb strings.Builder

	sb.WriteString("You have a new inquiry from the AMJ HUB website.\n")
	sb.WriteString("─────────────────────────────────────────────\n\n")
	sb.WriteString(fmt.Sprintf("Full Name : %s\n", data.FullName))
	sb.WriteString(fmt.Sprintf("Email     : %s\n", data.Email))
	sb.WriteString(fmt.Sprintf("Phone     : %s\n", ifEmpty(data.Phone, "(not provided)")))
	sb.WriteString(fmt.Sprintf("Service   : %s\n", ifEmpty(data.Service, "(not specified)")))
	sb.WriteString(fmt.Sprintf("Submitted : %s\n\n", data.Timestamp.Format("Monday, 2 January 2006 at 15:04 MST")))
	sb.WriteString("Message:\n")
	sb.WriteString(data.Message)
	sb.WriteString("\n\n─────────────────────────────────────────────\n")
	sb.WriteString("Reply directly to this email to respond to the client.\n")

	return sb.String()
}

// logEmailFallback writes the inquiry to the server log when email isn't
// configured, so no submission is silently lost during setup/testing.
func logEmailFallback(data ContactFormData, reason error) {
	log.Println("⚠️  EMAIL NOT SENT — falling back to console log only.")
	log.Printf("    Reason: %v", reason)
	log.Println("    To enable email delivery, set the following environment variables:")
	log.Println("      SMTP_HOST, SMTP_PORT, SMTP_USERNAME, SMTP_PASSWORD, CONTACT_RECIPIENT_EMAIL")
}
