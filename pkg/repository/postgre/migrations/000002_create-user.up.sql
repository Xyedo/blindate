CREATE TABLE users (
  user_id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
  email citext UNIQUE NOT NULL,
  dob DATE NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);