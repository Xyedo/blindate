CREATE TABLE locations (
  user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  geog geography(POINT, 4326) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  version bigint NOT NULL DEFAULT 0
);

CREATE INDEX ON locations USING GIST(geog);