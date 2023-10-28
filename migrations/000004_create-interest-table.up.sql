CREATE TABLE hobbies (
  id UUID PRIMARY KEY default  uuid_generate_v4(),
  account_id CITEXT NOT NULL REFERENCES account(id) ON DELETE CASCADE,
  hobbie CITEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  version bigint NOT NULL DEFAULT 0,
  CONSTRAINT hobbies_account_id_hobbie_unique UNIQUE (account_id, hobbie)   
);

CREATE TABLE movie_series (
  id UUID PRIMARY KEY default  uuid_generate_v4(),
  account_id CITEXT NOT NULL REFERENCES account(id) ON DELETE CASCADE,
  movie_serie CITEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  version bigint NOT NULL DEFAULT 0,
  CONSTRAINT movie_series_account_id_movie_serie_unique UNIQUE (account_id, movie_serie)   
);

CREATE TABLE traveling (
  id UUID PRIMARY KEY default  uuid_generate_v4(),
  account_id CITEXT NOT NULL REFERENCES account(id) ON DELETE CASCADE,
  travel CITEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  version bigint NOT NULL DEFAULT 0,
  CONSTRAINT traveling_account_id_travel_unique UNIQUE (account_id, travel)  
);

CREATE TABLE sports (
  id UUID PRIMARY KEY default  uuid_generate_v4(),
  account_id CITEXT NOT NULL REFERENCES account(id) ON DELETE CASCADE,
  sport CITEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  version bigint NOT NULL DEFAULT 0,
  CONSTRAINT sports_account_id_sport_unique UNIQUE (account_id, sport)  
);