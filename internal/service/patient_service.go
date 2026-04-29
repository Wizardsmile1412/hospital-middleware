package service

import (
	"context"

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
			return nil, err
		}
		return buildResponse(patient, params.HospitalID), nil
	}

	if params.PassportID != "" {
		patient, err := s.hospitalClient.SearchByPassportID(ctx, params.PassportID)
		if err != nil {
			return nil, err
		}
		return buildResponse(patient, params.HospitalID), nil
	}

	// All other filters → query our own DB (fast, reliable, always hospital-scoped)
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
