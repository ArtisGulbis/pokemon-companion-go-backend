DROP TABLE IF EXISTS version_group;
DROP TABLE IF EXISTS version;
DROP TABLE IF EXISTS pokedexes;
DROP TABLE IF EXISTS version_group_pokedexes;
DROP TABLE IF EXISTS species;
DROP TABLE IF EXISTS pokedex_entries;
DROP TABLE IF EXISTS pokemon;
DROP TABLE IF EXISTS pokemon_types;
DROP TABLE IF EXISTS pokemon_abilities;
DROP TABLE IF EXISTS evolution_chains;
DROP TABLE IF EXISTS evolutions;
DROP TABLE IF EXISTS moves;
DROP TABLE IF EXISTS pokemon_moves;
DROP TABLE IF EXISTS type_effectiveness;
DROP TABLE IF EXISTS abilities;
DROP TABLE IF EXISTS flavor_texts;

-- ============================================================================
-- GAME STRUCTURE TABLES
-- These define which games exist and what Pokemon are available in each
-- ============================================================================

-- Populated from: GET /version-group?limit=100
-- This is your primary "game" selector - Black/White share a version-group
CREATE TABLE version_groups (
    id INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,           -- e.g., "black-white", "sword-shield"
    generation_name TEXT NOT NULL        -- e.g., "generation-v"
);

-- Populated from: GET /version?limit=100
-- Individual games - you show these to users, but use version_group internally
CREATE TABLE versions (
    id INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,           -- e.g., "black", "white", "sword"
    display_name TEXT,                   -- e.g., "Pokemon Black" (from names[].name where language=en)
    version_group_id INTEGER NOT NULL REFERENCES version_groups(id)
);

-- Populated from: GET /pokedex?limit=100
-- Regional Pokedexes - each contains a list of Pokemon
CREATE TABLE pokedexes (
    id INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,           -- e.g., "original-unova", "national"
    region_name TEXT                     -- e.g., "unova", "kanto"
);

-- Populated from: version-group.pokedexes array
-- Links version-groups to their pokedexes (many-to-many)
CREATE TABLE version_group_pokedexes (
    version_group_id INTEGER NOT NULL REFERENCES version_groups(id),
    pokedex_id INTEGER NOT NULL REFERENCES pokedexes(id),
    PRIMARY KEY (version_group_id, pokedex_id)
);

-- ============================================================================
-- POKEMON DATA TABLES
-- Core Pokemon information from pokemon-species and pokemon endpoints
-- ============================================================================

-- Populated from: GET /pokemon-species/{id}
-- The "species" is the conceptual creature - Pikachu the species
CREATE TABLE species (
    id INTEGER PRIMARY KEY,              -- National dex number
    name TEXT UNIQUE NOT NULL,           -- e.g., "pikachu"
    evolution_chain_id INTEGER,          -- Links to evolution_chains table
    gender_rate INTEGER,                 -- -1 = genderless, 0-8 = female ratio
    capture_rate INTEGER,
    base_happiness INTEGER,
    is_baby BOOLEAN DEFAULT FALSE,
    is_legendary BOOLEAN DEFAULT FALSE,
    is_mythical BOOLEAN DEFAULT FALSE,
    growth_rate_name TEXT,               -- e.g., "medium-fast"
    generation_name TEXT                 -- When this species was introduced
);

