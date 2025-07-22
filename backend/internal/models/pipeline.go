package models

import "time"

// PipelineJob represents the main job tracking
type PipelineJob struct {
	ID               int       `json:"id" db:"id"`
	JobNo            string    `json:"job_no" db:"job_no"`
	CurrentStage     string    `json:"current_stage" db:"current_stage"`
	Status           string    `json:"status" db:"status"`
	CreatedBy        int       `json:"created_by" db:"created_by"`
	AssignedToStage2 *int      `json:"assigned_to_stage2" db:"assigned_to_stage2"`
	AssignedToStage3 *int      `json:"assigned_to_stage3" db:"assigned_to_stage3"`
	CustomerID       *int      `json:"customer_id" db:"customer_id"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// Stage1Data represents initial job creation data (Admin)
type Stage1Data struct {
	ID                    int        `json:"id" db:"id"`
	JobID                 int        `json:"job_id" db:"job_id"`
	JobNo                 string     `json:"job_no" db:"job_no"`
	JobDate               *time.Time `json:"job_date" db:"job_date"`
	EDIJobNo              *string    `json:"edi_job_no" db:"edi_job_no"`
	EDIDate               *time.Time `json:"edi_date" db:"edi_date"`
	Consignee             *string    `json:"consignee" db:"consignee"`
	Shipper               *string    `json:"shipper" db:"shipper"`
	PortOfDischarge       *string    `json:"port_of_discharge" db:"port_of_discharge"`
	FinalPlaceOfDelivery  *string    `json:"final_place_of_delivery" db:"final_place_of_delivery"`
	PortOfLoading         *string    `json:"port_of_loading" db:"port_of_loading"`
	CountryOfShipment     *string    `json:"country_of_shipment" db:"country_of_shipment"`
	HBLNo                 *string    `json:"hbl_no" db:"hbl_no"`
	HBLDate               *time.Time `json:"hbl_date" db:"hbl_date"`
	MBLNo                 *string    `json:"mbl_no" db:"mbl_no"`
	MBLDate               *time.Time `json:"mbl_date" db:"mbl_date"`
	ShippingLine          *string    `json:"shipping_line" db:"shipping_line"`
	Forwarder             *string    `json:"forwarder" db:"forwarder"`
	Weight                *float64   `json:"weight" db:"weight"`
	Packages              *int       `json:"packages" db:"packages"`
	InvoiceNo             *string    `json:"invoice_no" db:"invoice_no"`
	InvoiceDate           *time.Time `json:"invoice_date" db:"invoice_date"`
	GatewayIGM            *string    `json:"gateway_igm" db:"gateway_igm"`
	GatewayIGMDate        *time.Time `json:"gateway_igm_date" db:"gateway_igm_date"`
	LocalIGM              *string    `json:"local_igm" db:"local_igm"`
	LocalIGMDate          *time.Time `json:"local_igm_date" db:"local_igm_date"`
	Commodity             *string    `json:"commodity" db:"commodity"`
	ETA                   *time.Time `json:"eta" db:"eta"`
	CurrentStatus         *string    `json:"current_status" db:"current_status"`
	ContainerNo           *string    `json:"container_no" db:"container_no"`
	ContainerSize         *string    `json:"container_size" db:"container_size"`
	DateOfArrival         *time.Time `json:"date_of_arrival" db:"date_of_arrival"`
	InvoicePLDoc          *string    `json:"invoice_pl_doc" db:"invoice_pl_doc"`
	BLDoc                 *string    `json:"bl_doc" db:"bl_doc"`
	COODoc                *string    `json:"coo_doc" db:"coo_doc"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at" db:"updated_at"`
}

