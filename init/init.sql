CREATE TABLE IF NOT EXISTS service (
                                       id SERIAL PRIMARY KEY,
                                       user_id TEXT NOT NULL,
                                       name TEXT NOT NULL,
                                       login TEXT NOT NULL,
                                       password TEXT NOT NULL,
                                       hash TEXT NOT NULL,
                                       expiration_time TIMESTAMP NOT NULL
);