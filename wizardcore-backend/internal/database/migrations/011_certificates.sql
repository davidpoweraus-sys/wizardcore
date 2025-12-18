-- Create certificates table
CREATE TABLE IF NOT EXISTS certificates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    pathway_id UUID REFERENCES pathways(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    certificate_number VARCHAR(100) UNIQUE NOT NULL,
    verification_url TEXT,
    download_url TEXT,
    is_verified BOOLEAN DEFAULT false,
    verified_at TIMESTAMP,
    issued_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_certificates_user_id ON certificates(user_id);
CREATE INDEX idx_certificates_certificate_number ON certificates(certificate_number);