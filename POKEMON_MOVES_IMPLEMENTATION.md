# Pokemon Moves Implementation Guide

## Overview

This document explains how `pokemon_moves` insertion works in the syncing system.

## The Problem You Were Facing

You were confused about **where to put the InsertPokemonMove method** and **how to wire it up** because:

1. The method was in the wrong file (`pokemon_repository.go` instead of `move_repository.go`)
2. The method signature was incomplete (missing required parameters)
3. The receiver type was wrong (`*MoveRepository` in a `PokemonRepository` file)
4. The sync flow wasn't calling it properly

## Understanding the Data Model

### Tables Involved

```sql
-- The move itself (Thunderbolt, Tackle, etc.)
CREATE TABLE moves (
    id INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    ...
);

-- Which Pokemon learns which move in which game
CREATE TABLE pokemon_moves (
    pokemon_id INTEGER NOT NULL,          -- FK to pokemon(id)
    move_id INTEGER NOT NULL,             -- FK to moves(id)
    version_group_id INTEGER NOT NULL,    -- FK to version_groups(id) - IMPORTANT!
    learn_method TEXT NOT NULL,           -- "level-up", "machine", "egg", "tutor"
    level_learned_at INTEGER NOT NULL,    -- e.g., 15 (or 0 for non-level-up)
    PRIMARY KEY (pokemon_id, move_id, version_group_id, learn_method)
);
```

### Why version_group_id?

Pokemon learn **different moves in different games**! For example:
- Pikachu learns Thunder at level 41 in Red/Blue
- Pikachu learns Thunder at level 50 in X/Y
- Some moves aren't available in certain games at all

So we need to track **which game** (version_group) each pokemon_move belongs to.

## The Data Flow (from PokeAPI)

When you fetch a Pokemon from PokeAPI, you get:

```json
{
  "id": 25,
  "name": "pikachu",
  "moves": [
    {
      "move": { "name": "thunderbolt", "url": ".../move/85/" },
      "version_group_details": [
        {
          "level_learned_at": 26,
          "move_learn_method": { "name": "level-up" },
          "version_group": { "name": "red-blue", "url": ".../version-group/1/" }
        },
        {
          "level_learned_at": 29,
          "move_learn_method": { "name": "level-up" },
          "version_group": { "name": "gold-silver", "url": ".../version-group/3/" }
        }
      ]
    }
  ]
}
```

Notice:
- **One move** (Thunderbolt)
- **Multiple version_group_details** (one for each game it appears in)
- Each detail has: level, learn method, and version_group

## The Solution

### 1. Repository Layer (`db/move_repository.go`)

```go
func (r *MoveRepository) InsertPokemonMove(
    pokemonID, moveID, versionGroupID int,
    learnMethod string,
    levelLearnedAt int,
) error {
    // Uses the query from queries/sql/move/pokemon_moves.sql
    stmt, err := r.db.Prepare(queries.InsertPokemonMoves)
    // ... exec with all 5 parameters
}
```

**Why MoveRepository?** Because it's move-related data.

### 2. Interface Layer (`services/interfaces.go`)

```go
type MoveRepo interface {
    InsertMove(v *external.Move) error
    InsertPokemonMove(pokemonID, moveID, versionGroupID int, learnMethod string, levelLearnedAt int) error
    GetMoveByID(id int) (*dto.Move, error)
}
```

### 3. Service Layer (`services/game_syncer.go`)

The `syncPokemonData` function handles everything:

