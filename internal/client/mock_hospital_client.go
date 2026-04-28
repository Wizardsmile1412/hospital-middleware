package client

import (
	"context"

	"github.com/Wizardsmile1412/hospital-middleware/internal/model"
)

// MockHospitalClient is used in tests. It is fully interchangeable with the real
// client because both satisfy HospitalClient (Liskov Substitution).
type MockHospitalClient struct {
	SearchByNationalIDFunc func(ctx context.Context, nationalID string) (*model.Patient, error)
	SearchByPassportIDFunc func(ctx context.Context, passportID string) (*model.Patient, error)
}

func (m *MockHospitalClient) SearchByNationalID(ctx context.Context, nationalID string) (*model.Patient, error) {
	if m.SearchByNationalIDFunc != nil {
		return m.SearchByNationalIDFunc(ctx, nationalID)
	}
	return nil, nil
}

func (m *MockHospitalClient) SearchByPassportID(ctx context.Context, passportID string) (*model.Patient, error) {
	if m.SearchByPassportIDFunc != nil {
		return m.SearchByPassportIDFunc(ctx, passportID)
	}
	return nil, nil
}
