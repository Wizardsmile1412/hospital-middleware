package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Wizardsmile1412/hospital-middleware/internal/model"
	"github.com/Wizardsmile1412/hospital-middleware/internal/repository"
	"github.com/Wizardsmile1412/hospital-middleware/internal/token"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrHospitalNotFound  = errors.New("hospital not found")
	ErrUsernameTaken     = errors.New("username already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type StaffService interface {
	CreateStaff(ctx context.Context, req model.CreateStaffRequest) (string, error)
	Login(ctx context.Context, req model.LoginRequest) (*model.TokenResponse, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (*model.TokenResponse, error)
}

type staffService struct {
	repo             repository.StaffRepository
	jwtSecret        string
	accessExpirySec  int
	refreshExpiryeSec int
}

func NewStaffService(repo repository.StaffRepository, jwtSecret string, accessExpiry, refreshExpiry int) StaffService {
	return &staffService{
		repo:              repo,
		jwtSecret:         jwtSecret,
		accessExpirySec:   accessExpiry,
		refreshExpiryeSec: refreshExpiry,
	}
}

func (s *staffService) CreateStaff(ctx context.Context, req model.CreateStaffRequest) (string, error) {
	hospital, err := s.repo.GetHospitalBySlug(ctx, req.HospitalSlug)
	if err != nil {
		return "", err
	}
	if hospital == nil {
		return "", ErrHospitalNotFound
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	staffID, err := s.repo.CreateStaff(ctx, req.Username, string(hash), hospital.ID)
	if err != nil {
		if isDuplicateKeyError(err) {
			return "", ErrUsernameTaken
		}
		return "", err
	}
	return staffID, nil
}

func (s *staffService) Login(ctx context.Context, req model.LoginRequest) (*model.TokenResponse, string, error) {
	hospital, err := s.repo.GetHospitalBySlug(ctx, req.HospitalSlug)
	if err != nil {
		return nil, "", err
	}
	if hospital == nil {
		return nil, "", ErrHospitalNotFound
	}

	staff, err := s.repo.GetStaffByUsernameAndHospital(ctx, req.Username, hospital.ID)
	if err != nil {
		return nil, "", err
	}
	if staff == nil {
		return nil, "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(req.Password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	accessToken, err := token.GenerateAccessToken(staff.ID, staff.Username, staff.HospitalID, s.jwtSecret, s.accessExpirySec)
	if err != nil {
		return nil, "", err
	}

	refreshToken, err := token.GenerateRefreshToken(staff.ID, s.jwtSecret, s.refreshExpiryeSec)
	if err != nil {
		return nil, "", err
	}

	return &model.TokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   s.accessExpirySec,
	}, refreshToken, nil
}

func (s *staffService) RefreshToken(ctx context.Context, refreshToken string) (*model.TokenResponse, error) {
	staffID, err := token.ParseRefreshToken(refreshToken, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	staff, err := s.repo.GetStaffByID(ctx, staffID)
	if err != nil {
		return nil, err
	}
	if staff == nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := token.GenerateAccessToken(staff.ID, staff.Username, staff.HospitalID, s.jwtSecret, s.accessExpirySec)
	if err != nil {
		return nil, err
	}

	return &model.TokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   s.accessExpirySec,
	}, nil
}

// isDuplicateKeyError detects PostgreSQL unique constraint violations.
func isDuplicateKeyError(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint"))
}
