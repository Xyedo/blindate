CREATE TABLE valid_match_request (match_value VARCHAR(25) PRIMARY KEY);

INSERT INTO
  valid_match_request(match_value)
VALUES
  ('unknown'),
  ('requested'),
  ('accepted');

CREATE TABLE match(
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  request_from UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  request_to UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  request_status VARCHAR(25) NOT NULL REFERENCES valid_match_request(match_value) ON UPDATE CASCADE DEFAULT 'unknown',
  created_at TIMESTAMPTZ NOT NULL,
  accepted_at TIMESTAMPTZ,
  reveal_status VARCHAR(25) NOT NULL REFERENCES valid_match_request(match_value) ON UPDATE CASCADE DEFAULT 'unknown',
  revealed_at TIMESTAMPTZ,
  CONSTRAINT match_request_from_request_to_unique UNIQUE(request_from, request_to)
);