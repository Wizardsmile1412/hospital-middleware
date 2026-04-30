CREATE TABLE IF NOT EXISTS patients (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id   UUID NOT NULL REFERENCES hospitals(id) ON DELETE CASCADE,
    national_id   TEXT,
    passport_id   TEXT,
    first_name_th TEXT,
    middle_name_th TEXT,
    last_name_th  TEXT,
    first_name_en TEXT,
    middle_name_en TEXT,
    last_name_en  TEXT,
    date_of_birth DATE,
    patient_hn    TEXT,
    phone_number  TEXT,
    email         TEXT,
    gender        CHAR(1),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_patients_hospital_id ON patients(hospital_id);
CREATE INDEX IF NOT EXISTS idx_patients_national_id  ON patients(national_id);
CREATE INDEX IF NOT EXISTS idx_patients_passport_id  ON patients(passport_id);
