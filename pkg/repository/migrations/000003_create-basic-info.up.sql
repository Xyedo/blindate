CREATE TABLE basic_info(
  id UUID PRIMARY KEY NOT NULL default uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  gender VARCHAR(50),
  location 
)