CREATE TABLE onlines (
  user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  last_online TIMESTAMPTZ NOT NULL,
  is_online BOOLEAN NOT NULL
);