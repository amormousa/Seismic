-- Adds display_name, separate from username.
-- username: unique, used in URLs, lowercase, no spaces
-- display_name: shown on profile, can have spaces/caps, not unique
ALTER TABLE users
    ADD COLUMN display_name VARCHAR(50);