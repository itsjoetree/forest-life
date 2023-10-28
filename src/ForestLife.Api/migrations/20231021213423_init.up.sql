BEGIN;

CREATE TYPE theme AS ENUM ('forest', 'dark', 'standard');

CREATE TABLE IF NOT EXISTS profiles (
  id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
  username VARCHAR(255) NOT NULL UNIQUE,
  nickname VARCHAR(255),
  email VARCHAR(255) NOT NULL UNIQUE,
  theme theme NOT NULL DEFAULT 'standard'
);

ALTER TABLE users DROP COLUMN IF EXISTS username;
ALTER TABLE users DROP COLUMN IF EXISTS email;
ALTER TABLE users DROP COLUMN IF EXISTS name;

ALTER TABLE users ADD COLUMN IF NOT EXISTS profile_id uuid REFERENCES profiles(id) ON DELETE CASCADE;

COMMIT;