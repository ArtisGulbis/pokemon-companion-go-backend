DROP TABLE IF EXISTS pokemon;

CREATE TABLE pokemon (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    national_dex_number INTEGER,
    height INTEGER,
    weight INTEGER,
    base_experience INTEGER
);

CREATE TABLE types (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

INSERT INTO types (name, id) VALUES 
('normal', 1),
('fighting', 2),
('flying', 3),
('poison', 4),
('ground', 5),
('rock', 6),
('bug', 7),
('ghost', 8),
('steel', 9),
('fire', 10),
('water', 11),
('grass', 12),
('electric', 13),
('psychic', 14),
('ice', 15),
('dragon', 16),
('dark', 17),
('fairy', 18),
('stellar', 19);

CREATE TABLE pokemon_types (
    pokemon_id INTEGER NOT NULL,    -- Which Pokemon?
    type_id INTEGER NOT NULL,       -- Which Type?
    slot INTEGER NOT NULL,          -- Primary (1) or Secondary (2)?

    FOREIGN KEY (pokemon_id) REFERENCES pokemon(id) ON DELETE CASCADE,
    FOREIGN KEY (type_id) REFERENCES types(id),
    PRIMARY KEY (pokemon_id, slot)
);