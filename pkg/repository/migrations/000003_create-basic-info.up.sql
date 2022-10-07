CREATE TABLE valid_gender (
  gender VARCHAR(25) NOT NULL PRIMARY KEY UNIQUE
);

INSERT INTO
  valid_gender(gender)
VALUES
  ('Female'),
  ('Male'),
  ('Other');

CREATE TABLE valid_education_level (
  education VARCHAR(50) NOT NULL PRIMARY KEY UNIQUE
);

INSERT INTO
  valid_education_level(education)
VALUES
  ('Less than high school diploma'),
  ('High school'),
  ('Some college, no degree'),
  ('Assosiate''s Degree'),
  ('Bachelor''s Degree'),
  ('Master''s Degree'),
  ('Professional Degree'),
  ('Doctorate Degree');

CREATE TABLE drinking_level(
  drinking VARCHAR(50) NOT NULL PRIMARY KEY UNIQUE
);

INSERT INTO
  drinking_level(drinking)
VALUES
  ('Never'),
  ('Ocassionally'),
  ('Once a week'),
  ('More than 2/3 times a week'),
  ('Every day');

CREATE TABLE smoking_level(
  smoking VARCHAR(50) NOT NULL PRIMARY KEY UNIQUE
);

INSERT INTO
  smoking_level(smoking)
VALUES
  ('Never'),
  ('Ocassionally'),
  ('Once a week'),
  ('More than 2/3 times a week'),
  ('Every day');

CREATE TABLE relationship_preferences(
  preferences VARCHAR(50) NOT NULL PRIMARY KEY UNIQUE
);

INSERT INTO
  relationship_preferences(preferences)
VALUES
  ('One night Stand'),
  ('Having fun'),
  ('Serious');

CREATE TABLE zodiac_lookups (
  zodiac VARCHAR(50) PRIMARY KEY NOT NULL UNIQUE
);

INSERT INTO
  zodiac_lookups(zodiac)
VALUES
  ('Aries'),
  ('Taurus'),
  ('Gemini'),
  ('Cancer'),
  ('Leo'),
  ('Virgo'),
  ('Libra'),
  ('Scorpio'),
  ('Sagittarius'),
  ('Capricorn'),
  ('Aquarius'),
  ('Pisces');

CREATE TABLE basic_info(
  user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  gender VARCHAR(25) NOT NULL REFERENCES valid_gender(gender) ON UPDATE CASCADE,
  from_loc VARCHAR(100),
  height SMALLINT,
  education_level VARCHAR(50) REFERENCES valid_education_level(education) ON UPDATE CASCADE,
  drinking VARCHAR(50) REFERENCES drinking_level(drinking) ON UPDATE CASCADE,
  smoking VARCHAR(50) REFERENCES smoking_level(smoking) ON UPDATE CASCADE,
  relationship_pref VARCHAR(50) REFERENCES relationship_preferences(preferences) ON UPDATE CASCADE,
  looking_for VARCHAR(25) NOT NULL REFERENCES valid_gender(gender) ON UPDATE CASCADE,
  zodiac VARCHAR(50) REFERENCES zodiac_lookups(zodiac) ON UPDATE CASCADE,
  kids SMALLINT,
  work VARCHAR(50),
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);