-- Add up migration script here
CREATE TABLE sessions (
    id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    username VARCHAR(255) NOT NULL UNIQUE,
    expiry  TIMESTAMP
);