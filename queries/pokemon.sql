-- Insert or update a Pokemon
-- Uses SQLite's ON CONFLICT clause to handle duplicate IDs
-- The excluded.* syntax refers to the values we tried to insert
INSERT INTO pokemon (id, name, height, weight)
VALUES (?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    name = excluded.name,
    height = excluded.height,
    weight = excluded.weight;
