CREATE TABLE interests (
  id UUID PRIMARY KEY default uuid_generate_v4(),
  user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
  bio TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE hobbies (
  id UUID PRIMARY KEY default  uuid_generate_v4(),
  interest_id UUID NOT NULL REFERENCES interests(id) ON DELETE CASCADE,
  hobbie CITEXT NOT NULL,
  CONSTRAINT hobbies_interest_id_hobbie_unique UNIQUE(interest_id, hobbie)   
);

CREATE TABLE movie_series (
  id UUID PRIMARY KEY default  uuid_generate_v4(),
  interest_id UUID NOT NULL REFERENCES interests(id) ON DELETE CASCADE,
  movie_serie CITEXT NOT NULL,
  CONSTRAINT movie_series_interest_id_movie_serie_unique UNIQUE(interest_id, movie_serie)   
);

CREATE TABLE traveling (
  id UUID PRIMARY KEY default  uuid_generate_v4(),
  interest_id UUID NOT NULL REFERENCES interests(id) ON DELETE CASCADE,
  travel CITEXT NOT NULL,
  CONSTRAINT traveling_interest_id_travel_unique UNIQUE(interest_id, travel)  
);

CREATE TABLE sports (
  id UUID PRIMARY KEY default  uuid_generate_v4(),
  interest_id UUID NOT NULL REFERENCES interests(id) ON DELETE CASCADE,
  sport CITEXT NOT NULL,
  CONSTRAINT sports_interest_id_sport_unique UNIQUE(interest_id, sport)  
);