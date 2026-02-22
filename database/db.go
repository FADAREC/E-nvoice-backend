package database

import (
    "database/sql"
    "log"
    _ "github.com/lib/pq"
)

var DB *sql.DB

func Connect(databaseURL string) error {
    var err error
    DB, err = sql.Open("postgres", databaseURL)
    if err != nil {
        return err
    }

    if err = DB.Ping(); err != nil {
        return err
    }

    log.Println("✅ Database connected")
    return nil
}

func InitSchema() error {
    schema := `
    CREATE TABLE IF NOT EXISTS users (
        id UUID PRIMARY KEY,
        email VARCHAR(255) UNIQUE NOT NULL,
        password_hash VARCHAR(255) NOT NULL,
        business_name VARCHAR(255),
        created_at TIMESTAMP DEFAULT NOW(),
        updated_at TIMESTAMP DEFAULT NOW()
    );

    CREATE TABLE IF NOT EXISTS mutation_log (
        id BIGSERIAL PRIMARY KEY,
        user_id UUID NOT NULL REFERENCES users(id),
        entity_type VARCHAR(50) NOT NULL,
        entity_id VARCHAR(255) NOT NULL,
        operation VARCHAR(20) NOT NULL,
        payload JSONB NOT NULL,
        device_id VARCHAR(255) NOT NULL,
        client_timestamp TIMESTAMP NOT NULL,
        server_timestamp TIMESTAMP DEFAULT NOW(),
        CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );

    CREATE INDEX IF NOT EXISTS idx_mutation_user_timestamp 
        ON mutation_log(user_id, server_timestamp);
    CREATE INDEX IF NOT EXISTS idx_mutation_entity 
        ON mutation_log(entity_type, entity_id);

    CREATE TABLE IF NOT EXISTS invoices (
        id UUID PRIMARY KEY,
        user_id UUID NOT NULL REFERENCES users(id),
        invoice_number VARCHAR(50) NOT NULL,
        branch_id UUID NOT NULL,
        customer_id UUID NOT NULL,
        status VARCHAR(20) NOT NULL,
        subtotal DECIMAL(10,2) NOT NULL,
        tax DECIMAL(10,2) NOT NULL,
        discount DECIMAL(10,2) NOT NULL,
        total DECIMAL(10,2) NOT NULL,
        notes TEXT,
        due_date TIMESTAMP,
        paid_at TIMESTAMP,
        payment_method VARCHAR(50),
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        device_id VARCHAR(255) NOT NULL,
        deleted BOOLEAN DEFAULT FALSE
    );

    CREATE INDEX IF NOT EXISTS idx_invoice_user ON invoices(user_id);
    CREATE INDEX IF NOT EXISTS idx_invoice_status ON invoices(status);

    CREATE TABLE IF NOT EXISTS customers (
        id UUID PRIMARY KEY,
        user_id UUID NOT NULL REFERENCES users(id),
        name VARCHAR(255) NOT NULL,
        email VARCHAR(255),
        phone VARCHAR(50),
        address TEXT,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        deleted BOOLEAN DEFAULT FALSE
    );

    CREATE INDEX IF NOT EXISTS idx_customer_user ON customers(user_id);

    CREATE TABLE IF NOT EXISTS line_items (
        id UUID PRIMARY KEY,
        invoice_id UUID NOT NULL,
        name VARCHAR(255) NOT NULL,
        description TEXT,
        quantity DECIMAL(10,2) NOT NULL,
        unit_price DECIMAL(10,2) NOT NULL,
        total DECIMAL(10,2) NOT NULL
    );

    CREATE INDEX IF NOT EXISTS idx_line_items_invoice ON line_items(invoice_id);

    CREATE TABLE IF NOT EXISTS receipts (
        id UUID PRIMARY KEY,
        user_id UUID NOT NULL REFERENCES users(id),
        invoice_id UUID NOT NULL,
        receipt_number VARCHAR(50) NOT NULL,
        amount DECIMAL(10,2) NOT NULL,
        payment_date TIMESTAMP NOT NULL,
        payment_method VARCHAR(50) NOT NULL,
        created_at TIMESTAMP NOT NULL
    );

    CREATE INDEX IF NOT EXISTS idx_receipt_user ON receipts(user_id);
    CREATE INDEX IF NOT EXISTS idx_receipt_invoice ON receipts(invoice_id);
    `

    _, err := DB.Exec(schema)
    if err != nil {
        return err
    }

    log.Println("✅ Database schema initialized")
    return nil
}