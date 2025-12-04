SELECT
    p.id,
    p.name,
    p.height,
    p.weight,
    t.name AS type_name,
    pt.slot AS type_slot
FROM pokemon p
LEFT JOIN pokemon_types pt ON p.id = pt.pokemon_id
LEFT JOIN types t ON pt.type_id = t.id
WHERE p.id = ?
ORDER BY pt.slot;