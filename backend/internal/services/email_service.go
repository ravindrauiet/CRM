package services

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	dialer *gomail.Dialer
	from   string
}

type StageCompletionEmail struct {
	JobNo        string
	JobTitle     string
	Stage        string
	StageName    string
	CompletedBy  string
	CompletedAt  string
	NextStage    string
	AdminEmail   string
}

func NewEmailService() *EmailService {
	// Get SMTP configuration from environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	fromEmail := os.Getenv("FROM_EMAIL")

	if smtpHost == "" {
		smtpHost = "smtp.gmail.com" // Default to Gmail
	}

	smtpPort := 587 // Default port
	if smtpPortStr != "" {
		if port, err := strconv.Atoi(smtpPortStr); err == nil {
			smtpPort = port
		}
	}

	if fromEmail == "" {
		fromEmail = smtpUser // Use SMTP user as from email if not specified
	}

	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	return &EmailService{
		dialer: dialer,
		from:   fromEmail,
	}
}

func (es *EmailService) SendStageCompletionEmail(emailData StageCompletionEmail) error {
	subject := fmt.Sprintf("Stage Completion Notification - Job %s", emailData.JobNo)
	
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Stage Completion Notification</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #2563eb; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background-color: #f8fafc; padding: 20px; border-radius: 0 0 8px 8px; }
        .info-row { margin: 10px 0; }
        .label { font-weight: bold; color: #374151; }
        .value { color: #1f2937; }
        .stage-badge { 
            display: inline-block; 
            padding: 4px 12px; 
            background-color: #10b981; 
            color: white; 
            border-radius: 20px; 
            font-size: 14px; 
            font-weight: bold; 
        }
        .footer { margin-top: 20px; padding-top: 20px; border-top: 1px solid #e5e7eb; font-size: 12px; color: #6b7280; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎉 Stage Completion Notification</h1>
        </div>
        <div class="content">
            <p>Hello Admin,</p>
            
            <p>A pipeline stage has been successfully completed. Here are the details:</p>
            
            <div class="info-row">
                <span class="label">Job Number:</span>
                <span class="value">%s</span>
            </div>
            
            <div class="info-row">
                <span class="label">Job Title:</span>
                <span class="value">%s</span>
            </div>
            
            <div class="info-row">
                <span class="label">Completed Stage:</span>
                <span class="stage-badge">%s</span>
            </div>
            
            <div class="info-row">
                <span class="label">Completed By:</span>
                <span class="value">%s</span>
            </div>
            
            <div class="info-row">
                <span class="label">Completion Time:</span>
                <span class="value">%s</span>
            </div>
            
            <div class="info-row">
                <span class="label">Next Stage:</span>
                <span class="value">%s</span>
            </div>
            
            <p style="margin-top: 20px;">
                <strong>Action Required:</strong> Please review the completed stage and proceed with the next stage if everything is in order.
            </p>
            
            <div class="footer">
                <p>This is an automated notification from the MayDiv CRM System.</p>
                <p>If you have any questions, please contact the system administrator.</p>
            </div>
        </div>
    </div>
</body>
</html>
`, emailData.JobNo, emailData.JobTitle, emailData.StageName, emailData.CompletedBy, emailData.CompletedAt, emailData.NextStage)

	m := gomail.NewMessage()
	m.SetHeader("From", es.from)
	m.SetHeader("To", emailData.AdminEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Try to send email
	if err := es.dialer.DialAndSend(m); err != nil {
		log.Printf("Failed to send email notification: %v", err)
		return err
	}

	log.Printf("Stage completion email sent successfully to %s for job %s", emailData.AdminEmail, emailData.JobNo)
	return nil
}

func (es *EmailService) SendJobCreationEmail(jobNo, createdBy, adminEmail string) error {
	subject := fmt.Sprintf("New Pipeline Job Created - %s", jobNo)
	
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>New Job Creation Notification</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #059669; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background-color: #f8fafc; padding: 20px; border-radius: 0 0 8px 8px; }
        .info-row { margin: 10px 0; }
        .label { font-weight: bold; color: #374151; }
        .value { color: #1f2937; }
        .footer { margin-top: 20px; padding-top: 20px; border-top: 1px solid #e5e7eb; font-size: 12px; color: #6b7280; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🆕 New Pipeline Job Created</h1>
        </div>
        <div class="content">
            <p>Hello Admin,</p>
            
            <p>A new pipeline job has been created in the system. Here are the details:</p>
            
            <div class="info-row">
                <span class="label">Job Number:</span>
                <span class="value">%s</span>
            </div>
            
            <div class="info-row">
                <span class="label">Created By:</span>
                <span class="value">%s</span>
            </div>
            
            <div class="info-row">
                <span class="label">Status:</span>
                <span class="value">Stage 1 - Initial Setup</span>
            </div>
            
            <p style="margin-top: 20px;">
                <strong>Action Required:</strong> Please review the new job and assign it to the appropriate team members.
            </p>
            
            <div class="footer">
                <p>This is an automated notification from the MayDiv CRM System.</p>
                <p>If you have any questions, please contact the system administrator.</p>
            </div>
        </div>
    </div>
</body>
</html>
`, jobNo, createdBy)

	m := gomail.NewMessage()
	m.SetHeader("From", es.from)
	m.SetHeader("To", adminEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	if err := es.dialer.DialAndSend(m); err != nil {
		log.Printf("Failed to send job creation email: %v", err)
		return err
	}

	log.Printf("Job creation email sent successfully to %s for job %s", adminEmail, jobNo)
	return nil
}

// Test email configuration
func (es *EmailService) TestEmailConnection() error {
	// Try to connect to SMTP server
	s, err := es.dialer.Dial()
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %v", err)
	}
	defer s.Close()
	
	log.Println("Email service connection test successful")
	return nil
} 