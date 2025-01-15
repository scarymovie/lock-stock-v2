ALTER TABLE users
    ADD COLUMN name VARCHAR(255) not null DEFAULT md5(random()::text || now()::text);
