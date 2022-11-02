CREATE TABLE conversations (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  from_id UUID NOT NULL REFERENCES users(id),
  to_id UUID NOT NULL REFERENCES users(id),
  chat_rows INTEGER NOT NULL DEFAULT 0,
  day_pass INTEGER NOT NULL DEFAULT 0,
  CONSTRAINT conversations_from_id_to_id_unique UNIQUE(from_id, to_id)
);