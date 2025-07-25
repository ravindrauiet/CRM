package database

import (
	"log"
	"strings"
)

func (db *DB) Migrate() error {
	log.Println("Running database migrations...")

	// SQL schema for 4-stage pipeline workflow
	schema := `
	-- Drop existing tables if they exist
	DROP TABLE IF EXISTS job_files;
	DROP TABLE IF EXISTS task_updates;
	DROP TABLE IF EXISTS task_assignments;
	DROP TABLE IF EXISTS tasks;
	DROP TABLE IF EXISTS job_updates;
	DROP TABLE IF EXISTS stage4_data;
	DROP TABLE IF EXISTS stage3_containers;
	DROP TABLE IF EXISTS stage3_data;
	DROP TABLE IF EXISTS stage2_data;
	DROP TABLE IF EXISTS stage1_data;
	DROP TABLE IF EXISTS pipeline_jobs;

	-- Users table (updated)
	DROP TABLE IF EXISTS users;
	CREATE TABLE users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(50) NOT NULL UNIQUE,
		password_hash VARCHAR(255) NOT NULL,
		designation VARCHAR(50) NOT NULL,
		is_admin BOOLEAN DEFAULT FALSE,
		role ENUM('admin', 'subadmin', 'stage1_employee', 'stage2_employee', 'stage3_employee', 'customer') DEFAULT 'stage1_employee',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Pipeline Jobs table (Main job tracking)
	CREATE TABLE pipeline_jobs (
		id INT AUTO_INCREMENT PRIMARY KEY,
		job_no VARCHAR(50) NOT NULL UNIQUE,
		current_stage ENUM('stage1', 'stage2', 'stage3', 'stage4', 'completed') DEFAULT 'stage1',
		status ENUM('active', 'on_hold', 'completed', 'cancelled') DEFAULT 'active',
		created_by INT NOT NULL,
		assigned_to_stage2 INT NULL,
		assigned_to_stage3 INT NULL,
		customer_id INT NULL,
		notification_email VARCHAR(255) DEFAULT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (created_by) REFERENCES users(id),
		FOREIGN KEY (assigned_to_stage2) REFERENCES users(id),
		FOREIGN KEY (assigned_to_stage3) REFERENCES users(id),
		FOREIGN KEY (customer_id) REFERENCES users(id)
	);

	-- Stage 1: Initial Job Creation (Admin)
	CREATE TABLE stage1_data (
		id INT AUTO_INCREMENT PRIMARY KEY,
		job_id INT NOT NULL,
		job_no VARCHAR(50) NOT NULL,
		job_date DATE,
		edi_job_no VARCHAR(50),
		edi_date DATE,
		consignee TEXT,
		shipper TEXT,
		port_of_discharge VARCHAR(100),
		final_place_of_delivery VARCHAR(100),
		port_of_loading VARCHAR(100),
		country_of_shipment VARCHAR(100),
		hbl_no VARCHAR(50),
		hbl_date DATE,
		mbl_no VARCHAR(50),
		mbl_date DATE,
		shipping_line VARCHAR(100),
		forwarder VARCHAR(100),
		weight DECIMAL(10,2),
		packages INT,
		invoice_no VARCHAR(50),
		invoice_date DATE,
		gateway_igm VARCHAR(50),
		gateway_igm_date DATE,
		local_igm VARCHAR(50),
		local_igm_date DATE,
		commodity TEXT,
		eta DATETIME,
		current_status VARCHAR(100),
		container_no VARCHAR(50),
		container_size ENUM('20', '40', 'LCL'),
		date_of_arrival DATE,
		invoice_pl_doc VARCHAR(255),
		bl_doc VARCHAR(255),
		coo_doc VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (job_id) REFERENCES pipeline_jobs(id) ON DELETE CASCADE
	);

	-- Stage 2: Customs & Documentation (Employee)
	CREATE TABLE stage2_data (
		id INT AUTO_INCREMENT PRIMARY KEY,
		job_id INT NOT NULL,
		hsn_code VARCHAR(20),
		filing_requirement TEXT,
		checklist_sent_date DATE,
		approval_date DATE,
		bill_of_entry_no VARCHAR(50),
		bill_of_entry_date DATE,
		debit_note VARCHAR(50),
		debit_paid_by VARCHAR(100),
		duty_amount DECIMAL(10,2),
		duty_paid_by VARCHAR(100),
		ocean_freight DECIMAL(10,2),
		destination_charges DECIMAL(10,2),
		original_doct_recd_date DATE,
		drn_no VARCHAR(50),
		irn_no VARCHAR(50),
		documents_type VARCHAR(100),
		document_1 VARCHAR(255),
		document_2 VARCHAR(255),
		document_3 VARCHAR(255),
		document_4 VARCHAR(255),
		document_5 VARCHAR(255),
		document_6 VARCHAR(255),
		query_upload VARCHAR(255),
		reply_upload VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (job_id) REFERENCES pipeline_jobs(id) ON DELETE CASCADE
	);

	-- Stage 3: Clearance & Logistics (Employee)
	CREATE TABLE stage3_data (
		id INT AUTO_INCREMENT PRIMARY KEY,
		job_id INT NOT NULL,
		exam_date DATE,
		out_of_charge DATE,
		clearance_exps DECIMAL(10,2),
		stamp_duty DECIMAL(10,2),
		custodian VARCHAR(100),
		offloading_charges DECIMAL(10,2),
		transport_detention DECIMAL(10,2),
		dispatch_info TEXT,
		bill_of_entry_upload VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (job_id) REFERENCES pipeline_jobs(id) ON DELETE CASCADE
	);

	-- Stage 3: Container Details (Multiple containers per job)
	CREATE TABLE stage3_containers (
		id INT AUTO_INCREMENT PRIMARY KEY,
		job_id INT NOT NULL,
		container_no VARCHAR(50),
		size ENUM('20', '40', 'LCL'),
		vehicle_no VARCHAR(50),
		date_of_offloading DATE,
		empty_return_date DATE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (job_id) REFERENCES pipeline_jobs(id) ON DELETE CASCADE
	);

	-- Stage 4: Billing & Customer (Customer/Admin)
	CREATE TABLE stage4_data (
		id INT AUTO_INCREMENT PRIMARY KEY,
		job_id INT NOT NULL,
		bill_no VARCHAR(50),
		bill_date DATE,
		amount_taxable DECIMAL(10,2),
		gst_5_percent DECIMAL(10,2),
		gst_18_percent DECIMAL(10,2),
		bill_mail VARCHAR(255),
		bill_courier VARCHAR(100),
		courier_date DATE,
		acknowledge_date DATE,
		acknowledge_name VARCHAR(100),
		bill_copy_upload VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (job_id) REFERENCES pipeline_jobs(id) ON DELETE CASCADE
	);

	-- Job Files (File uploads for each stage)
	CREATE TABLE job_files (
		id INT AUTO_INCREMENT PRIMARY KEY,
		job_id INT NOT NULL,
		stage ENUM('stage1', 'stage2', 'stage3', 'stage4') NOT NULL,
		uploaded_by INT NOT NULL,
		file_name VARCHAR(255) NOT NULL,
		original_name VARCHAR(255) NOT NULL,
		file_path VARCHAR(500) NOT NULL,
		file_size BIGINT NOT NULL,
		file_type VARCHAR(100),
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (job_id) REFERENCES pipeline_jobs(id) ON DELETE CASCADE,
		FOREIGN KEY (uploaded_by) REFERENCES users(id)
	);

	-- Job Updates/Comments (Timeline)
	CREATE TABLE job_updates (
		id INT AUTO_INCREMENT PRIMARY KEY,
		job_id INT NOT NULL,
		user_id INT NOT NULL,
		stage ENUM('stage1', 'stage2', 'stage3', 'stage4') NOT NULL,
		update_type ENUM('status_change', 'data_update', 'comment', 'stage_completion', 'file_upload') NOT NULL,
		message TEXT,
		old_value VARCHAR(255),
		new_value VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (job_id) REFERENCES pipeline_jobs(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	// Split schema into individual statements and execute them
	statements := strings.Split(schema, ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		
		log.Printf("Executing: %s", strings.Split(stmt, "\n")[0])
		_, err := db.Exec(stmt)
		if err != nil {
			log.Printf("Error executing statement: %v", err)
			return err
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}

func (db *DB) Seed() error {
	log.Println("Seeding database with sample data...")

	// Clear existing users first
	_, err := db.Exec("DELETE FROM users")
	if err != nil {
		log.Printf("Error clearing users: %v", err)
	}

	// Insert sample users for different roles
	seedSQL := `
	INSERT INTO users (username, password_hash, designation, is_admin, role) VALUES
	('admin', '123456', 'Administrator', TRUE, 'admin'),
	('stage2_emp', '123456', 'Customs Officer', FALSE, 'stage2_employee'),
	('stage3_emp', '123456', 'Logistics Coordinator', FALSE, 'stage3_employee'),
	('customer1', '123456', 'Client', FALSE, 'customer'),
	('subadmin', '123456', 'Sub Administrator', FALSE, 'subadmin');

	INSERT INTO pipeline_jobs (job_no, current_stage, created_by, assigned_to_stage2, assigned_to_stage3, customer_id) 
	VALUES ('JOB001', 'stage1', 1, 2, 3, 4);

	INSERT INTO stage1_data (job_id, job_no, job_date, consignee, shipper, commodity, current_status) 
	VALUES (1, 'JOB001', CURDATE(), 'ABC Import Co.', 'XYZ Export Ltd.', 'Electronics', 'Documents Received');`

	statements := strings.Split(seedSQL, ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		
		_, err := db.Exec(stmt)
		if err != nil {
			log.Printf("Error seeding data: %v", err)
			return err
		}
	}

	log.Println("Database seeded successfully")
	return nil
} 