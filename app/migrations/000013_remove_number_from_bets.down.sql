BEGIN;
ALTER TABLE bets
    ADD COLUMN number VARCHAR(255) not null DEFAULT md5(random()::text || now()::text);

commit;