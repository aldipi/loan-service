CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE employees (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE investors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE loans (
    id SERIAL PRIMARY KEY,
    state SMALLINT NOT NULL,
    borrower_id BIGINT NOT NULL,
    principal_amount INT NOT NULL,
    rate DECIMAL(10, 2) NOT NULL,
    roi DECIMAL(10, 2) NOT NULL,
    approval_proof VARCHAR(255),
    approved_by BIGINT,
    agreement_letter VARCHAR(255),
    disbursed_by BIGINT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    approved_at TIMESTAMP,
    invested_at TIMESTAMP,
    disbursed_at TIMESTAMP,
    last_updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_borrower_id ON loans(borrower_id);

CREATE TABLE investments (
    id SERIAL PRIMARY KEY,
    amount INT NOT NULL,
    investor_id BIGINT NOT NULL,
    loan_id BIGINT NOT NULL,
    agreement_letter VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_investor_id ON investments(investor_id);
CREATE INDEX idx_loan_id ON investments(loan_id);

CREATE TABLE loan_products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    rate DECIMAL(10, 2) NOT NULL,
    roi DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);