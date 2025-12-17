INSERT OR IGNORE INTO pokemon (
    id,
    species_id,
    name,
    is_default,
    height,
    weight,
    base_experience,
    hp,
    attack,
    defense,
    special_attack,
    special_defense,
    speed,
    sprite_front_default,
    sprite_front_shiny,
    sprite_artwork
 )
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
