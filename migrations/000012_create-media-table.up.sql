CREATE TABLE valid_media_type (media_type VARCHAR(25) PRIMARY KEY);

INSERT INTO
  valid_media_type(media_type)
VALUES
  ('application/ogg'),
  ('audio/mpeg');

CREATE TABLE media(
  chat_id UUID PRIMARY KEY REFERENCES chats(id) ON DELETE CASCADE,
  blob_link CITEXT NOT NULL,
  media_type VARCHAR(25) NOT NULL REFERENCES valid_media_type(media_type) ON UPDATE CASCADE
);