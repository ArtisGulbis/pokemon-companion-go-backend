# Special Game Versions Handling

## Overview

Some Pokemon games like **Pokemon Colosseum** and **Pokemon XD: Gale of Darkness** don't have traditional Pokedexes in the PokeAPI database. These games require special handling with hardcoded Pokemon lists.

## How It Works

### 1. Pokemon Lists (`models/special_games.go`)

This file contains the hardcoded Pokemon IDs for special game versions:

- **ColosseumPokemonIDs**: 50 Pokemon available in Pokemon Colosseum
- **XDPokemonIDs**: Pokemon available in Pokemon XD (needs to be populated)

### 2. Automatic Detection (`services/game_syncer.go`)

The `SyncGame` function automatically detects special game versions:

```go
if specialPokemonIDs := models.GetSpecialGamePokemon(version.Name); specialPokemonIDs != nil {
    log.Printf("Special game detected: %s - Processing %d Pokemon...", version.Name, len(specialPokemonIDs))
    return g.syncSpecialGamePokemon(specialPokemonIDs)
}
```

When syncing a game version:
1. Checks if the version name matches a special game ("colosseum", "xd", etc.)
2. If found, uses the hardcoded Pokemon list instead of fetching from Pokedexes
3. Syncs each Pokemon with its species, types, moves, and abilities

### 3. Version Names

The version names must match exactly:
- **Colosseum**: `"colosseum"`
- **XD**: `"xd"`

These names come from the PokeAPI and are also used as keys in the IGDB game mapping.

## Adding Pokemon XD Support

To add support for Pokemon XD, you need to populate the `XDPokemonIDs` slice in `models/special_games.go`:

```go
var XDPokemonIDs = []int{
    // Add all Pokemon IDs available in XD here
    // Example: 1, 2, 3, 4, 5, etc.
}
```

### Finding XD Pokemon IDs

1. **Manual Method**: Look up which Pokemon are available in XD and note their National Dex numbers
2. **Bulbapedia**: Check the Pokemon XD page on Bulbapedia for a complete list
3. **Database Query**: If you have XD data in your database from another source

## How Special Game Syncing Works

When a special game is detected:

1. **Insert Version & Version Group**: Normal flow
2. **Create Virtual Pokedex**: Creates a fake pokedex with ID = 1000 + version_group_id
   - Name format: `"{version-name}-pokedex"` (e.g., "colosseum-pokedex")
   - Region: "unknown"
3. **Link Pokedex to Version Group**: Inserts into `version_group_pokedexes` table
4. **Use Hardcoded List**: Iterate through the Pokemon ID list
5. **For Each Pokemon**:
   - Sync species data
   - Sync Pokemon data (stats, sprites, etc.)
   - Insert types
   - Sync and link moves
   - Insert abilities
   - **Insert pokedex entry** with sequential entry number

### Benefits of Caching

All the optimization caching still applies:
- Species are cached (many Pokemon share species)
- Pokemon are cached (if they appear in multiple games)
- Moves are cached (significant reduction in API calls)

## Adding More Special Games

To add support for another special game version:

1. **Add Pokemon List** in `models/special_games.go`:
   ```go
   var MyGamePokemonIDs = []int{
       // Pokemon IDs
   }
   ```

2. **Add to Switch Statement** in `GetSpecialGamePokemon`:
   ```go
   case "my-game-version-name":
       return MyGamePokemonIDs
   ```

3. **Ensure Version Name Matches**: Check that the version name in PokeAPI matches your case statement

## Example: Colosseum Sync

When syncing Colosseum:
```
Syncing game colosseum (X/Y)...
Inserting version group...
Inserting version...
Special game detected: colosseum - Processing 50 Pokemon...
Creating virtual pokedex 1016 (colosseum-pokedex)...
  [1/50] Processing Pokemon 153...
    Inserting 1 types for pokemon 153 (Bayleef)...
    Syncing 45 moves for pokemon 153 (Bayleef)...
    Inserting 1 abilities for pokemon 153 (Bayleef)...
    Inserting pokedex entry for species 153...
  [2/50] Processing Pokemon 156...
  ...
✓ Completed syncing 50 Pokemon for special game
✓ Completed colosseum
```

## Querying Special Games

After implementing virtual pokedexes, **all your existing queries work unchanged**!

The query you mentioned:
```sql
SELECT
  p.id,
  p.name
FROM versions v
JOIN version_groups vg ON v.version_group_id = vg.id
JOIN version_group_pokedexes vgp ON vg.id = vgp.version_group_id
JOIN pokedex_entries pe ON vgp.pokedex_id = pe.pokedex_id
JOIN species s ON pe.species_id = s.id
JOIN pokemon p ON s.id = p.species_id
WHERE v.name = "colosseum"
  AND p.is_default = TRUE
ORDER BY pe.entry_number;
```

This will now work for Colosseum because:
- ✅ `version_group_pokedexes` has an entry linking to the virtual pokedex
- ✅ `pokedex_entries` has entries for all 50 Colosseum Pokemon
- ✅ Entry numbers are sequential (1-50)

## Files Modified

1. **`models/special_games.go`** (NEW)
   - Contains Pokemon lists for special games
   - `GetSpecialGamePokemon()` function for lookup

2. **`services/game_syncer.go`**
   - Added detection logic in `SyncGame()`
   - Added `syncSpecialGamePokemon()` helper
   - Added `syncPokemonData()` helper

3. **`db/pokemon_repository.go`**
   - Removed duplicate `colosseumIds` array (moved to models)

## Notes

- Special games get version and version_group entries in the database (same as regular games)
- **Virtual pokedexes** are created with IDs starting at 1000 to avoid conflicts with real pokedexes (which have IDs < 100)
- Pokedex entries ARE created with sequential entry numbers (1, 2, 3, ...)
- All Pokemon data (species, pokemon, types, moves, abilities) is fully synced
- The sync process respects all caching mechanisms for performance
- All existing queries work unchanged - they don't need to know about "virtual" vs "real" pokedexes
