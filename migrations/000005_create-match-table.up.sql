CREATE TABLE match(
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  request_from CITEXT NOT NULL REFERENCES account(id) ON DELETE CASCADE,
  request_to CITEXT NOT NULL REFERENCES account(id) ON DELETE CASCADE,
  request_status VARCHAR(25),
  accepted_at TIMESTAMPTZ,
  reveal_status VARCHAR(25),
  revealed_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  version BIGINT NOT NULL DEFAULT 0,
  CONSTRAINT match_request_from_request_to_unique UNIQUE (request_from, request_to) 
);