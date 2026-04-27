package model

import "time"

type Patient struct {
	ID           string     `json:"id" db:"id"`
	HospitalID   string     `json:"hospital_id" db:"hospital_id"`
	NationalID   *string    `json:"national_id" db:"national_id"`
	PassportID   *string    `json:"passport_id" db:"passport_id"`
	FirstNameTH  *string    `json:"first_name_th" db:"first_name_th"`
	MiddleNameTH *string    `json:"middle_name_th" db:"middle_name_th"`
	LastNameTH   *string    `json:"last_name_th" db:"last_name_th"`
	FirstNameEN  *string    `json:"first_name_en" db:"first_name_en"`
	MiddleNameEN *string    `json:"middle_name_en" db:"middle_name_en"`
	LastNameEN   *string    `json:"last_name_en" db:"last_name_en"`
	DateOfBirth  *time.Time `json:"date_of_birth" db:"date_of_birth"`
	PatientHN    *string    `json:"patient_hn" db:"patient_hn"`
	PhoneNumber  *string    `json:"phone_number" db:"phone_number"`
	Email        *string    `json:"email" db:"email"`
	Gender       *string    `json:"gender" db:"gender"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
}

type SearchPatientParams struct {
	NationalID  string
	PassportID  string
	FirstName   string
	LastName    string
	DateOfBirth string
	PhoneNumber string
	Email       string
	HospitalID  string
}

type SearchPatientResponse struct {
	Patients []Patient `json:"patients"`
	Total    int       `json:"total"`
}
