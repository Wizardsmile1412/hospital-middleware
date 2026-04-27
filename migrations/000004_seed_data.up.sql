-- Seed hospitals
INSERT INTO hospitals (id, name, slug, api_url) VALUES
    ('a0000000-0000-0000-0000-000000000001', 'Hospital A', 'hospital-a', 'https://hospital-a.api.co.th'),
    ('b0000000-0000-0000-0000-000000000002', 'Hospital B', 'hospital-b', 'https://hospital-b.api.co.th')
ON CONFLICT DO NOTHING;

-- Seed patients for Hospital A
INSERT INTO patients (hospital_id, national_id, passport_id, first_name_th, last_name_th, first_name_en, last_name_en, date_of_birth, patient_hn, phone_number, email, gender) VALUES
    ('a0000000-0000-0000-0000-000000000001', '1234567890123', NULL, 'สมชาย', 'ใจดี', 'Somchai', 'Jaidee', '1990-01-15', 'HN001234', '0812345678', 'somchai@example.com', 'M'),
    ('a0000000-0000-0000-0000-000000000001', '9876543210987', NULL, 'สมหญิง', 'รักสุข', 'Somying', 'Raksuk', '1985-07-22', 'HN001235', '0898765432', 'somying@example.com', 'F'),
    ('a0000000-0000-0000-0000-000000000001', NULL, 'A12345678', 'จอห์น', 'สมิธ', 'John', 'Smith', '1978-03-10', 'HN001236', '0861234567', 'john.smith@example.com', 'M')
ON CONFLICT DO NOTHING;

-- Seed patients for Hospital B
INSERT INTO patients (hospital_id, national_id, first_name_th, last_name_th, first_name_en, last_name_en, date_of_birth, patient_hn, phone_number, email, gender) VALUES
    ('b0000000-0000-0000-0000-000000000002', '1111111111111', 'วิไล', 'สุขสันต์', 'Wilai', 'Suksan', '1992-11-05', 'HN002001', '0823456789', 'wilai@example.com', 'F')
ON CONFLICT DO NOTHING;
