-- Create the players table
 CREATE TABLE IF NOT EXISTS players (
     id SERIAL PRIMARY KEY,
     balance INTEGER NOT NULL,
     status TEXT NOT NULL,
     room_user_id INTEGER NOT NULL,
     game_id INTEGER NOT NULL,
     FOREIGN KEY (room_user_id) REFERENCES room_users(id) ON UPDATE CASCADE ON DELETE CASCADE,
     FOREIGN KEY (game_id) REFERENCES lock_stock_games(id) ON UPDATE CASCADE ON DELETE CASCADE
 );
