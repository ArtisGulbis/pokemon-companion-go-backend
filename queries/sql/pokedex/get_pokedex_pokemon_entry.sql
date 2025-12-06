SELECT
    name,
    entry_number,
    pokemon_id,
    pokedex_id
FROM pokedex_pokemon_entries
WHERE pokedex_id = ?
