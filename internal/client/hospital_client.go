package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Wizardsmile1412/hospital-middleware/internal/model"
)

// HospitalClient is the interface every hospital adapter must satisfy.
// Adding Hospital B/C only requires a new implementation — no existing code changes (Open/Closed).
type HospitalClient interface {
	SearchByNationalID(ctx context.Context, nationalID string) (*model.Patient, error)
	SearchByPassportID(ctx context.Context, passportID string) (*model.Patient, error)
}

// hospitalAClient calls the real Hospital A external REST API.
type hospitalAClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHospitalAClient(baseURL string) HospitalClient {
	return &hospitalAClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *hospitalAClient) SearchByNationalID(ctx context.Context, nationalID string) (*model.Patient, error) {
	url := fmt.Sprintf("%s/patient/search/%s", c.baseURL, nationalID)
	return c.doRequest(ctx, url)
}

func (c *hospitalAClient) SearchByPassportID(ctx context.Context, passportID string) (*model.Patient, error) {
	url := fmt.Sprintf("%s/patient/search/%s", c.baseURL, passportID)
	return c.doRequest(ctx, url)
}

func (c *hospitalAClient) doRequest(ctx context.Context, url string) (*model.Patient, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("hospital API returned status %d", resp.StatusCode)
	}

	var patient model.Patient
	if err := json.NewDecoder(resp.Body).Decode(&patient); err != nil {
		return nil, err
	}
	return &patient, nil
}
