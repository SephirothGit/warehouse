ALTER TABLE stock_movements
    ALTER COLUMN created_at TYPE TIMESTAMP USING created_at AT TIME ZONE 'UTC';

ALTER TABLE stock_items
    ALTER COLUMN updated_at TYPE TIMESTAMP USING updated_at AT TIME ZONE 'UTC';

ALTER TABLE refresh_tokens
    ALTER COLUMN expires_at TYPE TIMESTAMP USING expires_at AT TIME ZONE 'UTC',
    ALTER COLUMN issued_at  TYPE TIMESTAMP USING issued_at  AT TIME ZONE 'UTC';