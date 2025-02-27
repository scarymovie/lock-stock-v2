-- Create the rounds table
CREATE TABLE IF NOT EXISTS rounds (
    id SERIAL PRIMARY KEY,
    uid VARCHAR(255) NOT NULL,
    number INTEGER NOT NULL,
    buy_in INTEGER NOT NULL,
    game_id INTEGER NOT NULL,
    pot INTEGER NOT NULL,
    FOREIGN KEY (game_id) REFERENCES lock_stock_games(id) ON UPDATE CASCADE ON DELETE CASCADE
);