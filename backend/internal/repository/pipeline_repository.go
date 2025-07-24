package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"maydiv-crm/internal/models"
)

type PipelineRepository struct {
	db *sql.DB
}

func NewPipelineRepository(db *sql.DB) *PipelineRepository {
	return &PipelineRepository{db: db}
}

// GetAllJobs retrieves all pipeline jobs with their current stage data
func (r *PipelineRepository) GetAllJobs() ([]models.PipelineJobResponse, error) {
	query := `
		SELECT 
			pj.id, pj.job_no, pj.current_stage, pj.status, pj.created_by, 
			pj.assigned_to_stage2, pj.assigned_to_stage3, pj.customer_id,
			pj.created_at, pj.updated_at,
			u1.username as created_by_user,
			u2.username as stage2_user_name,
			u3.username as stage3_user_name,
			u4.username as customer_name
		FROM pipeline_jobs pj
		LEFT JOIN users u1 ON pj.created_by = u1.id
		LEFT JOIN users u2 ON pj.assigned_to_stage2 = u2.id
		LEFT JOIN users u3 ON pj.assigned_to_stage3 = u3.id
		LEFT JOIN users u4 ON pj.customer_id = u4.id
		ORDER BY pj.created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []models.PipelineJobResponse
	for rows.Next() {
		var job models.PipelineJobResponse
		var stage2UserName, stage3UserName, customerName sql.NullString

		err := rows.Scan(
			&job.ID, &job.JobNo, &job.CurrentStage, &job.Status, &job.CreatedBy,
			&job.AssignedToStage2, &job.AssignedToStage3, &job.CustomerID,
			&job.CreatedAt, &job.UpdatedAt,
			&job.CreatedByUser, &stage2UserName, &stage3UserName, &customerName,
		)
		if err != nil {
			return nil, err
		}

		if stage2UserName.Valid {
			job.Stage2UserName = stage2UserName.String
		}
		if stage3UserName.Valid {
			job.Stage3UserName = stage3UserName.String
		}
		if customerName.Valid {
			job.CustomerName = customerName.String
		}

		// Load stage data based on current stage
		if err := r.loadJobStageData(&job); err != nil {
			log.Printf("Error loading stage data for job %d: %v", job.ID, err)
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

// GetJobByID retrieves a specific job with all its stage data
func (r *PipelineRepository) GetJobByID(jobID int) (*models.PipelineJobResponse, error) {
	log.Printf("GetJobByID called with jobID: %d", jobID)
	
	query := `
		SELECT 
			pj.id, pj.job_no, pj.current_stage, pj.status, pj.created_by, 
			pj.assigned_to_stage2, pj.assigned_to_stage3, pj.customer_id,
			pj.created_at, pj.updated_at,
			u1.username as created_by_user,
			u2.username as stage2_user_name,
			u3.username as stage3_user_name,
			u4.username as customer_name
		FROM pipeline_jobs pj
		LEFT JOIN users u1 ON pj.created_by = u1.id
		LEFT JOIN users u2 ON pj.assigned_to_stage2 = u2.id
		LEFT JOIN users u3 ON pj.assigned_to_stage3 = u3.id
		LEFT JOIN users u4 ON pj.customer_id = u4.id
		WHERE pj.id = ?
	`

	var job models.PipelineJobResponse
	var stage2UserName, stage3UserName, customerName sql.NullString

	log.Printf("Executing query for job ID: %d", jobID)
	err := r.db.QueryRow(query, jobID).Scan(
		&job.ID, &job.JobNo, &job.CurrentStage, &job.Status, &job.CreatedBy,
		&job.AssignedToStage2, &job.AssignedToStage3, &job.CustomerID,
		&job.CreatedAt, &job.UpdatedAt,
		&job.CreatedByUser, &stage2UserName, &stage3UserName, &customerName,
	)
	if err != nil {
		log.Printf("Error scanning job: %v", err)
		return nil, err
	}
	log.Printf("Job found: ID=%d, JobNo=%s, Stage=%s", job.ID, job.JobNo, job.CurrentStage)

	if stage2UserName.Valid {
		job.Stage2UserName = stage2UserName.String
	}
	if stage3UserName.Valid {
		job.Stage3UserName = stage3UserName.String
	}
	if customerName.Valid {
		job.CustomerName = customerName.String
	}

	// Load all stage data
	if err := r.loadJobStageData(&job); err != nil {
		return nil, err
	}

	// Load job updates
	updates, err := r.getJobUpdates(jobID)
	if err == nil {
		job.Updates = updates
	}

	return &job, nil
}

// GetJobsByUserRole retrieves jobs assigned to a specific user based on their role
func (r *PipelineRepository) GetJobsByUserRole(userID int, role string) ([]models.PipelineJobResponse, error) {
	log.Printf("GetJobsByUserRole called with userID=%d, role=%s", userID, role)
	
	var query string
	switch role {
	case "subadmin":
		// Subadmin sees all jobs (same as admin)
		query = `
			SELECT 
				pj.id, pj.job_no, pj.current_stage, pj.status, pj.created_by, 
				pj.assigned_to_stage2, pj.assigned_to_stage3, pj.customer_id,
				pj.created_at, pj.updated_at,
				u1.username as created_by_user,
				u2.username as stage2_user_name,
				u3.username as stage3_user_name,
				u4.username as customer_name
			FROM pipeline_jobs pj
			LEFT JOIN users u1 ON pj.created_by = u1.id
			LEFT JOIN users u2 ON pj.assigned_to_stage2 = u2.id
			LEFT JOIN users u3 ON pj.assigned_to_stage3 = u3.id
			LEFT JOIN users u4 ON pj.customer_id = u4.id
			ORDER BY pj.created_at DESC
		`
		log.Printf("Subadmin query: %s", query)
	case "stage1_employee":
		query = `
			SELECT 
				pj.id, pj.job_no, pj.current_stage, pj.status, pj.created_by, 
				pj.assigned_to_stage2, pj.assigned_to_stage3, pj.customer_id,
				pj.created_at, pj.updated_at,
				u1.username as created_by_user
			FROM pipeline_jobs pj
			LEFT JOIN users u1 ON pj.created_by = u1.id
			WHERE pj.created_by = ? AND pj.current_stage IN ('stage1', 'stage2', 'stage3', 'stage4')
			ORDER BY pj.created_at DESC
		`
		log.Printf("Stage1 query: %s", query)
		log.Printf("Querying with userID: %d", userID)
	case "stage2_employee":
		query = `
			SELECT 
				pj.id, pj.job_no, pj.current_stage, pj.status, pj.created_by, 
				pj.assigned_to_stage2, pj.assigned_to_stage3, pj.customer_id,
				pj.created_at, pj.updated_at,
				u1.username as created_by_user
			FROM pipeline_jobs pj
			LEFT JOIN users u1 ON pj.created_by = u1.id
			WHERE pj.assigned_to_stage2 = ? AND pj.current_stage IN ('stage1', 'stage2', 'stage3', 'stage4')
			ORDER BY pj.created_at DESC
		`
		log.Printf("Stage2 query: %s", query)
		log.Printf("Querying with userID: %d", userID)
	case "stage3_employee":
		query = `
			SELECT 
				pj.id, pj.job_no, pj.current_stage, pj.status, pj.created_by, 
				pj.assigned_to_stage2, pj.assigned_to_stage3, pj.customer_id,
				pj.created_at, pj.updated_at,
				u1.username as created_by_user
			FROM pipeline_jobs pj
			LEFT JOIN users u1 ON pj.created_by = u1.id
			WHERE pj.assigned_to_stage3 = ? AND pj.current_stage IN ('stage2', 'stage3', 'stage4')
			ORDER BY pj.created_at DESC
		`
	case "customer":
		query = `
			SELECT 
				pj.id, pj.job_no, pj.current_stage, pj.status, pj.created_by, 
				pj.assigned_to_stage2, pj.assigned_to_stage3, pj.customer_id,
				pj.created_at, pj.updated_at,
				u1.username as created_by_user
			FROM pipeline_jobs pj
			LEFT JOIN users u1 ON pj.created_by = u1.id
			WHERE pj.customer_id = ? AND pj.current_stage IN ('stage3', 'stage4')
			ORDER BY pj.created_at DESC
		`
	default:
		return nil, fmt.Errorf("invalid role: %s", role)
	}

	// For subadmin, no userID parameter is needed since they see all jobs
	var rows *sql.Rows
	var err error
	if role == "subadmin" {
		rows, err = r.db.Query(query)
	} else {
		rows, err = r.db.Query(query, userID)
	}
	
	if err != nil {
		log.Printf("Query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	log.Printf("Query executed successfully, checking rows...")
	
	// First, let's check if there are any jobs at all
	var count int
	err = r.db.QueryRow("SELECT COUNT(*) FROM pipeline_jobs").Scan(&count)
	if err != nil {
		log.Printf("Error counting jobs: %v", err)
	} else {
		log.Printf("Total jobs in database: %d", count)
	}
	
	// Check if there are any jobs assigned to this user
	var assignedCount int
	err = r.db.QueryRow("SELECT COUNT(*) FROM pipeline_jobs WHERE assigned_to_stage2 = ?", userID).Scan(&assignedCount)
	if err != nil {
		log.Printf("Error counting assigned jobs: %v", err)
	} else {
		log.Printf("Jobs assigned to user %d: %d", userID, assignedCount)
	}
	
	// Check if there are any jobs created by this user (for stage1 employees)
	var createdCount int
	err = r.db.QueryRow("SELECT COUNT(*) FROM pipeline_jobs WHERE created_by = ?", userID).Scan(&createdCount)
	if err != nil {
		log.Printf("Error counting created jobs: %v", err)
	} else {
		log.Printf("Jobs created by user %d: %d", userID, createdCount)
	}
	
	var jobs []models.PipelineJobResponse
	for rows.Next() {
		var job models.PipelineJobResponse
		
		if role == "subadmin" {
			// Subadmin query includes all user names
			var stage2UserName, stage3UserName, customerName sql.NullString
			err := rows.Scan(
				&job.ID, &job.JobNo, &job.CurrentStage, &job.Status, &job.CreatedBy,
				&job.AssignedToStage2, &job.AssignedToStage3, &job.CustomerID,
				&job.CreatedAt, &job.UpdatedAt, &job.CreatedByUser,
				&stage2UserName, &stage3UserName, &customerName,
			)
			if err != nil {
				return nil, err
			}
			
			if stage2UserName.Valid {
				job.Stage2UserName = stage2UserName.String
			}
			if stage3UserName.Valid {
				job.Stage3UserName = stage3UserName.String
			}
			if customerName.Valid {
				job.CustomerName = customerName.String
			}
		} else {
			// Employee queries only include created_by_user
			err := rows.Scan(
				&job.ID, &job.JobNo, &job.CurrentStage, &job.Status, &job.CreatedBy,
				&job.AssignedToStage2, &job.AssignedToStage3, &job.CustomerID,
				&job.CreatedAt, &job.UpdatedAt, &job.CreatedByUser,
			)
			if err != nil {
				return nil, err
			}
		}

		// Load stage data based on current stage
		if err := r.loadJobStageData(&job); err != nil {
			log.Printf("Error loading stage data for job %d: %v", job.ID, err)
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

// CreateJob creates a new pipeline job with stage 1 data
func (r *PipelineRepository) CreateJob(req *models.Stage1CreateRequest, createdBy int) (*models.PipelineJobResponse, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create pipeline job
	jobResult, err := tx.Exec(`
		INSERT INTO pipeline_jobs (job_no, current_stage, status, created_by, assigned_to_stage2, assigned_to_stage3, customer_id)
		VALUES (?, 'stage1', 'active', ?, ?, ?, ?)
	`, req.JobNo, createdBy, nullInt(req.AssignedToStage2), nullInt(req.AssignedToStage3), nullInt(req.CustomerID))
	if err != nil {
		return nil, err
	}

	jobID, err := jobResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Create stage1 data
	_, err = tx.Exec(`
		INSERT INTO stage1_data (
			job_id, job_no, job_date, edi_job_no, edi_date, consignee, shipper,
			port_of_discharge, final_place_of_delivery, port_of_loading, country_of_shipment,
			hbl_no, hbl_date, mbl_no, mbl_date, shipping_line, forwarder,
			weight, packages, invoice_no, invoice_date, gateway_igm, gateway_igm_date,
			local_igm, local_igm_date, commodity, eta, current_status,
			container_no, container_size, date_of_arrival
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		jobID, req.JobNo, parseDate(req.JobDate), req.EDIJobNo, parseDate(req.EDIDate),
		req.Consignee, req.Shipper, req.PortOfDischarge, req.FinalPlaceOfDelivery,
		req.PortOfLoading, req.CountryOfShipment, req.HBLNo, parseDate(req.HBLDate),
		req.MBLNo, parseDate(req.MBLDate), req.ShippingLine, req.Forwarder,
		req.Weight, req.Packages, req.InvoiceNo, parseDate(req.InvoiceDate),
		req.GatewayIGM, parseDate(req.GatewayIGMDate), req.LocalIGM, parseDate(req.LocalIGMDate),
		req.Commodity, parseDateTime(req.ETA), req.CurrentStatus,
		req.ContainerNo, req.ContainerSize, parseDate(req.DateOfArrival),
	)
	if err != nil {
		return nil, err
	}

	// Add job update
	_, err = tx.Exec(`
		INSERT INTO job_updates (job_id, user_id, stage, update_type, message)
		VALUES (?, ?, 'stage1', 'status_change', 'Job created')
	`, jobID, createdBy)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Return the created job
	return r.GetJobByID(int(jobID))
}

