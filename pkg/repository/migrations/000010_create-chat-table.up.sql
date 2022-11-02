CREATE TABLE chats (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
  messages TEXT NOT NULL,
  reply_to UUID REFERENCES chats(id),
  sent_at TIMESTAMPTZ NOT NULL,
  seen_at TIMESTAMPTZ NOT NULL
);