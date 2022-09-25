CREATE TABLE users (
  user_id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
  email citext UNIQUE NOT NULL,
  password TEXT NOT NULL,
  dob DATE NOT NULL,
  active BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);