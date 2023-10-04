CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS posts (
  id          uuid          PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
  text        VARCHAR(500)  NOT NULL UNIQUE,
  image       VARCHAR(500)  NOT NULL,
  author_id   uuid          NOT NULL,
  created_at  TIMESTAMP     WITH TIME ZONE DEFAULT NOW(),
  updated_at  TIMESTAMP     WITH TIME ZONE DEFAULT NOW(),
  CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
);