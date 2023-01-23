CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  full_name TEXT NOT NULL,
  alias TEXT NOT NULL,
  email citext UNIQUE NOT NULL,
  password TEXT NOT NULL,
  dob DATE NOT NULL,
  active BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);