package repository

import (
	"context"
	"errors"

	"github.com/Wizardsmile1412/hospital-middleware/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StaffRepository interface {
	GetHospitalBySlug(ctx context.Context, slug string) (*model.Hospital, error)
	CreateStaff(ctx context.Context, username, passwordHash, hospitalID string) (string, error)
	GetStaffByUsernameAndHospital(ctx context.Context, username, hospitalID string) (*model.Staff, error)
}

type staffRepo struct {
	db *pgxpool.Pool
}

func NewStaffRepository(db *pgxpool.Pool) StaffRepository {
	return &staffRepo{db: db}
}

func (r *staffRepo) GetHospitalBySlug(ctx context.Context, slug string) (*model.Hospital, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, name, slug, api_url, created_at FROM hospitals WHERE slug = $1`, slug)

	var h model.Hospital
	err := row.Scan(&h.ID, &h.Name, &h.Slug, &h.APIURL, &h.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (r *staffRepo) CreateStaff(ctx context.Context, username, passwordHash, hospitalID string) (string, error) {
	var id string
	err := r.db.QueryRow(ctx,
		`INSERT INTO staff (username, password_hash, hospital_id)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		username, passwordHash, hospitalID,
	).Scan(&id)
	return id, err
}

func (r *staffRepo) GetStaffByUsernameAndHospital(ctx context.Context, username, hospitalID string) (*model.Staff, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, username, password_hash, hospital_id, created_at
		 FROM staff WHERE username = $1 AND hospital_id = $2`,
		username, hospitalID)

	var s model.Staff
	err := row.Scan(&s.ID, &s.Username, &s.PasswordHash, &s.HospitalID, &s.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}
