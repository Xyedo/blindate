CREATE TABLE profile_picture(
  id SERIAL PRIMARY KEY 
  user_id UUID NOT NULL REFERENCES users(id),
  selected BOOLEAN NOT NULL DEFAULT FALSE,
  picture_ref CITEXT NOT NULL
);