-- +migrate Up
CREATE TABLE invoices (
    id VARCHAR(26) NOT NULL PRIMARY KEY,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    image_path VARCHAR(500) NOT NULL,
    extracted_data JSON,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +migrate Down
DROP TABLE IF EXISTS invoices;
