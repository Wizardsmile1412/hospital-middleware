package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wizardsmile1412/hospital-middleware/internal/config"
	"github.com/Wizardsmile1412/hospital-middleware/internal/handler"
	"github.com/Wizardsmile1412/hospital-middleware/internal/model"
	"github.com/Wizardsmile1412/hospital-middleware/internal/routes"
	"github.com/Wizardsmile1412/hospital-middleware/internal/service"
	"github.com/Wizardsmile1412/hospital-middleware/internal/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testJWTSecret = "test-secret"

// --- Mock services ---

type mockStaffService struct {
	CreateStaffFunc  func(ctx context.Context, req model.CreateStaffRequest) (string, error)
	LoginFunc        func(ctx context.Context, req model.LoginRequest) (*model.TokenResponse, string, error)
	RefreshTokenFunc func(ctx context.Context, refreshToken string) (*model.TokenResponse, error)
}

func (m *mockStaffService) CreateStaff(ctx context.Context, req model.CreateStaffRequest) (string, error) {
	if m.CreateStaffFunc != nil {
		return m.CreateStaffFunc(ctx, req)
	}
	return "", nil
}

func (m *mockStaffService) Login(ctx context.Context, req model.LoginRequest) (*model.TokenResponse, string, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(ctx, req)
	}
	return nil, "", nil
}

func (m *mockStaffService) RefreshToken(ctx context.Context, refreshToken string) (*model.TokenResponse, error) {
	if m.RefreshTokenFunc != nil {
		return m.RefreshTokenFunc(ctx, refreshToken)
	}
	return nil, nil
}

type mockPatientService struct {
	SearchPatientsFunc func(ctx context.Context, params model.SearchPatientParams) (*model.SearchPatientResponse, error)
}

func (m *mockPatientService) SearchPatients(ctx context.Context, params model.SearchPatientParams) (*model.SearchPatientResponse, error) {
	if m.SearchPatientsFunc != nil {
		return m.SearchPatientsFunc(ctx, params)
	}
	return &model.SearchPatientResponse{Patients: []model.Patient{}, Total: 0}, nil
}

// --- Test helpers ---

func setupRouter(staffSvc service.StaffService, patientSvc service.PatientService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	cfg := &config.Config{AppEnv: "development", JWTSecret: testJWTSecret}
	staffHandler := handler.NewStaffHandler(staffSvc, cfg)
	patientHandler := handler.NewPatientHandler(patientSvc)
	routes.Register(r, staffHandler, patientHandler, cfg)
	return r
}

func validAccessToken(t *testing.T, hospitalID string) string {
	t.Helper()
	tok, err := token.GenerateAccessToken("staff-uuid", "john", hospitalID, testJWTSecret, 900)
	require.NoError(t, err)
	return tok
}

func expiredAccessToken(t *testing.T) string {
	t.Helper()
	// Negative expiry places ExpiresAt in the past.
	tok, err := token.GenerateAccessToken("staff-uuid", "john", "hospital-uuid", testJWTSecret, -3600)
	require.NoError(t, err)
	return tok
}

// --- Staff handler tests ---

