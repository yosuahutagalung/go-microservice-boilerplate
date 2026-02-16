-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE greeters (
    id         VARCHAR(36) PRIMARY KEY,
    hello      VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE greeters;