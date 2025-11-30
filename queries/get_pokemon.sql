-- Retrieve a single Pokemon by ID
SELECT id, name, height, weight
FROM pokemon
WHERE id = ?;