func TestCreateStaff_Success(t *testing.T) {
	svc := &mockStaffService{
		CreateStaffFunc: func(_ context.Context, _ model.CreateStaffRequest) (string, error) {
			return "staff-uuid", nil
		},
	}
	r := setupRouter(svc, &mockPatientService{})

	body, _ := json.Marshal(map[string]string{
		"username":      "john",
		"password":      "password123",
		"hospital_slug": "hospital-a",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/staff/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "staff-uuid", resp["staff_id"])
	assert.Equal(t, "Staff created successfully", resp["message"])
}

func TestCreateStaff_Conflict(t *testing.T) {
	svc := &mockStaffService{
		CreateStaffFunc: func(_ context.Context, _ model.CreateStaffRequest) (string, error) {
			return "", service.ErrUsernameTaken
		},
	}
	r := setupRouter(svc, &mockPatientService{})

	body, _ := json.Marshal(map[string]string{
		"username":      "john",
		"password":      "password123",
		"hospital_slug": "hospital-a",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/staff/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestCreateStaff_HospitalNotFound(t *testing.T) {
	svc := &mockStaffService{
		CreateStaffFunc: func(_ context.Context, _ model.CreateStaffRequest) (string, error) {
			return "", service.ErrHospitalNotFound
		},
	}
	r := setupRouter(svc, &mockPatientService{})

	body, _ := json.Marshal(map[string]string{
		"username":      "john",
		"password":      "password123",
		"hospital_slug": "unknown",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/staff/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestLogin_Success_SetsCookie(t *testing.T) {
	svc := &mockStaffService{
		LoginFunc: func(_ context.Context, _ model.LoginRequest) (*model.TokenResponse, string, error) {
			return &model.TokenResponse{
				AccessToken: "access-token",
				TokenType:   "Bearer",
				ExpiresIn:   900,
			}, "refresh-token", nil
		},
	}
	r := setupRouter(svc, &mockPatientService{})

	body, _ := json.Marshal(map[string]string{
		"username":      "john",
		"password":      "password123",
		"hospital_slug": "hospital-a",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/staff/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp model.TokenResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.NotEmpty(t, resp.AccessToken)
	assert.Equal(t, "Bearer", resp.TokenType)

	var refreshCookieFound bool
	for _, c := range w.Result().Cookies() {
		if c.Name == "refresh_token" {
			refreshCookieFound = true
			assert.True(t, c.HttpOnly)
		}
	}
	assert.True(t, refreshCookieFound, "refresh_token cookie should be set")
}

func TestLogin_Unauthorized(t *testing.T) {
	svc := &mockStaffService{
		LoginFunc: func(_ context.Context, _ model.LoginRequest) (*model.TokenResponse, string, error) {
			return nil, "", service.ErrInvalidCredentials
		},
	}
	r := setupRouter(svc, &mockPatientService{})

	body, _ := json.Marshal(map[string]string{
		"username":      "john",
		"password":      "wrong",
		"hospital_slug": "hospital-a",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/staff/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// --- Patient handler tests ---

func TestSearchPatients_Authorized(t *testing.T) {
	hospitalID := "hospital-uuid"
	svc := &mockPatientService{
		SearchPatientsFunc: func(_ context.Context, params model.SearchPatientParams) (*model.SearchPatientResponse, error) {
			// Verify hospital_id is injected from JWT claims, not query params.
			assert.Equal(t, hospitalID, params.HospitalID)
			return &model.SearchPatientResponse{
				Patients: []model.Patient{{ID: "p1", HospitalID: hospitalID}},
				Total:    1,
			}, nil
		},
	}
	r := setupRouter(&mockStaffService{}, svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/patient/search", nil)
	req.Header.Set("Authorization", "Bearer "+validAccessToken(t, hospitalID))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp model.SearchPatientResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 1, resp.Total)
}

func TestSearchPatients_NoToken(t *testing.T) {
	r := setupRouter(&mockStaffService{}, &mockPatientService{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/patient/search", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSearchPatients_ExpiredToken(t *testing.T) {
	r := setupRouter(&mockStaffService{}, &mockPatientService{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/patient/search", nil)
	req.Header.Set("Authorization", "Bearer "+expiredAccessToken(t))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "token has expired", resp["error"])
}

func TestSearchPatients_HospitalScoped(t *testing.T) {
	// Even if the client sends hospital_id=other-hospital in query params,
	// the service receives hospital_id from JWT claims only.
	hospitalID := "hospital-a"
	svc := &mockPatientService{
		SearchPatientsFunc: func(_ context.Context, params model.SearchPatientParams) (*model.SearchPatientResponse, error) {
			assert.Equal(t, hospitalID, params.HospitalID)
			return &model.SearchPatientResponse{Patients: []model.Patient{}, Total: 0}, nil
		},
	}
	r := setupRouter(&mockStaffService{}, svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/patient/search?hospital_id=other-hospital", nil)
	req.Header.Set("Authorization", "Bearer "+validAccessToken(t, hospitalID))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
