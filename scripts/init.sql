CREATE TABLE IF NOT EXISTS players (
    id SERIAL PRIMARY KEY,
    name varchar NOT NULL,   
    wins integer,
    UNIQUE(name)
);