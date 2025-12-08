SELECT
    p.id,
    p.name,
    p.height,
    p.weight,
    p.base_experience,
    pt.type_name,
    pt.slot AS type_slot
FROM pokemon p
LEFT JOIN pokemon_types pt ON p.id = pt.pokemon_id
WHERE p.id = ?
ORDER BY pt.slot;