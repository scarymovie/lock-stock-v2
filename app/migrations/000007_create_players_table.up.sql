-- Create the players table
 CREATE TABLE IF NOT EXISTS players (
     id SERIAL PRIMARY KEY,
     uid VARCHAR(255) NOT NULL,
     balance INTEGER NOT NULL,
     status TEXT NOT NULL,
     user_id INTEGER NOT NULL,
     game_id INTEGER NOT NULL,
     FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
     FOREIGN KEY (game_id) REFERENCES lock_stock_games(id) ON UPDATE CASCADE ON DELETE CASCADE
 );
