-- +goose Up
ALTER TABLE admin RENAME COLUMN password TO password_hash;

UPDATE admin
SET password_hash = '$2a$10$m6SFhM5SdeMD34udQNveG.4GLRBoQyry3.2r/GsYNKxY/pyMhxl72'
WHERE login = 'admin';

-- +goose Down
-- A bcrypt hash cannot be converted back to the original plaintext password.
ALTER TABLE admin RENAME COLUMN password_hash TO password;
