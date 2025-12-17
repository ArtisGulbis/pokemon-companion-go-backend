SELECT
    p.id,
    p.species_id,
    p.name,
    p.is_default,
    p.height,
    p.weight,
    p.base_experience,
    p.hp,
    p.attack,
    p.defense,
    p.special_attack,
    p.special_defense,
    p.speed,
    p.sprite_front_default,
    p.sprite_front_shiny,
    p.sprite_artwork
FROM pokemon p
WHERE p.id = ?;
