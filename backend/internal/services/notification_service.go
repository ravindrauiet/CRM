package services

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type NotificationService struct {
	EmailService *EmailService
	db           *sql.DB
}

func NewNotificationService(db *sql.DB) *NotificationService {
	return &NotificationService{
		EmailService: NewEmailService(),
		db:           db,
	}
}

// NotifyStageCompletion sends email to admin when a stage is completed
func (ns *NotificationService) NotifyStageCompletion(jobID int, stage string, completedByUserID int) error {
	// Get job details
	job, err := ns.getJobDetails(jobID)
	if err != nil {
		log.Printf("Failed to get job details for notification: %v", err)
		return err
	}

	// Get user details who completed the stage
	completedByUser, err := ns.getUserDetails(completedByUserID)
	if err != nil {
		log.Printf("Failed to get user details for notification: %v", err)
		return err
	}

	// Get notification email (job-specific or admin email as fallback)
	notificationEmail, err := ns.getNotificationEmail(jobID)
	if err != nil {
		log.Printf("Failed to get notification email: %v", err)
		return err
	}

	// Determine next stage
	nextStage := ns.getNextStage(stage)
	stageName := ns.getStageName(stage)

	// Prepare email data
	emailData := StageCompletionEmail{
		JobNo:       job.JobNo,
		JobTitle:    fmt.Sprintf("Import/Export Job - %s", job.JobNo),
		Stage:       stage,
		StageName:   stageName,
		CompletedBy: completedByUser.Username,
		CompletedAt: time.Now().Format("2006-01-02 15:04:05"),
		NextStage:   nextStage,
		AdminEmail:  notificationEmail,
	}

	// Send email notification
	if err := ns.EmailService.SendStageCompletionEmail(emailData); err != nil {
		log.Printf("Failed to send stage completion email: %v", err)
		return err
	}

	// Log the notification
	log.Printf("Stage completion notification sent for job %s, stage %s, completed by %s to %s", 
		job.JobNo, stageName, completedByUser.Username, notificationEmail)

	return nil
}

// NotifyJobCreation sends email to admin when a new job is created
func (ns *NotificationService) NotifyJobCreation(jobID int, createdByUserID int) error {
	// Get job details
	job, err := ns.getJobDetails(jobID)
	if err != nil {
		log.Printf("Failed to get job details for creation notification: %v", err)
		return err
	}

	// Get user details who created the job
	createdByUser, err := ns.getUserDetails(createdByUserID)
	if err != nil {
		log.Printf("Failed to get user details for creation notification: %v", err)
		return err
	}

	// Get notification email (job-specific or admin email as fallback)
	notificationEmail, err := ns.getNotificationEmail(jobID)
	if err != nil {
		log.Printf("Failed to get notification email: %v", err)
		return err
	}

	// Send email notification
	if err := ns.EmailService.SendJobCreationEmail(job.JobNo, createdByUser.Username, notificationEmail); err != nil {
		log.Printf("Failed to send job creation email: %v", err)
		return err
	}

	// Log the notification
	log.Printf("Job creation notification sent for job %s, created by %s to %s", 
		job.JobNo, createdByUser.Username, notificationEmail)

	return nil
}

// Helper functions
func (ns *NotificationService) getJobDetails(jobID int) (*JobDetails, error) {
	query := `
		SELECT pj.id, pj.job_no, pj.current_stage, pj.status, pj.created_at,
		       s1.consignee, s1.shipper, s1.commodity, pj.notification_email
		FROM pipeline_jobs pj
		LEFT JOIN stage1_data s1 ON pj.id = s1.job_id
		WHERE pj.id = ?
	`
	
	var job JobDetails
	err := ns.db.QueryRow(query, jobID).Scan(
		&job.ID, &job.JobNo, &job.CurrentStage, &job.Status, &job.CreatedAt,
		&job.Consignee, &job.Shipper, &job.Commodity, &job.NotificationEmail,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &job, nil
}

func (ns *NotificationService) getUserDetails(userID int) (*UserDetails, error) {
	query := `SELECT id, username, designation, role FROM users WHERE id = ?`
	
	var user UserDetails
	err := ns.db.QueryRow(query, userID).Scan(
		&user.ID, &user.Username, &user.Designation, &user.Role,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

func (ns *NotificationService) getNotificationEmail(jobID int) (string, error) {
	// First, try to get the job-specific notification email from the job details
	job, err := ns.getJobDetails(jobID)
	if err != nil {
		log.Printf("Failed to get job details for notification email: %v", err)
		return "", err
	}

	if job.NotificationEmail.Valid {
		return job.NotificationEmail.String, nil
	}

	// If no job-specific email, return a default admin email
	// In a real system, you might want to store admin email in database or config
	adminEmail := "admin@maydiv.com" // Default admin email
	
	// You can also query the database for admin users and get their email
	// query := `SELECT email FROM users WHERE role = 'admin' LIMIT 1`
	// err := ns.db.QueryRow(query).Scan(&adminEmail)
	
	return adminEmail, nil
}

func (ns *NotificationService) getNextStage(currentStage string) string {
	switch currentStage {
	case "stage1":
		return "Stage 2 - Customs & Documentation"
	case "stage2":
		return "Stage 3 - Clearance & Logistics"
	case "stage3":
		return "Stage 4 - Billing & Completion"
	case "stage4":
		return "Completed"
	default:
		return "Unknown"
	}
}

func (ns *NotificationService) getStageName(stage string) string {
	switch stage {
	case "stage1":
		return "Stage 1 - Initial Setup"
	case "stage2":
		return "Stage 2 - Customs & Documentation"
	case "stage3":
		return "Stage 3 - Clearance & Logistics"
	case "stage4":
		return "Stage 4 - Billing & Completion"
	case "completed":
		return "Completed"
	default:
		return stage
	}
}

// Data structures for job and user details
type JobDetails struct {
	ID           int
	JobNo        string
	CurrentStage string
	Status       string
	CreatedAt    time.Time
	Consignee    sql.NullString
	Shipper      sql.NullString
	Commodity    sql.NullString
	NotificationEmail sql.NullString // Added for job-specific notification email
}

type UserDetails struct {
	ID          int
	Username    string
	Designation string
	Role        string
} 