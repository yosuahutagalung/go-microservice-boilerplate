-- +goose Up
CREATE TABLE service_boilerplate_outbox_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    topic VARCHAR(255) NOT NULL,
    
    -- BYTEA is required because we are storing raw Protobuf binary bytes, not JSON strings
    payload BYTEA NOT NULL, 
    
    -- Status will toggle from 'PENDING' to 'PUBLISHED'
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- CRITICAL: The Relay Worker constantly queries WHERE status = 'PENDING' ORDER BY created_at.
-- This index ensures that query takes microseconds instead of scanning the whole table.
CREATE INDEX idx_service_boilerplate_outbox_events_status_created_at ON service_boilerplate_outbox_events (status, created_at);

-- +goose Down
DROP TABLE IF EXISTS service_boilerplate_outbox_events;