// Stage2Data represents customs & documentation data (Employee)
type Stage2Data struct {
	ID                   int        `json:"id" db:"id"`
	JobID                int        `json:"job_id" db:"job_id"`
	HSNCode              *string    `json:"hsn_code" db:"hsn_code"`
	FilingRequirement    *string    `json:"filing_requirement" db:"filing_requirement"`
	ChecklistSentDate    *time.Time `json:"checklist_sent_date" db:"checklist_sent_date"`
	ApprovalDate         *time.Time `json:"approval_date" db:"approval_date"`
	BillOfEntryNo        *string    `json:"bill_of_entry_no" db:"bill_of_entry_no"`
	BillOfEntryDate      *time.Time `json:"bill_of_entry_date" db:"bill_of_entry_date"`
	DebitNote            *string    `json:"debit_note" db:"debit_note"`
	DebitPaidBy          *string    `json:"debit_paid_by" db:"debit_paid_by"`
	DutyAmount           float64    `json:"duty_amount" db:"duty_amount"`
	DutyPaidBy           *string    `json:"duty_paid_by" db:"duty_paid_by"`
	OceanFreight         float64    `json:"ocean_freight" db:"ocean_freight"`
	DestinationCharges   float64    `json:"destination_charges" db:"destination_charges"`
	OriginalDoctRecdDate *time.Time `json:"original_doct_recd_date" db:"original_doct_recd_date"`
	DRNNo                *string    `json:"drn_no" db:"drn_no"`
	IRNNo                *string    `json:"irn_no" db:"irn_no"`
	DocumentsType        *string    `json:"documents_type" db:"documents_type"`
	Document1            *string    `json:"document_1" db:"document_1"`
	Document2            *string    `json:"document_2" db:"document_2"`
	Document3            *string    `json:"document_3" db:"document_3"`
	Document4            *string    `json:"document_4" db:"document_4"`
	Document5            *string    `json:"document_5" db:"document_5"`
	Document6            *string    `json:"document_6" db:"document_6"`
	QueryUpload          *string    `json:"query_upload" db:"query_upload"`
	ReplyUpload          *string    `json:"reply_upload" db:"reply_upload"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
}

// Stage3Data represents clearance & logistics data (Employee)
type Stage3Data struct {
	ID                  int        `json:"id" db:"id"`
	JobID               int        `json:"job_id" db:"job_id"`
	ExamDate            *time.Time `json:"exam_date" db:"exam_date"`
	OutOfCharge         *time.Time `json:"out_of_charge" db:"out_of_charge"`
	ClearanceExps       *float64   `json:"clearance_exps" db:"clearance_exps"`
	StampDuty           *float64   `json:"stamp_duty" db:"stamp_duty"`
	Custodian           *string    `json:"custodian" db:"custodian"`
	OffloadingCharges   *float64   `json:"offloading_charges" db:"offloading_charges"`
	TransportDetention  *float64   `json:"transport_detention" db:"transport_detention"`
	DispatchInfo        *string    `json:"dispatch_info" db:"dispatch_info"`
	BillOfEntryUpload   *string    `json:"bill_of_entry_upload" db:"bill_of_entry_upload"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`
}

// Stage3Container represents container details for stage 3
type Stage3Container struct {
	ID                int        `json:"id" db:"id"`
	JobID             int        `json:"job_id" db:"job_id"`
	ContainerNo       *string    `json:"container_no" db:"container_no"`
	Size              *string    `json:"size" db:"size"`
	VehicleNo         *string    `json:"vehicle_no" db:"vehicle_no"`
	DateOfOffloading  *time.Time `json:"date_of_offloading" db:"date_of_offloading"`
	EmptyReturnDate   *time.Time `json:"empty_return_date" db:"empty_return_date"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
}

