DROP TABLE IF EXISTS pokemon;

CREATE TABLE pokemon (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    national_dex_number INTEGER,
    height INTEGER,
    weight INTEGER,
    base_experience INTEGER,
    p_order INTEGER
);