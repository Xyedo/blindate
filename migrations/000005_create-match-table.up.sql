CREATE TABLE match(
  id TEXT PRIMARY KEY,
  request_from TEXT NOT NULL REFERENCES account(id) ON DELETE CASCADE,
  request_to TEXT NOT NULL REFERENCES account(id) ON DELETE CASCADE,
  request_status VARCHAR(50) NOT NULL,
  accepted_at TIMESTAMPTZ,
  reveal_status VARCHAR(50),
  revealed_declined_count INT,
  revealed_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  updated_by TEXT REFERENCES account(id) ON DELETE CASCADE,
  version BIGINT NOT NULL DEFAULT 0,
  CONSTRAINT match_request_from_request_to_unique UNIQUE (request_from, request_to) 
);