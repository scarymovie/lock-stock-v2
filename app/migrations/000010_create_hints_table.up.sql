-- Create the hints table
CREATE TABLE IF NOT EXISTS hints (
    id SERIAL PRIMARY KEY,
    text TEXT NOT NULL,
    round_id INTEGER,
    FOREIGN KEY (round_id) REFERENCES rounds(id) ON UPDATE CASCADE ON DELETE SET NULL
);