```go
func (g *GameSyncer) syncPokemonData(pokemonID int, versionGroupID int) error {
    // 1. Sync the Pokemon itself
    pokemon, err := g.pokemonSyncer.SyncPokemon(pokemonID)

    // 2. Insert types
    // ...

    // 3. Sync moves
    for _, m := range pokemon.Moves {
        moveId := extractID(m.Move.Url)

        // 3a. Sync the move itself (INSERT INTO moves)
        g.moveSyncer.SyncMove(moveId)

        // 3b. Insert pokemon_moves for THIS version_group only
        for _, vgDetail := range m.VersionGroupDetails {
            vgID := extractID(vgDetail.VersionGroup.Url)

            // Only insert if this move is learned in the CURRENT game
            if vgID == versionGroupID {
                g.moveSyncer.repo.InsertPokemonMove(
                    pokemon.ID,
                    moveId,
                    versionGroupID,
                    vgDetail.MoveLearnMethod.Name,
                    vgDetail.LevelLearnedAt,
                )
            }
        }
    }

    // 4. Insert abilities
    // ...
}
```

### Key Points:

1. **We filter by versionGroupID** - only insert moves for the current game being synced
2. **Two inserts per move**:
   - `InsertMove()` - inserts into `moves` table (the move itself)
   - `InsertPokemonMove()` - inserts into `pokemon_moves` table (the relationship)
3. **Caching still works** - moves are cached, so duplicate API calls are avoided

## How It's Called

### Normal Games (with Pokedexes)

```go
// In SyncGame, for each Pokemon in the pokedex:
g.pokemonSyncer.SyncSpecies(speciesID)
g.syncPokemonData(speciesID, versionGroupId)  // ← passes version_group_id
```

### Special Games (Colosseum, XD)

```go
// In syncSpecialGamePokemon, for each hardcoded Pokemon:
g.pokemonSyncer.SyncSpecies(pokemonID)
g.syncPokemonData(pokemonID, versionGroupID)  // ← passes version_group_id
```

## Query Protection

The SQL query uses `INSERT OR IGNORE`:

```sql
INSERT OR IGNORE INTO pokemon_moves (pokemon_id, move_id, version_group_id, learn_method, level_learned_at)
VALUES (?, ?, ?, ?, ?)
```

This prevents errors if:
- The same Pokemon appears in multiple pokedexes in the same game
- You re-sync the same game multiple times

## Files Modified

1. **`db/pokemon_repository.go`**
   - Removed broken `InsertPokemonMove` method (was on wrong receiver)

2. **`db/move_repository.go`**
   - Added correct `InsertPokemonMove` method

3. **`services/interfaces.go`**
   - Added `InsertPokemonMove` to `MoveRepo` interface
   - Removed incorrect method from `PokemonRepo` interface

4. **`services/pokemon_syncer.go`**
   - Removed broken `SyncPokemonMove` method
   - Removed unused `syncedMoves` cache

5. **`services/game_syncer.go`**
   - Updated `syncPokemonData` to accept `versionGroupID` parameter
   - Added pokemon_moves insertion logic with version_group filtering
   - Refactored normal flow to use `syncPokemonData` (removed duplicate code)
   - Fixed broken pokemon_moves insertion attempt in normal flow

6. **`queries/sql/move/pokemon_moves.sql`**
   - Added `INSERT OR IGNORE` to prevent duplicate errors

## Testing

After syncing a game, you can verify pokemon_moves were inserted:

```sql
-- Check how many moves Pikachu has in Pokemon Red
SELECT
    p.name as pokemon,
    m.name as move,
    pm.learn_method,
    pm.level_learned_at
FROM pokemon_moves pm
JOIN pokemon p ON pm.pokemon_id = p.id
JOIN moves m ON pm.move_id = m.id
JOIN version_groups vg ON pm.version_group_id = vg.id
WHERE p.name = 'pikachu'
  AND vg.name = 'red-blue'
ORDER BY pm.level_learned_at;
```

## Why This Approach?

✅ **Separation of Concerns**: Move-related data goes in MoveRepository
✅ **Version-specific**: Only stores moves relevant to each game
✅ **No Duplicates**: `INSERT OR IGNORE` prevents constraint errors
✅ **DRY**: Both normal and special games use the same `syncPokemonData` function
✅ **Cached**: Move syncing is cached, reducing API calls
✅ **Complete**: Stores all 5 required fields for each pokemon_move
