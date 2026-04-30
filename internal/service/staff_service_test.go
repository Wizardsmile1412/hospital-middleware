package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Wizardsmile1412/hospital-middleware/internal/model"
	"github.com/Wizardsmile1412/hospital-middleware/internal/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// mockStaffRepo implements repository.StaffRepository for tests.
type mockStaffRepo struct {
	GetHospitalBySlugFunc             func(ctx context.Context, slug string) (*model.Hospital, error)
	CreateStaffFunc                   func(ctx context.Context, username, passwordHash, hospitalID string) (string, error)
	GetStaffByUsernameAndHospitalFunc func(ctx context.Context, username, hospitalID string) (*model.Staff, error)
	GetStaffByIDFunc                  func(ctx context.Context, id string) (*model.Staff, error)
}

func (m *mockStaffRepo) GetHospitalBySlug(ctx context.Context, slug string) (*model.Hospital, error) {
	if m.GetHospitalBySlugFunc != nil {
		return m.GetHospitalBySlugFunc(ctx, slug)
	}
	return nil, nil
}

func (m *mockStaffRepo) CreateStaff(ctx context.Context, username, passwordHash, hospitalID string) (string, error) {
	if m.CreateStaffFunc != nil {
		return m.CreateStaffFunc(ctx, username, passwordHash, hospitalID)
	}
	return "", nil
}

func (m *mockStaffRepo) GetStaffByUsernameAndHospital(ctx context.Context, username, hospitalID string) (*model.Staff, error) {
	if m.GetStaffByUsernameAndHospitalFunc != nil {
		return m.GetStaffByUsernameAndHospitalFunc(ctx, username, hospitalID)
	}
	return nil, nil
}

func (m *mockStaffRepo) GetStaffByID(ctx context.Context, id string) (*model.Staff, error) {
	if m.GetStaffByIDFunc != nil {
		return m.GetStaffByIDFunc(ctx, id)
	}
	return nil, nil
}

func TestCreateStaff_Success(t *testing.T) {
	repo := &mockStaffRepo{
		GetHospitalBySlugFunc: func(_ context.Context, slug string) (*model.Hospital, error) {
			return &model.Hospital{ID: "hospital-uuid", Slug: slug}, nil
		},
		CreateStaffFunc: func(_ context.Context, _, _, _ string) (string, error) {
			return "staff-uuid", nil
		},
	}
	svc := NewStaffService(repo, "secret", 900, 604800)

	id, err := svc.CreateStaff(context.Background(), model.CreateStaffRequest{
		Username:     "john",
		Password:     "password123",
		HospitalSlug: "hospital-a",
	})

	require.NoError(t, err)
	assert.Equal(t, "staff-uuid", id)
}

func TestCreateStaff_HospitalNotFound(t *testing.T) {
	repo := &mockStaffRepo{
		GetHospitalBySlugFunc: func(_ context.Context, _ string) (*model.Hospital, error) {
			return nil, nil
		},
	}
	svc := NewStaffService(repo, "secret", 900, 604800)

	_, err := svc.CreateStaff(context.Background(), model.CreateStaffRequest{
		Username:     "john",
		Password:     "password123",
		HospitalSlug: "unknown-hospital",
	})

	assert.ErrorIs(t, err, ErrHospitalNotFound)
}

func TestCreateStaff_DuplicateUsername(t *testing.T) {
	repo := &mockStaffRepo{
		GetHospitalBySlugFunc: func(_ context.Context, slug string) (*model.Hospital, error) {
			return &model.Hospital{ID: "hospital-uuid", Slug: slug}, nil
		},
		CreateStaffFunc: func(_ context.Context, _, _, _ string) (string, error) {
			return "", errors.New("duplicate key value violates unique constraint")
		},
	}
	svc := NewStaffService(repo, "secret", 900, 604800)

	_, err := svc.CreateStaff(context.Background(), model.CreateStaffRequest{
		Username:     "john",
		Password:     "password123",
		HospitalSlug: "hospital-a",
	})

	assert.ErrorIs(t, err, ErrUsernameTaken)
}

