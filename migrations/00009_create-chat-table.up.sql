CREATE TABLE chats (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  conversation_id UUID NOT NULL REFERENCES conversations(match_id) ON DELETE CASCADE,
  author CITEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  messages TEXT NOT NULL,
  reply_to UUID REFERENCES chats(id) ON DELETE SET NULL,
  sent_at TIMESTAMPTZ NOT NULL,
  seen_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  version BIGINT NOT NULL DEFAULT 0 
);

CREATE INDEX sent_at_idx ON chats(id, sent_at);