// UpdateStage2Data updates stage 2 data and advances job to stage 2
func (r *PipelineRepository) UpdateStage2Data(jobID int, req *models.Stage2UpdateRequest, userID int) error {
	log.Printf("UpdateStage2Data called with jobID: %d, userID: %d", jobID, userID)
	log.Printf("Stage2 data: %+v", req)
	
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert or update stage2 data
	log.Printf("Executing stage2 data insert/update for job %d", jobID)
	log.Printf("HSNCode: %s", req.HSNCode)
	log.Printf("FilingRequirement: %s", req.FilingRequirement)
	log.Printf("DutyAmount: %f", req.DutyAmount)
	log.Printf("OceanFreight: %f", req.OceanFreight)
	
	stage2Result, err := tx.Exec(`
		INSERT INTO stage2_data (
			job_id, hsn_code, filing_requirement, checklist_sent_date, approval_date,
			bill_of_entry_no, bill_of_entry_date, debit_note, debit_paid_by,
			duty_amount, duty_paid_by, ocean_freight, destination_charges,
			original_doct_recd_date, drn_no, irn_no, documents_type
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			hsn_code = VALUES(hsn_code),
			filing_requirement = VALUES(filing_requirement),
			checklist_sent_date = VALUES(checklist_sent_date),
			approval_date = VALUES(approval_date),
			bill_of_entry_no = VALUES(bill_of_entry_no),
			bill_of_entry_date = VALUES(bill_of_entry_date),
			debit_note = VALUES(debit_note),
			debit_paid_by = VALUES(debit_paid_by),
			duty_amount = VALUES(duty_amount),
			duty_paid_by = VALUES(duty_paid_by),
			ocean_freight = VALUES(ocean_freight),
			destination_charges = VALUES(destination_charges),
			original_doct_recd_date = VALUES(original_doct_recd_date),
			drn_no = VALUES(drn_no),
			irn_no = VALUES(irn_no),
			documents_type = VALUES(documents_type),
			updated_at = CURRENT_TIMESTAMP
	`,
		jobID, req.HSNCode, req.FilingRequirement, parseDate(req.ChecklistSentDate),
		parseDate(req.ApprovalDate), req.BillOfEntryNo, parseDate(req.BillOfEntryDate),
		req.DebitNote, req.DebitPaidBy, req.DutyAmount, req.DutyPaidBy,
		req.OceanFreight, req.DestinationCharges, parseDate(req.OriginalDoctRecdDate),
		req.DRNNo, req.IRNNo, req.DocumentsType,
	)
	
	if err != nil {
		log.Printf("Error executing stage2 data insert/update: %v", err)
		return err
	}
	
	stage2RowsAffected, _ := stage2Result.RowsAffected()
	log.Printf("Stage2 data insert/update affected %d rows", stage2RowsAffected)
	if err != nil {
		return err
	}

	// Update job stage if not already in stage2 or beyond
	log.Printf("Updating job stage to stage2 for job %d", jobID)
	jobResult, err := tx.Exec(`
		UPDATE pipeline_jobs 
		SET current_stage = 'stage2', updated_at = CURRENT_TIMESTAMP 
		WHERE id = ? AND current_stage = 'stage1'
	`, jobID)
	if err != nil {
		log.Printf("Error updating job stage: %v", err)
		return err
	}
	
	jobRowsAffected, _ := jobResult.RowsAffected()
	log.Printf("Job stage update affected %d rows", jobRowsAffected)

	// Add job update
	_, err = tx.Exec(`
		INSERT INTO job_updates (job_id, user_id, stage, update_type, message)
		VALUES (?, ?, 'stage2', 'data_update', 'Stage 2 data updated')
	`, jobID, userID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	log.Printf("Stage2 data updated successfully for job %d", jobID)
	return nil
}

// UpdateStage3Data updates stage 3 data and advances job to stage 3
func (r *PipelineRepository) UpdateStage3Data(jobID int, req *models.Stage3UpdateRequest, userID int) error {
	log.Printf("UpdateStage3Data called with jobID: %d, userID: %d", jobID, userID)
	log.Printf("Stage3 data: %+v", req)
	
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert or update stage3 data
	log.Printf("Executing stage3 data insert/update for job %d", jobID)
	log.Printf("Custodian: %s", req.Custodian)
	log.Printf("DispatchInfo: %s", req.DispatchInfo)
	log.Printf("ClearanceExps: %f", req.ClearanceExps)
	log.Printf("StampDuty: %f", req.StampDuty)
	
	stage3Result, err := tx.Exec(`
		INSERT INTO stage3_data (
			job_id, exam_date, out_of_charge, clearance_exps, stamp_duty,
			custodian, offloading_charges, transport_detention, dispatch_info
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			exam_date = VALUES(exam_date),
			out_of_charge = VALUES(out_of_charge),
			clearance_exps = VALUES(clearance_exps),
			stamp_duty = VALUES(stamp_duty),
			custodian = VALUES(custodian),
			offloading_charges = VALUES(offloading_charges),
			transport_detention = VALUES(transport_detention),
			dispatch_info = VALUES(dispatch_info),
			updated_at = CURRENT_TIMESTAMP
	`,
		jobID, parseDate(req.ExamDate), parseDate(req.OutOfCharge),
		req.ClearanceExps, req.StampDuty, req.Custodian,
		req.OffloadingCharges, req.TransportDetention, req.DispatchInfo,
	)
	if err != nil {
		log.Printf("Error executing stage3 data insert/update: %v", err)
		return err
	}
	stage3RowsAffected, _ := stage3Result.RowsAffected()
	log.Printf("Stage3 data insert/update affected %d rows", stage3RowsAffected)

	// Delete existing containers and add new ones
	_, err = tx.Exec("DELETE FROM stage3_containers WHERE job_id = ?", jobID)
	if err != nil {
		return err
	}

	for _, container := range req.Containers {
		_, err = tx.Exec(`
			INSERT INTO stage3_containers (job_id, container_no, size, vehicle_no, date_of_offloading, empty_return_date)
			VALUES (?, ?, ?, ?, ?, ?)
		`,
			jobID, container.ContainerNo, container.Size, container.VehicleNo,
			parseDate(container.DateOfOffloading), parseDate(container.EmptyReturnDate),
		)
		if err != nil {
			return err
		}
	}

	// Update job stage if not already in stage3 or beyond
	log.Printf("Updating job stage to stage3 for job %d", jobID)
	jobResult, err := tx.Exec(`
		UPDATE pipeline_jobs 
		SET current_stage = 'stage3', updated_at = CURRENT_TIMESTAMP 
		WHERE id = ? AND current_stage IN ('stage1', 'stage2')
	`, jobID)
	if err != nil {
		log.Printf("Error updating job stage: %v", err)
		return err
	}
	jobRowsAffected, _ := jobResult.RowsAffected()
	log.Printf("Job stage update affected %d rows", jobRowsAffected)

	// Add job update
	_, err = tx.Exec(`
		INSERT INTO job_updates (job_id, user_id, stage, update_type, message)
		VALUES (?, ?, 'stage3', 'data_update', 'Stage 3 data updated')
	`, jobID, userID)
	if err != nil {
		log.Printf("Error adding job update: %v", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}
	log.Printf("Stage3 data updated successfully for job %d", jobID)
	return nil
}

// UpdateStage4Data updates stage 4 data and completes the job
func (r *PipelineRepository) UpdateStage4Data(jobID int, req *models.Stage4UpdateRequest, userID int) error {
	log.Printf("UpdateStage4Data called with jobID: %d, userID: %d", jobID, userID)
	log.Printf("Stage4 data: %+v", req)
	
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert or update stage4 data
	log.Printf("Executing stage4 data insert/update for job %d", jobID)
	log.Printf("BillNo: %s", req.BillNo)
	log.Printf("BillMail: %s", req.BillMail)
	log.Printf("AmountTaxable: %f", req.AmountTaxable)
	log.Printf("GST5Percent: %f", req.GST5Percent)
	log.Printf("GST18Percent: %f", req.GST18Percent)
	
	stage4Result, err := tx.Exec(`
		INSERT INTO stage4_data (
			job_id, bill_no, bill_date, amount_taxable, gst_5_percent, gst_18_percent,
			bill_mail, bill_courier, courier_date, acknowledge_date, acknowledge_name
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			bill_no = VALUES(bill_no),
			bill_date = VALUES(bill_date),
			amount_taxable = VALUES(amount_taxable),
			gst_5_percent = VALUES(gst_5_percent),
			gst_18_percent = VALUES(gst_18_percent),
			bill_mail = VALUES(bill_mail),
			bill_courier = VALUES(bill_courier),
			courier_date = VALUES(courier_date),
			acknowledge_date = VALUES(acknowledge_date),
			acknowledge_name = VALUES(acknowledge_name),
			updated_at = CURRENT_TIMESTAMP
	`,
		jobID, req.BillNo, parseDate(req.BillDate), req.AmountTaxable,
		req.GST5Percent, req.GST18Percent, req.BillMail, req.BillCourier,
		parseDate(req.CourierDate), parseDate(req.AcknowledgeDate), req.AcknowledgeName,
	)
	if err != nil {
		log.Printf("Error executing stage4 data insert/update: %v", err)
		return err
	}
	stage4RowsAffected, _ := stage4Result.RowsAffected()
	log.Printf("Stage4 data insert/update affected %d rows", stage4RowsAffected)

	// Update job stage to stage4 and potentially completed
	var newStage string
	if req.AcknowledgeDate != "" {
		newStage = "completed"
	} else {
		newStage = "stage4"
	}

	log.Printf("Updating job stage to %s for job %d", newStage, jobID)
	jobResult, err := tx.Exec(`
		UPDATE pipeline_jobs 
		SET current_stage = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?
	`, newStage, jobID)
	if err != nil {
		log.Printf("Error updating job stage: %v", err)
		return err
	}
	jobRowsAffected, _ := jobResult.RowsAffected()
	log.Printf("Job stage update affected %d rows", jobRowsAffected)

	// Add job update
	message := "Stage 4 data updated"
	if newStage == "completed" {
		message = "Job completed"
	}

	_, err = tx.Exec(`
		INSERT INTO job_updates (job_id, user_id, stage, update_type, message)
		VALUES (?, ?, 'stage4', 'data_update', ?)
	`, jobID, userID, message)
	if err != nil {
		log.Printf("Error adding job update: %v", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}
	log.Printf("Stage4 data updated successfully for job %d", jobID)
	return nil
}

// Helper functions
func (r *PipelineRepository) loadJobStageData(job *models.PipelineJobResponse) error {
	log.Printf("loadJobStageData called for job %d (current stage: %s)", job.ID, job.CurrentStage)
	
	// Load Stage 1 data
	stage1, err := r.getStage1Data(job.ID)
	if err == nil {
		job.Stage1 = stage1
		log.Printf("Stage1 data loaded successfully for job %d", job.ID)
	} else {
		log.Printf("Failed to load Stage1 data for job %d: %v", job.ID, err)
	}

	// Load Stage 2 data if applicable
	if job.CurrentStage == "stage2" || job.CurrentStage == "stage3" || job.CurrentStage == "stage4" || job.CurrentStage == "completed" {
		stage2, err := r.getStage2Data(job.ID)
		if err == nil {
			job.Stage2 = stage2
		}
	}

	// Load Stage 3 data if applicable
	if job.CurrentStage == "stage3" || job.CurrentStage == "stage4" || job.CurrentStage == "completed" {
		stage3, err := r.getStage3Data(job.ID)
		if err == nil {
			job.Stage3 = stage3
			log.Printf("Stage3 data loaded successfully for job %d", job.ID)
		} else {
			log.Printf("Failed to load Stage3 data for job %d: %v", job.ID, err)
		}

		containers, err := r.getStage3Containers(job.ID)
		if err == nil {
			job.Stage3Containers = containers
			log.Printf("Stage3 containers loaded successfully for job %d", job.ID)
		} else {
			log.Printf("Failed to load Stage3 containers for job %d: %v", job.ID, err)
		}
	}

	// Load Stage 4 data if applicable
	if job.CurrentStage == "stage4" || job.CurrentStage == "completed" {
		stage4, err := r.getStage4Data(job.ID)
		if err == nil {
			job.Stage4 = stage4
			log.Printf("Stage4 data loaded successfully for job %d", job.ID)
		} else {
			log.Printf("Failed to load Stage4 data for job %d: %v", job.ID, err)
		}
	}

	return nil
}

func (r *PipelineRepository) getStage1Data(jobID int) (*models.Stage1Data, error) {
	log.Printf("getStage1Data called for jobID: %d", jobID)
	
	var stage1 models.Stage1Data
	query := `
		SELECT id, job_id, job_no, job_date, edi_job_no, edi_date, consignee, shipper,
			   port_of_discharge, final_place_of_delivery, port_of_loading, country_of_shipment,
			   hbl_no, hbl_date, mbl_no, mbl_date, shipping_line, forwarder,
			   weight, packages, invoice_no, invoice_date, gateway_igm, gateway_igm_date,
			   local_igm, local_igm_date, commodity, eta, current_status,
			   container_no, container_size, date_of_arrival, invoice_pl_doc, bl_doc, coo_doc,
			   created_at, updated_at
		FROM stage1_data WHERE job_id = ?
	`

	err := r.db.QueryRow(query, jobID).Scan(
		&stage1.ID, &stage1.JobID, &stage1.JobNo, &stage1.JobDate, &stage1.EDIJobNo,
		&stage1.EDIDate, &stage1.Consignee, &stage1.Shipper, &stage1.PortOfDischarge,
		&stage1.FinalPlaceOfDelivery, &stage1.PortOfLoading, &stage1.CountryOfShipment,
		&stage1.HBLNo, &stage1.HBLDate, &stage1.MBLNo, &stage1.MBLDate,
		&stage1.ShippingLine, &stage1.Forwarder, &stage1.Weight, &stage1.Packages,
		&stage1.InvoiceNo, &stage1.InvoiceDate, &stage1.GatewayIGM, &stage1.GatewayIGMDate,
		&stage1.LocalIGM, &stage1.LocalIGMDate, &stage1.Commodity, &stage1.ETA,
		&stage1.CurrentStatus, &stage1.ContainerNo, &stage1.ContainerSize,
		&stage1.DateOfArrival, &stage1.InvoicePLDoc, &stage1.BLDoc, &stage1.COODoc,
		&stage1.CreatedAt, &stage1.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error getting stage1 data for job %d: %v", jobID, err)
		return nil, err
	}

	log.Printf("Stage1 data found for job %d: Consignee=%s, Commodity=%s", jobID, stage1.Consignee, stage1.Commodity)
	return &stage1, nil
}

func (r *PipelineRepository) getStage2Data(jobID int) (*models.Stage2Data, error) {
	log.Printf("getStage2Data called for jobID: %d", jobID)
	
	var stage2 models.Stage2Data
	query := `
		SELECT id, job_id, hsn_code, filing_requirement, checklist_sent_date, approval_date,
			   bill_of_entry_no, bill_of_entry_date, debit_note, debit_paid_by,
			   duty_amount, duty_paid_by, ocean_freight, destination_charges,
			   original_doct_recd_date, drn_no, irn_no, documents_type,
			   document_1, document_2, document_3, document_4, document_5, document_6,
			   query_upload, reply_upload, created_at, updated_at
		FROM stage2_data WHERE job_id = ?
	`

	err := r.db.QueryRow(query, jobID).Scan(
		&stage2.ID, &stage2.JobID, &stage2.HSNCode, &stage2.FilingRequirement,
		&stage2.ChecklistSentDate, &stage2.ApprovalDate, &stage2.BillOfEntryNo,
		&stage2.BillOfEntryDate, &stage2.DebitNote, &stage2.DebitPaidBy,
		&stage2.DutyAmount, &stage2.DutyPaidBy, &stage2.OceanFreight,
		&stage2.DestinationCharges, &stage2.OriginalDoctRecdDate, &stage2.DRNNo,
		&stage2.IRNNo, &stage2.DocumentsType, &stage2.Document1, &stage2.Document2,
		&stage2.Document3, &stage2.Document4, &stage2.Document5, &stage2.Document6,
		&stage2.QueryUpload, &stage2.ReplyUpload, &stage2.CreatedAt, &stage2.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error getting stage2 data for job %d: %v", jobID, err)
		return nil, err
	}

	log.Printf("Stage2 data found for job %d: HSNCode=%s, FilingRequirement=%s", jobID, stage2.HSNCode, stage2.FilingRequirement)
	return &stage2, nil
}

func (r *PipelineRepository) getStage3Data(jobID int) (*models.Stage3Data, error) {
	log.Printf("getStage3Data called for jobID: %d", jobID)
	
	var stage3 models.Stage3Data
	query := `
		SELECT id, job_id, exam_date, out_of_charge, clearance_exps, stamp_duty,
			   custodian, offloading_charges, transport_detention, dispatch_info,
			   bill_of_entry_upload, created_at, updated_at
		FROM stage3_data WHERE job_id = ?
	`

	err := r.db.QueryRow(query, jobID).Scan(
		&stage3.ID, &stage3.JobID, &stage3.ExamDate, &stage3.OutOfCharge,
		&stage3.ClearanceExps, &stage3.StampDuty, &stage3.Custodian,
		&stage3.OffloadingCharges, &stage3.TransportDetention, &stage3.DispatchInfo,
		&stage3.BillOfEntryUpload, &stage3.CreatedAt, &stage3.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error getting stage3 data for job %d: %v", jobID, err)
		return nil, err
	}

	log.Printf("Stage3 data found for job %d: Custodian=%s, DispatchInfo=%s", jobID, stage3.Custodian, stage3.DispatchInfo)
	return &stage3, err
}

func (r *PipelineRepository) getStage3Containers(jobID int) ([]models.Stage3Container, error) {
	query := `
		SELECT id, job_id, container_no, size, vehicle_no, date_of_offloading, empty_return_date, created_at
		FROM stage3_containers WHERE job_id = ?
	`

	rows, err := r.db.Query(query, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var containers []models.Stage3Container
	for rows.Next() {
		var container models.Stage3Container
		err := rows.Scan(
			&container.ID, &container.JobID, &container.ContainerNo, &container.Size,
			&container.VehicleNo, &container.DateOfOffloading, &container.EmptyReturnDate,
			&container.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		containers = append(containers, container)
	}

	return containers, nil
}

func (r *PipelineRepository) getStage4Data(jobID int) (*models.Stage4Data, error) {
	log.Printf("getStage4Data called for jobID: %d", jobID)
	
	var stage4 models.Stage4Data
	query := `
		SELECT id, job_id, bill_no, bill_date, amount_taxable, gst_5_percent, gst_18_percent,
			   bill_mail, bill_courier, courier_date, acknowledge_date, acknowledge_name,
			   bill_copy_upload, created_at, updated_at
		FROM stage4_data WHERE job_id = ?
	`

	err := r.db.QueryRow(query, jobID).Scan(
		&stage4.ID, &stage4.JobID, &stage4.BillNo, &stage4.BillDate,
		&stage4.AmountTaxable, &stage4.GST5Percent, &stage4.GST18Percent,
		&stage4.BillMail, &stage4.BillCourier, &stage4.CourierDate,
		&stage4.AcknowledgeDate, &stage4.AcknowledgeName, &stage4.BillCopyUpload,
		&stage4.CreatedAt, &stage4.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error getting stage4 data for job %d: %v", jobID, err)
		return nil, err
	}

	log.Printf("Stage4 data found for job %d: BillNo=%s, BillMail=%s", jobID, stage4.BillNo, stage4.BillMail)
	return &stage4, err
}

func (r *PipelineRepository) getJobUpdates(jobID int) ([]models.JobUpdate, error) {
	query := `
		SELECT id, job_id, user_id, stage, update_type, message, old_value, new_value, created_at
		FROM job_updates WHERE job_id = ? ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var updates []models.JobUpdate
	for rows.Next() {
		var update models.JobUpdate
		err := rows.Scan(
			&update.ID, &update.JobID, &update.UserID, &update.Stage, &update.UpdateType,
			&update.Message, &update.OldValue, &update.NewValue, &update.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		updates = append(updates, update)
	}

	return updates, nil
}

// Utility functions
func parseDate(dateStr string) interface{} {
	if dateStr == "" {
		return nil
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil
	}
	return date
}

func parseDateTime(dateStr string) interface{} {
	if dateStr == "" {
		return nil
	}
	date, err := time.Parse("2006-01-02T15:04:05", dateStr)
	if err != nil {
		// Try date only format
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil
		}
	}
	return date
}

func nullInt(val int) interface{} {
	if val == 0 {
		return nil
	}
	return val
} 