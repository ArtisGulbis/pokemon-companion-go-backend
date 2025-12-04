DROP TABLE IF EXISTS pokemon;
DROP TABLE IF EXISTS pokedex;
DROP TABLE IF EXISTS pokedex_names;
DROP TABLE IF EXISTS pokedex_descriptions;
DROP TABLE IF EXISTS pokedex_pokemon;
DROP TABLE IF EXISTS pokemon_types;
DROP TABLE IF EXISTS types;

CREATE TABLE pokemon (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    height INTEGER,
    weight INTEGER,
    base_experience INTEGER
);

CREATE TABLE pokedex (
    id INTEGER PRIMARY KEY,
    is_main_series BOOLEAN NOT NULL,
    name TEXT NOT NULL
);

CREATE TABLE pokedex_names (
    language TEXT NOT NULL,
    name TEXT NOT NULL,
    pokedex_id INTEGER NOT NULL,
    FOREIGN KEY (pokedex_id) REFERENCES pokedex(id) ON DELETE CASCADE
);

CREATE TABLE pokedex_descriptions (
    language TEXT NOT NULL,
    description TEXT NOT NULL,
    pokedex_id INTEGER NOT NULL,
    FOREIGN KEY (pokedex_id) REFERENCES pokedex(id) ON DELETE CASCADE
);

CREATE TABLE pokedex_pokemon (
    pokemon_species_id INTEGER NOT NULL,
    entry_number INTEGER NOT NULL,
    pokedex_id INTEGER NOT NULL,
    FOREIGN KEY (pokedex_id) REFERENCES pokedex(id) ON DELETE CASCADE
);

CREATE TABLE pokemon_types (
    pokemon_id INTEGER NOT NULL,    -- Which Pokemon?
    type_id INTEGER NOT NULL,       -- Which Type?
    slot INTEGER NOT NULL,          -- Primary (1) or Secondary (2)?

    FOREIGN KEY (pokemon_id) REFERENCES pokemon(id) ON DELETE CASCADE,
    FOREIGN KEY (type_id) REFERENCES types(id),
    PRIMARY KEY (pokemon_id, slot)
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