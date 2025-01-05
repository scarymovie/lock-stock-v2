CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       uid VARCHAR(255) UNIQUE NOT NULL DEFAULT md5(random()::text)
);