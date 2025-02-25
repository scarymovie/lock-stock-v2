-- Create the rounds table
CREATE TABLE IF NOT EXISTS rounds (
    id SERIAL PRIMARY KEY,
    round_number INTEGER NOT NULL,
    question TEXT NOT NULL,
    price INTEGER NOT NULL,
    game_id INTEGER NOT NULL,
    FOREIGN KEY (game_id) REFERENCES lock_stock_games(id) ON UPDATE CASCADE ON DELETE CASCADE
);