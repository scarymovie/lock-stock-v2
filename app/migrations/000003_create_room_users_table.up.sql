CREATE TABLE room_users (
                            id SERIAL PRIMARY KEY,
                            room_id INT NOT NULL,
                            user_id INT NOT NULL,
                            FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
                            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                            UNIQUE (room_id, user_id) -- Уникальная пара room_id и user_id
);