// Stage4Data represents billing & customer data
type Stage4Data struct {
	ID               int        `json:"id" db:"id"`
	JobID            int        `json:"job_id" db:"job_id"`
	BillNo           *string    `json:"bill_no" db:"bill_no"`
	BillDate         *time.Time `json:"bill_date" db:"bill_date"`
	AmountTaxable    *float64   `json:"amount_taxable" db:"amount_taxable"`
	GST5Percent      *float64   `json:"gst_5_percent" db:"gst_5_percent"`
	GST18Percent     *float64   `json:"gst_18_percent" db:"gst_18_percent"`
	BillMail         *string    `json:"bill_mail" db:"bill_mail"`
	BillCourier      *string    `json:"bill_courier" db:"bill_courier"`
	CourierDate      *time.Time `json:"courier_date" db:"courier_date"`
	AcknowledgeDate  *time.Time `json:"acknowledge_date" db:"acknowledge_date"`
	AcknowledgeName  *string    `json:"acknowledge_name" db:"acknowledge_name"`
	BillCopyUpload   *string    `json:"bill_copy_upload" db:"bill_copy_upload"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

// JobUpdate represents job timeline/comments
type JobUpdate struct {
	ID         int       `json:"id" db:"id"`
	JobID      int       `json:"job_id" db:"job_id"`
	UserID     int       `json:"user_id" db:"user_id"`
	Stage      string    `json:"stage" db:"stage"`
	UpdateType string    `json:"update_type" db:"update_type"`
	Message    string    `json:"message" db:"message"`
	OldValue   string    `json:"old_value" db:"old_value"`
	NewValue   string    `json:"new_value" db:"new_value"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// PipelineJobResponse represents complete job data for API responses
type PipelineJobResponse struct {
	PipelineJob `json:",inline"`
	Stage1      *Stage1Data        `json:"stage1,omitempty"`
	Stage2      *Stage2Data        `json:"stage2,omitempty"`
	Stage3      *Stage3Data        `json:"stage3,omitempty"`
	Stage3Containers []Stage3Container `json:"stage3_containers,omitempty"`
	Stage4      *Stage4Data        `json:"stage4,omitempty"`
	Updates     []JobUpdate        `json:"updates,omitempty"`
	CreatedByUser    string        `json:"created_by_user,omitempty"`
	Stage2UserName   string        `json:"stage2_user_name,omitempty"`
	Stage3UserName   string        `json:"stage3_user_name,omitempty"`
	CustomerName     string        `json:"customer_name,omitempty"`
}

// Create request structs
type Stage1CreateRequest struct {
	JobNo                 string `json:"job_no" binding:"required"`
	JobDate               string `json:"job_date"`
	EDIJobNo              string `json:"edi_job_no"`
	EDIDate               string `json:"edi_date"`
	Consignee             string `json:"consignee"`
	Shipper               string `json:"shipper"`
	PortOfDischarge       string `json:"port_of_discharge"`
	FinalPlaceOfDelivery  string `json:"final_place_of_delivery"`
	PortOfLoading         string `json:"port_of_loading"`
	CountryOfShipment     string `json:"country_of_shipment"`
	HBLNo                 string `json:"hbl_no"`
	HBLDate               string `json:"hbl_date"`
	MBLNo                 string `json:"mbl_no"`
	MBLDate               string `json:"mbl_date"`
	ShippingLine          string `json:"shipping_line"`
	Forwarder             string `json:"forwarder"`
	Weight                float64 `json:"weight"`
	Packages              int    `json:"packages"`
	InvoiceNo             string `json:"invoice_no"`
	InvoiceDate           string `json:"invoice_date"`
	GatewayIGM            string `json:"gateway_igm"`
	GatewayIGMDate        string `json:"gateway_igm_date"`
	LocalIGM              string `json:"local_igm"`
	LocalIGMDate          string `json:"local_igm_date"`
	Commodity             string `json:"commodity"`
	ETA                   string `json:"eta"`
	CurrentStatus         string `json:"current_status"`
	ContainerNo           string `json:"container_no"`
	ContainerSize         string `json:"container_size"`
	DateOfArrival         string `json:"date_of_arrival"`
	AssignedToStage2      int    `json:"assigned_to_stage2"`
	AssignedToStage3      int    `json:"assigned_to_stage3"`
	CustomerID            int    `json:"customer_id"`
}

type Stage2UpdateRequest struct {
	HSNCode              string  `json:"hsn_code"`
	FilingRequirement    string  `json:"filing_requirement"`
	ChecklistSentDate    string  `json:"checklist_sent_date"`
	ApprovalDate         string  `json:"approval_date"`
	BillOfEntryNo        string  `json:"bill_of_entry_no"`
	BillOfEntryDate      string  `json:"bill_of_entry_date"`
	DebitNote            string  `json:"debit_note"`
	DebitPaidBy          string  `json:"debit_paid_by"`
	DutyAmount           float64 `json:"duty_amount"`
	DutyPaidBy           string  `json:"duty_paid_by"`
	OceanFreight         float64 `json:"ocean_freight"`
	DestinationCharges   float64 `json:"destination_charges"`
	OriginalDoctRecdDate string  `json:"original_doct_recd_date"`
	DRNNo                string  `json:"drn_no"`
	IRNNo                string  `json:"irn_no"`
	DocumentsType        string  `json:"documents_type"`
}

type Stage3UpdateRequest struct {
	ExamDate           string  `json:"exam_date"`
	OutOfCharge        string  `json:"out_of_charge"`
	ClearanceExps      float64 `json:"clearance_exps"`
	StampDuty          float64 `json:"stamp_duty"`
	Custodian          string  `json:"custodian"`
	OffloadingCharges  float64 `json:"offloading_charges"`
	TransportDetention float64 `json:"transport_detention"`
	DispatchInfo       string  `json:"dispatch_info"`
	Containers         []Stage3ContainerRequest `json:"containers"`
}

type Stage3ContainerRequest struct {
	ContainerNo      string `json:"container_no"`
	Size             string `json:"size"`
	VehicleNo        string `json:"vehicle_no"`
	DateOfOffloading string `json:"date_of_offloading"`
	EmptyReturnDate  string `json:"empty_return_date"`
}

type Stage4UpdateRequest struct {
	BillNo          string  `json:"bill_no"`
	BillDate        string  `json:"bill_date"`
	AmountTaxable   float64 `json:"amount_taxable"`
	GST5Percent     float64 `json:"gst_5_percent"`
	GST18Percent    float64 `json:"gst_18_percent"`
	BillMail        string  `json:"bill_mail"`
	BillCourier     string  `json:"bill_courier"`
	CourierDate     string  `json:"courier_date"`
	AcknowledgeDate string  `json:"acknowledge_date"`
	AcknowledgeName string  `json:"acknowledge_name"`
} 