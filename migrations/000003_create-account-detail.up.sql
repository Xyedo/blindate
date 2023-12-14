CREATE TABLE account_detail(
  account_id TEXT PRIMARY KEY REFERENCES account(id) ON DELETE CASCADE,
  geog geography(POINT, 4326) NOT NULL,
  alias TEXT NOT NULL,
  bio TEXT NOT NULL,
  last_online TIMESTAMPTZ NOT NULL,
  gender TEXT NOT NULL,
  from_loc TEXT,
  height SMALLINT,
  education_level TEXT,
  drinking TEXT,
  smoking TEXT,
  relationship_pref TEXT,
  looking_for TEXT,
  zodiac TEXT,
  kids SMALLINT,
  work VARCHAR(50),
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  version BIGINT NOT NULL DEFAULT 0
);

CREATE INDEX ON account_detail USING GIST(geog);