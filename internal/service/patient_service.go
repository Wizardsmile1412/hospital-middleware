package service

import (
	"context"
	"log"

	"github.com/Wizardsmile1412/hospital-middleware/internal/client"
	"github.com/Wizardsmile1412/hospital-middleware/internal/model"
	"github.com/Wizardsmile1412/hospital-middleware/internal/repository"
)

type PatientService interface {
	SearchPatients(ctx context.Context, params model.SearchPatientParams) (*model.SearchPatientResponse, error)
}

type patientService struct {
	repo           repository.PatientRepository
	hospitalClient client.HospitalClient
}

func NewPatientService(repo repository.PatientRepository, hospitalClient client.HospitalClient) PatientService {
	return &patientService{repo: repo, hospitalClient: hospitalClient}
}

func (s *patientService) SearchPatients(ctx context.Context, params model.SearchPatientParams) (*model.SearchPatientResponse, error) {
	// national_id or passport_id → real-time call to Hospital external API
	if params.NationalID != "" {
		patient, err := s.hospitalClient.SearchByNationalID(ctx, params.NationalID)
		if err != nil {
			log.Printf("hospital API error for national_id, falling back to DB: %v", err)
		} else if resp := buildResponse(patient, params.HospitalID); resp.Total > 0 {
			return resp, nil
		}
	}

	if params.PassportID != "" {
		patient, err := s.hospitalClient.SearchByPassportID(ctx, params.PassportID)
		if err != nil {
			log.Printf("hospital API error for passport_id, falling back to DB: %v", err)
		} else if resp := buildResponse(patient, params.HospitalID); resp.Total > 0 {
			return resp, nil
		}
	}

	// Fallback: query our own DB (also handles cases where external API errored or returned no result)
	patients, err := s.repo.SearchPatients(ctx, params)
	if err != nil {
		return nil, err
	}

	return &model.SearchPatientResponse{
		Patients: patients,
		Total:    len(patients),
	}, nil
}

// buildResponse wraps a single external-API patient result into the standard response.
// It filters out patients that don't belong to the requesting hospital for safety.
func buildResponse(patient *model.Patient, hospitalID string) *model.SearchPatientResponse {
	if patient == nil || patient.HospitalID != hospitalID {
		return &model.SearchPatientResponse{Patients: []model.Patient{}, Total: 0}
	}
	return &model.SearchPatientResponse{
		Patients: []model.Patient{*patient},
		Total:    1,
	}
}
