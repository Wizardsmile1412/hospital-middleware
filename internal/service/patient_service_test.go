package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Wizardsmile1412/hospital-middleware/internal/client"
	"github.com/Wizardsmile1412/hospital-middleware/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockPatientRepo implements repository.PatientRepository for tests.
type mockPatientRepo struct {
	SearchPatientsFunc func(ctx context.Context, params model.SearchPatientParams) ([]model.Patient, error)
}

func (m *mockPatientRepo) SearchPatients(ctx context.Context, params model.SearchPatientParams) ([]model.Patient, error) {
	if m.SearchPatientsFunc != nil {
		return m.SearchPatientsFunc(ctx, params)
	}
	return []model.Patient{}, nil
}

func TestSearchPatients_ByNationalID_CallsHospitalClient(t *testing.T) {
	nationalID := "1234567890123"
	hospitalID := "hospital-uuid"
	expectedPatient := &model.Patient{
		ID:         "patient-uuid",
		HospitalID: hospitalID,
		NationalID: &nationalID,
	}

	mockClient := &client.MockHospitalClient{
		SearchByNationalIDFunc: func(_ context.Context, id string) (*model.Patient, error) {
			assert.Equal(t, nationalID, id)
			return expectedPatient, nil
		},
	}
	svc := NewPatientService(&mockPatientRepo{}, mockClient)

	result, err := svc.SearchPatients(context.Background(), model.SearchPatientParams{
		HospitalID: hospitalID,
		NationalID: nationalID,
	})

	require.NoError(t, err)
	assert.Equal(t, 1, result.Total)
	assert.Equal(t, expectedPatient.ID, result.Patients[0].ID)
}

func TestSearchPatients_ByPassportID_CallsHospitalClient(t *testing.T) {
	passportID := "AB123456"
	hospitalID := "hospital-uuid"
	expectedPatient := &model.Patient{
		ID:         "patient-uuid",
		HospitalID: hospitalID,
		PassportID: &passportID,
	}

	mockClient := &client.MockHospitalClient{
		SearchByPassportIDFunc: func(_ context.Context, id string) (*model.Patient, error) {
			assert.Equal(t, passportID, id)
			return expectedPatient, nil
		},
	}
	svc := NewPatientService(&mockPatientRepo{}, mockClient)

	result, err := svc.SearchPatients(context.Background(), model.SearchPatientParams{
		HospitalID: hospitalID,
		PassportID: passportID,
	})

	require.NoError(t, err)
	assert.Equal(t, 1, result.Total)
}

func TestSearchPatients_ByDB_NoIdentifier(t *testing.T) {
	firstName := "John"
	patients := []model.Patient{
		{ID: "p1", HospitalID: "hospital-uuid", FirstNameEN: &firstName},
	}

	repo := &mockPatientRepo{
		SearchPatientsFunc: func(_ context.Context, _ model.SearchPatientParams) ([]model.Patient, error) {
			return patients, nil
		},
	}
	svc := NewPatientService(repo, &client.MockHospitalClient{})

	result, err := svc.SearchPatients(context.Background(), model.SearchPatientParams{
		HospitalID: "hospital-uuid",
		FirstName:  "John",
	})

	require.NoError(t, err)
	assert.Equal(t, 1, result.Total)
	assert.Equal(t, "p1", result.Patients[0].ID)
}

func TestSearchPatients_HospitalIsolation(t *testing.T) {
	// Patient belongs to a different hospital — buildResponse filters it out.
	otherHospitalID := "other-hospital-uuid"
	nationalID := "1234567890123"
	patientFromOtherHospital := &model.Patient{
		ID:         "patient-uuid",
		HospitalID: otherHospitalID,
		NationalID: &nationalID,
	}

	mockClient := &client.MockHospitalClient{
		SearchByNationalIDFunc: func(_ context.Context, _ string) (*model.Patient, error) {
			return patientFromOtherHospital, nil
		},
	}
	svc := NewPatientService(&mockPatientRepo{}, mockClient)

	result, err := svc.SearchPatients(context.Background(), model.SearchPatientParams{
		HospitalID: "requesting-hospital-uuid",
		NationalID: nationalID,
	})

	require.NoError(t, err)
	assert.Equal(t, 0, result.Total)
	assert.Empty(t, result.Patients)
}

func TestSearchPatients_HospitalClientError(t *testing.T) {
	mockClient := &client.MockHospitalClient{
		SearchByNationalIDFunc: func(_ context.Context, _ string) (*model.Patient, error) {
			return nil, errors.New("hospital API unreachable")
		},
	}
	svc := NewPatientService(&mockPatientRepo{}, mockClient)

	_, err := svc.SearchPatients(context.Background(), model.SearchPatientParams{
		HospitalID: "hospital-uuid",
		NationalID: "1234567890123",
	})

	assert.Error(t, err)
}
