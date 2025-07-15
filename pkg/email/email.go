package email

import (
	"brainloop-api/pkg/config"
	"fmt"

	"gopkg.in/mail.v2"
)

func SendVerificationEmail(toEmail, token string) error {
	verificationLink := fmt.Sprintf("%s/verify-email?token=%s", config.AppConfig.FrontendURL, token)
	htmlBody := `
	<!DOCTYPE html>
	<html>
	<head>
		<style>
			body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
			.container { max-width: 600px; margin: 20px auto; padding: 20px; border: 1px solid #ddd; border-radius: 10px; }
			.header { font-size: 24px; font-weight: bold; color: #0ea5e9; }
			.button { display: inline-block; padding: 12px 24px; margin: 20px 0; font-size: 16px; color: #fff; background-color: #0ea5e9; border-radius: 5px; text-decoration: none; }
			.footer { font-size: 12px; color: #777; margin-top: 20px; }
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">Welcome to BrainLoop!</div>
			<p>Thanks for signing up! We're excited to have you on board.</p>
			<p>Please click the button below to verify your email address and complete your registration:</p>
			<a href="` + verificationLink + `" class="button">Verify My Email</a>
			<p>If the button above doesn't work, you can also copy and paste this link into your browser:</p>
			<p><a href="` + verificationLink + `">` + verificationLink + `</a></p>
			<p>This link will expire in 24 hours.</p>
			<p class="footer">If you did not sign up for this account, you can safely ignore this email.</p>
		</div>
	</body>
	</html>
	`
	m := mail.NewMessage()
	m.SetHeader("From", "BrainLoop <noreply.brainloop@midaytech.com>")
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Verify Your Email for BrainLoop")
	m.SetBody("text/html", htmlBody)
	d := mail.NewDialer(config.AppConfig.SMTPHost, config.AppConfig.SMTPPort, config.AppConfig.SMTPUsername, config.AppConfig.SMTPPassword)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func SendPasswordResetEmail(toEmail, token string) error {
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", config.AppConfig.FrontendURL, token)
	htmlBody := `
	<!DOCTYPE html>
	<html>
	<head>
		<style>
			body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
			.container { max-width: 600px; margin: 20px auto; padding: 20px; border: 1px solid #ddd; border-radius: 10px; }
			.header { font-size: 24px; font-weight: bold; color: #0ea5e9; }
			.button { display: inline-block; padding: 12px 24px; margin: 20px 0; font-size: 16px; color: #fff; background-color: #0ea5e9; border-radius: 5px; text-decoration: none; }
			.footer { font-size: 12px; color: #777; margin-top: 20px; }
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">Reset Your BrainLoop Password</div>
			<p>We received a request to reset the password for your account.</p>
			<p>Please click the button below to set a new password. This link is only valid for 30 minutes.</p>
			<a href="` + resetLink + `" class="button">Reset My Password</a>
			<p>If the button above doesn't work, you can also copy and paste this link into your browser:</p>
			<p><a href="` + resetLink + `">` + resetLink + `</a></p>
			<p class="footer">If you did not request a password reset, you can safely ignore this email. Your password will not be changed.</p>
		</div>
	</body>
	</html>
	`
	m := mail.NewMessage()
	m.SetHeader("From", "BrainLoop <noreply.brainloop@midaytech.com>")
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Your BrainLoop Password Reset Request")
	m.SetBody("text/html", htmlBody)
	d := mail.NewDialer(config.AppConfig.SMTPHost, config.AppConfig.SMTPPort, config.AppConfig.SMTPUsername, config.AppConfig.SMTPPassword)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
