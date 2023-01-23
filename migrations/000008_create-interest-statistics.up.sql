CREATE TABLE interest_statistics (
  interest_id UUID PRIMARY KEY REFERENCES interests(id) ON DELETE CASCADE,
  hobbie_count SMALLINT NOT NULL DEFAULT 0,
  movie_serie_count SMALLINT NOT NULL DEFAULT 0,
  traveling_count SMALLINT NOT NULL DEFAULT 0,
  sport_count SMALLINT NOT NULL DEFAULT 0,
  CONSTRAINT interest_statistics_interest_id_hobbie_count UNIQUE(interest_id, hobbie_count),
  CONSTRAINT interest_statistics_interest_id_movie_serie_count UNIQUE(interest_id, movie_serie_count),
  CONSTRAINT interest_statistics_interest_id_traveling_count UNIQUE(interest_id, traveling_count),
  CONSTRAINT interest_statistics_interest_id_sport_count UNIQUE(interest_id, sport_count),
  CONSTRAINT interest_statistics_hobbie_count_chk CHECK (hobbie_count BETWEEN 0 AND 10),
  CONSTRAINT interest_statistics_movie_serie_count_chk CHECK (movie_serie_count BETWEEN 0 AND 10),
  CONSTRAINT interest_statistics_traveling_count_chk CHECK (traveling_count BETWEEN 0 AND 10),
  CONSTRAINT interest_statistics_sport_count_chk CHECK (sport_count BETWEEN 0 AND 10)
);