func TestLogin_Success(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	repo := &mockStaffRepo{
		GetHospitalBySlugFunc: func(_ context.Context, slug string) (*model.Hospital, error) {
			return &model.Hospital{ID: "hospital-uuid", Slug: slug}, nil
		},
		GetStaffByUsernameAndHospitalFunc: func(_ context.Context, _, _ string) (*model.Staff, error) {
			return &model.Staff{
				ID:           "staff-uuid",
				Username:     "john",
				PasswordHash: string(hash),
				HospitalID:   "hospital-uuid",
			}, nil
		},
	}
	svc := NewStaffService(repo, "secret", 900, 604800)

	tokenResp, refreshToken, err := svc.Login(context.Background(), model.LoginRequest{
		Username:     "john",
		Password:     "password123",
		HospitalSlug: "hospital-a",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, tokenResp.AccessToken)
	assert.Equal(t, "Bearer", tokenResp.TokenType)
	assert.Equal(t, 900, tokenResp.ExpiresIn)
	assert.NotEmpty(t, refreshToken)
}

func TestLogin_HospitalNotFound(t *testing.T) {
	repo := &mockStaffRepo{
		GetHospitalBySlugFunc: func(_ context.Context, _ string) (*model.Hospital, error) {
			return nil, nil
		},
	}
	svc := NewStaffService(repo, "secret", 900, 604800)

	_, _, err := svc.Login(context.Background(), model.LoginRequest{
		Username:     "john",
		Password:     "password123",
		HospitalSlug: "unknown",
	})

	assert.ErrorIs(t, err, ErrHospitalNotFound)
}

func TestLogin_WrongPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct_password"), bcrypt.MinCost)
	repo := &mockStaffRepo{
		GetHospitalBySlugFunc: func(_ context.Context, slug string) (*model.Hospital, error) {
			return &model.Hospital{ID: "hospital-uuid", Slug: slug}, nil
		},
		GetStaffByUsernameAndHospitalFunc: func(_ context.Context, _, _ string) (*model.Staff, error) {
			return &model.Staff{
				ID:           "staff-uuid",
				Username:     "john",
				PasswordHash: string(hash),
				HospitalID:   "hospital-uuid",
			}, nil
		},
	}
	svc := NewStaffService(repo, "secret", 900, 604800)

	_, _, err := svc.Login(context.Background(), model.LoginRequest{
		Username:     "john",
		Password:     "wrong_password",
		HospitalSlug: "hospital-a",
	})

	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestRefreshToken_Success(t *testing.T) {
	refreshToken, err := token.GenerateRefreshToken("staff-uuid", "secret", 604800)
	require.NoError(t, err)

	repo := &mockStaffRepo{
		GetStaffByIDFunc: func(_ context.Context, id string) (*model.Staff, error) {
			return &model.Staff{
				ID:         id,
				Username:   "john",
				HospitalID: "hospital-uuid",
			}, nil
		},
	}
	svc := NewStaffService(repo, "secret", 900, 604800)

	tokenResp, err := svc.RefreshToken(context.Background(), refreshToken)

	require.NoError(t, err)
	assert.NotEmpty(t, tokenResp.AccessToken)
	assert.Equal(t, "Bearer", tokenResp.TokenType)
	assert.Equal(t, 900, tokenResp.ExpiresIn)
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	svc := NewStaffService(&mockStaffRepo{}, "secret", 900, 604800)

	_, err := svc.RefreshToken(context.Background(), "not-a-valid-token")

	assert.ErrorIs(t, err, token.ErrTokenInvalid)
}

func TestRefreshToken_StaffNotFound(t *testing.T) {
	refreshToken, err := token.GenerateRefreshToken("ghost-uuid", "secret", 604800)
	require.NoError(t, err)

	repo := &mockStaffRepo{
		GetStaffByIDFunc: func(_ context.Context, _ string) (*model.Staff, error) {
			return nil, nil
		},
	}
	svc := NewStaffService(repo, "secret", 900, 604800)

	_, err = svc.RefreshToken(context.Background(), refreshToken)

	assert.ErrorIs(t, err, ErrInvalidCredentials)
}
