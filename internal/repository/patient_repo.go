package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Wizardsmile1412/hospital-middleware/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PatientRepository interface {
	SearchPatients(ctx context.Context, params model.SearchPatientParams) ([]model.Patient, error)
}

type patientRepo struct {
	db *pgxpool.Pool
}

func NewPatientRepository(db *pgxpool.Pool) PatientRepository {
	return &patientRepo{db: db}
}

func (r *patientRepo) SearchPatients(ctx context.Context, params model.SearchPatientParams) ([]model.Patient, error) {
	conditions := []string{"hospital_id = $1"}
	args := []any{params.HospitalID}
	idx := 2

	if params.NationalID != "" {
		conditions = append(conditions, fmt.Sprintf("national_id = $%d", idx))
		args = append(args, params.NationalID)
		idx++
	}
	if params.PassportID != "" {
		conditions = append(conditions, fmt.Sprintf("passport_id = $%d", idx))
		args = append(args, params.PassportID)
		idx++
	}
	if params.FirstName != "" {
		conditions = append(conditions, fmt.Sprintf(
			"(LOWER(first_name_th) LIKE LOWER($%d) OR LOWER(first_name_en) LIKE LOWER($%d))",
			idx, idx,
		))
		args = append(args, "%"+params.FirstName+"%")
		idx++
	}
	if params.LastName != "" {
		conditions = append(conditions, fmt.Sprintf(
			"(LOWER(last_name_th) LIKE LOWER($%d) OR LOWER(last_name_en) LIKE LOWER($%d))",
			idx, idx,
		))
		args = append(args, "%"+params.LastName+"%")
		idx++
	}
	if params.DateOfBirth != "" {
		conditions = append(conditions, fmt.Sprintf("date_of_birth = $%d", idx))
		args = append(args, params.DateOfBirth)
		idx++
	}
	if params.PhoneNumber != "" {
		conditions = append(conditions, fmt.Sprintf("phone_number = $%d", idx))
		args = append(args, params.PhoneNumber)
		idx++
	}
	if params.Email != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(email) = LOWER($%d)", idx))
		args = append(args, params.Email)
		idx++
	}
	_ = idx

	query := fmt.Sprintf(`
		SELECT id, hospital_id, national_id, passport_id,
		       first_name_th, middle_name_th, last_name_th,
		       first_name_en, middle_name_en, last_name_en,
		       date_of_birth, patient_hn, phone_number, email, gender, created_at
		FROM patients
		WHERE %s
		ORDER BY created_at DESC
	`, strings.Join(conditions, " AND "))

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []model.Patient
	for rows.Next() {
		var p model.Patient
		err := rows.Scan(
			&p.ID, &p.HospitalID, &p.NationalID, &p.PassportID,
			&p.FirstNameTH, &p.MiddleNameTH, &p.LastNameTH,
			&p.FirstNameEN, &p.MiddleNameEN, &p.LastNameEN,
			&p.DateOfBirth, &p.PatientHN, &p.PhoneNumber, &p.Email, &p.Gender,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		patients = append(patients, p)
	}

	if patients == nil {
		patients = []model.Patient{}
	}
	return patients, rows.Err()
}