-- Populated from: pokedex.pokemon_entries array
-- Which species appear in which pokedex (with their regional dex number)
CREATE TABLE pokedex_entries (
    pokedex_id INTEGER NOT NULL REFERENCES pokedexes(id),
    species_id INTEGER NOT NULL REFERENCES species(id),
    entry_number INTEGER NOT NULL,       -- Regional dex number (e.g., Victini is #000 in Unova)
    PRIMARY KEY (pokedex_id, species_id)
);

-- Populated from: GET /pokemon/{id}
-- A "pokemon" is a specific form with stats - Pikachu vs Alolan-Raichu
-- One species can have multiple pokemon (varieties)
CREATE TABLE pokemon (
    id INTEGER PRIMARY KEY,
    species_id INTEGER NOT NULL REFERENCES species(id),
    name TEXT UNIQUE NOT NULL,           -- e.g., "pikachu", "pikachu-gmax", "meowth-alola"
    is_default BOOLEAN NOT NULL,         -- TRUE for the "main" form of each species
    height INTEGER,                      -- In decimeters
    weight INTEGER,                      -- In hectograms
    base_experience INTEGER,
    -- Stats stored directly for easy querying
    hp INTEGER,
    attack INTEGER,
    defense INTEGER,
    special_attack INTEGER,
    special_defense INTEGER,
    speed INTEGER,
    -- Sprite URLs
    sprite_front_default TEXT,
    sprite_front_shiny TEXT,
    sprite_artwork TEXT                  -- official-artwork.front_default
);

-- Populated from: pokemon.types array
CREATE TABLE pokemon_types (
    pokemon_id INTEGER NOT NULL REFERENCES pokemon(id),
    type_name TEXT NOT NULL,             -- e.g., "electric", "fire"
    slot INTEGER NOT NULL,               -- 1 = primary, 2 = secondary
    PRIMARY KEY (pokemon_id, slot)
);

-- Populated from: pokemon.abilities array
CREATE TABLE pokemon_abilities (
    pokemon_id INTEGER NOT NULL REFERENCES pokemon(id),
    ability_name TEXT NOT NULL,          -- e.g., "static", "lightning-rod"
    is_hidden BOOLEAN NOT NULL,          -- Hidden abilities are rarer
    slot INTEGER NOT NULL,
    PRIMARY KEY (pokemon_id, slot)
);

-- ============================================================================
-- EVOLUTION TABLES
-- Evolution chains and requirements
-- ============================================================================

-- Populated from: GET /evolution-chain/{id}
-- Just tracks which chains exist
CREATE TABLE evolution_chains (
    id INTEGER PRIMARY KEY
);

-- Populated from: evolution-chain.chain (recursive structure flattened)
-- Each row = one evolution step
CREATE TABLE evolutions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    chain_id INTEGER NOT NULL REFERENCES evolution_chains(id),
    from_species_id INTEGER NOT NULL REFERENCES species(id),
    to_species_id INTEGER NOT NULL REFERENCES species(id),
    -- Evolution requirements (most are nullable)
    trigger_name TEXT NOT NULL,          -- "level-up", "use-item", "trade", etc.
    min_level INTEGER,                   -- For level-up evolutions
    item_name TEXT,                      -- Evolution stone or held item
    held_item_name TEXT,                 -- Item that must be held
    time_of_day TEXT,                    -- "day" or "night"
    min_happiness INTEGER,               -- Friendship evolutions
    min_affection INTEGER,
    location_name TEXT,                  -- Specific location required
    known_move_name TEXT,                -- Must know this move
    known_move_type_name TEXT,           -- Must know a move of this type
    gender TEXT,                         -- "male" or "female"
    needs_overworld_rain BOOLEAN,
    turn_upside_down BOOLEAN,            -- Inkay → Malamar
    UNIQUE(chain_id, from_species_id, to_species_id, trigger_name)
);

-- ============================================================================
-- MOVE TABLES
-- Moves and how Pokemon learn them (version-specific!)
-- ============================================================================

-- Populated from: GET /move/{id}
CREATE TABLE moves (
    id INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,           -- e.g., "thunderbolt"
    type_name TEXT NOT NULL,             -- e.g., "electric"
    power INTEGER,                       -- NULL for status moves
    accuracy INTEGER,                    -- NULL for moves that can't miss
    pp INTEGER NOT NULL,
    damage_class TEXT NOT NULL,          -- "physical", "special", "status"
    effect_short TEXT,                   -- Brief effect description
    priority INTEGER DEFAULT 0           -- Move priority (-7 to +5)
);

-- Populated from: pokemon.moves array (filtered by version_group_details)
-- THIS IS VERSION-SPECIFIC - same Pokemon learns different moves in different games
CREATE TABLE pokemon_moves (
    pokemon_id INTEGER NOT NULL REFERENCES pokemon(id),
    move_id INTEGER NOT NULL REFERENCES moves(id),
    version_group_id INTEGER NOT NULL REFERENCES version_groups(id),
    learn_method TEXT NOT NULL,          -- "level-up", "machine", "egg", "tutor"
    level_learned_at INTEGER NOT NULL,   -- 0 for non-level-up methods
    PRIMARY KEY (pokemon_id, move_id, version_group_id, learn_method)
);

-- ============================================================================
-- SUPPLEMENTARY TABLES
-- Additional data for display purposes
-- ============================================================================

-- Populated from: GET /type/{name}
CREATE TABLE types (
    name TEXT PRIMARY KEY,
    damage_class TEXT                    -- "physical" or "special" (Gen 1-3 only)
);

-- Populated from: type.damage_relations
CREATE TABLE type_effectiveness (
    attacking_type TEXT NOT NULL REFERENCES types(name),
    defending_type TEXT NOT NULL REFERENCES types(name),
    multiplier REAL NOT NULL,            -- 0, 0.5, 1, or 2
    PRIMARY KEY (attacking_type, defending_type)
);

-- Populated from: GET /ability/{name}
CREATE TABLE abilities (
    name TEXT PRIMARY KEY,
    effect_short TEXT,                   -- Brief description
    effect_full TEXT                     -- Full description
);

-- Populated from: pokemon-species.flavor_text_entries (filtered by version and language)
CREATE TABLE flavor_texts (
    species_id INTEGER NOT NULL REFERENCES species(id),
    version_id INTEGER NOT NULL REFERENCES versions(id),
    flavor_text TEXT NOT NULL,
    PRIMARY KEY (species_id, version_id)
);

CREATE INDEX idx_pokemon_species ON pokemon(species_id);
CREATE INDEX idx_pokemon_default ON pokemon(is_default);
CREATE INDEX idx_pokedex_entries_pokedex ON pokedex_entries(pokedex_id);
CREATE INDEX idx_pokemon_moves_pokemon ON pokemon_moves(pokemon_id);
CREATE INDEX idx_pokemon_moves_version ON pokemon_moves(version_group_id);
CREATE INDEX idx_evolutions_from ON evolutions(from_species_id);
CREATE INDEX idx_evolutions_to ON evolutions(to_species_id);
