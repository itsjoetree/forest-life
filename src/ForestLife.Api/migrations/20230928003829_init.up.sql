CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
  id          uuid          PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
  username    VARCHAR(255)  NOT NULL UNIQUE,
  name        VARCHAR(255)  NOT NULL,
  email       VARCHAR(255)  NOT NULL UNIQUE,
  password    VARCHAR(255)  NOT NULL,
  created_at  TIMESTAMP     WITH TIME ZONE DEFAULT NOW(),
  updated_at  TIMESTAMP     WITH TIME ZONE DEFAULT NOW()
);