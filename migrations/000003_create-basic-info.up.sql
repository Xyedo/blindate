-- CREATE TABLE valid_gender (
--   gender VARCHAR(25) NOT NULL PRIMARY KEY UNIQUE
-- );

-- INSERT INTO
--   valid_gender(gender)
-- VALUES
--   ('Female'),
--   ('Male'),
--   ('Other');

-- CREATE TABLE valid_education_level (
--   education VARCHAR(50) NOT NULL PRIMARY KEY UNIQUE
-- );

-- INSERT INTO
--   valid_education_level(education)
-- VALUES
--   ('Less than high school diploma'),
--   ('High school'),
--   ('Some college, no degree'),
--   ('Assosiate''s Degree'),
--   ('Bachelor''s Degree'),
--   ('Master''s Degree'),
--   ('Professional Degree'),
--   ('Doctorate Degree');

-- CREATE TABLE drinking_level(
--   drinking VARCHAR(50) NOT NULL PRIMARY KEY UNIQUE
-- );

-- INSERT INTO
--   drinking_level(drinking)
-- VALUES
--   ('Never'),
--   ('Ocassionally'),
--   ('Once a week'),
--   ('More than 2/3 times a week'),
--   ('Every day');

-- CREATE TABLE smoking_level(
--   smoking VARCHAR(50) NOT NULL PRIMARY KEY UNIQUE
-- );

-- INSERT INTO
--   smoking_level(smoking)
-- VALUES
--   ('Never'),
--   ('Ocassionally'),
--   ('Once a week'),
--   ('More than 2/3 times a week'),
--   ('Every day');

-- CREATE TABLE relationship_preferences(
--   preferences VARCHAR(50) NOT NULL PRIMARY KEY UNIQUE
-- );

-- INSERT INTO
--   relationship_preferences(preferences)
-- VALUES
--   ('One night Stand'),
--   ('Having fun'),
--   ('Serious');

-- CREATE TABLE zodiac_lookups (
--   zodiac VARCHAR(50) PRIMARY KEY NOT NULL UNIQUE
-- );

-- INSERT INTO
--   zodiac_lookups(zodiac)
-- VALUES
--   ('Aries'),
--   ('Taurus'),
--   ('Gemini'),
--   ('Cancer'),
--   ('Leo'),
--   ('Virgo'),
--   ('Libra'),
--   ('Scorpio'),
--   ('Sagittarius'),
--   ('Capricorn'),
--   ('Aquarius'),
--   ('Pisces');

CREATE TABLE basic_info(
  user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
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