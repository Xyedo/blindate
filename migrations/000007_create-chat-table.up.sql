CREATE TABLE chat (
  id TEXT PRIMARY KEY,
  conversation_id TEXT NOT NULL REFERENCES conversations(match_id) ON DELETE CASCADE,
  author TEXT NOT NULL REFERENCES account(id) ON DELETE CASCADE,
  messages TEXT NOT NULL,
  reply_to TEXT REFERENCES chat(id) ON DELETE SET NULL,
  sent_at TIMESTAMPTZ NOT NULL,
  seen_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ NOT NULL,
  version BIGINT NOT NULL DEFAULT 0
);

CREATE INDEX ON chat (conversation_id, sent_at DESC, id);



