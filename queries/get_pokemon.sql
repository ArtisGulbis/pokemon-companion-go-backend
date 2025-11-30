SELECT
    p.id,
    p.name,
    p.height,
    p.weight,
    t.name AS type_name,
    pt.slot AS type_slot
FROM pokemon p
JOIN pokemon_types pt ON p.id = pt.pokemon_id
JOIN types t ON pt.type_id = t.id
WHERE p.id = ?
ORDER BY pt.slot;