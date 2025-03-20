BEGIN;

CREATE TABLE round_player_logs (
                                   id SERIAL PRIMARY KEY,
                                   player_id INTEGER NOT NULL REFERENCES players(id) ON DELETE CASCADE,
                                   round_id INTEGER NOT NULL REFERENCES rounds(id) ON DELETE CASCADE,
                                   number INTEGER NOT NULL,
                                   bets_value INTEGER NOT NULL,
                                   answer INTEGER NULL
);

COMMIT;
