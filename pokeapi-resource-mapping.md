# PokeAPI Resource Mapping: Complete Data Flow Guide

## The Core Confusion: What Each "Game-Related" Resource Actually Means

Before diving into mappings, let's clarify what these similar-sounding things actually represent:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         GAME-RELATED RESOURCES                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  GENERATION (e.g., "generation-v")                                          │
│  ├── What it is: A real-world release ERA (Gen 5 = 2010-2012)              │
│  ├── What it tells you: Which Pokemon were INTRODUCED in that era          │
│  └── NOT useful for: "What Pokemon are in Pokemon Black"                   │
│      (Gen 5 introduced ~150 new Pokemon, but Black has ~300+ available)    │
│                                                                             │
│  VERSION (e.g., "black", "white", "sword")                                  │
│  ├── What it is: A single game cartridge/title                             │
│  ├── What it tells you: Just the name and which version-group it belongs to│
│  └── Use it for: Looking up a game by name, then getting its version-group │
│                                                                             │
│  VERSION-GROUP (e.g., "black-white", "sword-shield")                        │
│  ├── What it is: Games that share IDENTICAL mechanics and Pokemon pools    │
│  ├── What it tells you: Which pokedexes apply, which move-learn-methods    │
│  └── THIS IS THE KEY CONNECTOR - it links to pokedexes and move data       │
│                                                                             │
│  POKEDEX (e.g., "original-unova", "national")                               │
│  ├── What it is: An in-game catalog of Pokemon for a specific region/game  │
│  ├── What it tells you: The actual list of Pokemon available               │
│  └── Use it for: "What Pokemon can I catch/obtain in this game"            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## MASTER FLOW DIAGRAM: From Game Selection to All Data

```
USER SELECTS: "Pokemon Black"
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ STEP 1: Game → Version-Group → Pokedex(es)                                  │
│ Purpose: Find which Pokemon are available in this game                      │
└─────────────────────────────────────────────────────────────────────────────┘
         │
         │  GET /version/black
         │  Response: { "name": "black", "version_group": { "name": "black-white" } }
         │
         ▼
         │  GET /version-group/black-white
         │  Response: {
         │    "pokedexes": [
         │      { "name": "original-unova" }   ◄── Regional dex for this game
         │    ],
         │    "move_learn_methods": [...],     ◄── How Pokemon learn moves in this game
         │    "generation": { "name": "generation-v" }
         │  }
         │
         ▼
         │  GET /pokedex/original-unova
         │  Response: {
         │    "pokemon_entries": [
         │      { "entry_number": 1, "pokemon_species": { "name": "victini" } },
         │      { "entry_number": 2, "pokemon_species": { "name": "snivy" } },
         │      ... (156 total for original-unova)
         │    ]
         │  }
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ RESULT: You now have a list of pokemon-species names available in Black    │
│ These are SPECIES (e.g., "pikachu"), not pokemon IDs yet                    │
└─────────────────────────────────────────────────────────────────────────────┘
         │
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ STEP 2: For each Species → Get Pokemon data (stats, abilities, sprites)    │
│ Purpose: Get the actual gameplay data for each Pokemon                      │
└─────────────────────────────────────────────────────────────────────────────┘
         │
         │  For each species from the pokedex...
         │
         │  GET /pokemon-species/pikachu
         │  Response: {
         │    "id": 25,
         │    "evolution_chain": { "url": "/evolution-chain/10/" },  ◄── For evolutions
         │    "varieties": [
         │      { "is_default": true, "pokemon": { "name": "pikachu" } },
         │      { "is_default": false, "pokemon": { "name": "pikachu-gmax" } },
         │      ... more forms
         │    ],
         │    "flavor_text_entries": [...]  ◄── Pokedex descriptions per game
         │  }
         │
         ▼
         │  GET /pokemon/pikachu  (the default variety)
         │  Response: {
         │    "id": 25,
         │    "stats": [
         │      { "base_stat": 35, "stat": { "name": "hp" } },
         │      { "base_stat": 55, "stat": { "name": "attack" } },
         │      { "base_stat": 40, "stat": { "name": "defense" } },
         │      { "base_stat": 50, "stat": { "name": "special-attack" } },
         │      { "base_stat": 50, "stat": { "name": "special-defense" } },
         │      { "base_stat": 90, "stat": { "name": "speed" } }
         │    ],
         │    "types": [
         │      { "slot": 1, "type": { "name": "electric" } }
         │    ],
         │    "abilities": [
         │      { "ability": { "name": "static" }, "is_hidden": false, "slot": 1 },
         │      { "ability": { "name": "lightning-rod" }, "is_hidden": true, "slot": 3 }
         │    ],
         │    "sprites": {
         │      "front_default": "https://raw.githubusercontent.com/.../25.png",
         │      "front_shiny": "...",
         │      "other": {
         │        "official-artwork": { "front_default": "..." }
         │      }
         │    },
         │    "moves": [...]  ◄── We'll filter this by version-group
         │  }
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ RESULT: You now have stats, types, abilities, sprites for this Pokemon     │
└─────────────────────────────────────────────────────────────────────────────┘
         │
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ STEP 3: Filter moves by version-group                                       │
│ Purpose: Get moves THIS Pokemon can learn IN THIS SPECIFIC GAME             │
└─────────────────────────────────────────────────────────────────────────────┘
         │
         │  The pokemon.moves array contains ALL moves across ALL games.
         │  You must filter by version_group_details.version_group.name
         │
         │  Example from /pokemon/pikachu:
         │  "moves": [
         │    {
         │      "move": { "name": "thunder-shock" },
         │      "version_group_details": [
         │        {
         │          "level_learned_at": 1,
         │          "move_learn_method": { "name": "level-up" },
         │          "version_group": { "name": "black-white" }    ◄── MATCH!
         │        },
         │        {
         │          "level_learned_at": 5,
         │          "move_learn_method": { "name": "level-up" },
         │          "version_group": { "name": "sword-shield" }   ◄── Different game
         │        }
         │      ]
         │    },
         │    {
         │      "move": { "name": "thunderbolt" },
         │      "version_group_details": [
         │        {
         │          "level_learned_at": 0,                        ◄── 0 = not level-up
         │          "move_learn_method": { "name": "machine" },   ◄── TM/HM
         │          "version_group": { "name": "black-white" }
         │        }
         │      ]
         │    }
         │  ]
         │
         │  FILTER LOGIC (pseudocode):
         │  for move in pokemon.moves:
         │      for detail in move.version_group_details:
         │          if detail.version_group.name == "black-white":
         │              save(move.name, detail.level_learned_at, detail.move_learn_method)
         │
         ▼
         │  To get move details (power, accuracy, type):
         │  GET /move/thunder-shock
         │  Response: {
         │    "name": "thunder-shock",
         │    "power": 40,
         │    "accuracy": 100,
         │    "pp": 30,
         │    "type": { "name": "electric" },
         │    "damage_class": { "name": "special" },
         │    "effect_entries": [{ "effect": "...", "short_effect": "..." }]
         │  }
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ RESULT: Moves for Pikachu in Pokemon Black:                                 │
│ - Thunder Shock: Level 1, level-up                                          │
│ - Thunderbolt: TM (level 0)                                                 │
│ - etc.                                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
         │
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ STEP 4: Get Evolution Chain                                                 │
│ Purpose: Show how this Pokemon evolves                                      │
└─────────────────────────────────────────────────────────────────────────────┘
         │
         │  From pokemon-species, we got: evolution_chain.url = "/evolution-chain/10/"
         │
         │  GET /evolution-chain/10
         │  Response: {
         │    "chain": {
         │      "species": { "name": "pichu" },
         │      "evolution_details": [],           ◄── Pichu is the base
         │      "evolves_to": [
         │        {
         │          "species": { "name": "pikachu" },
         │          "evolution_details": [{
         │            "trigger": { "name": "level-up" },
         │            "min_happiness": 220         ◄── Evolves with high friendship
         │          }],
         │          "evolves_to": [
         │            {
         │              "species": { "name": "raichu" },
         │              "evolution_details": [{
         │                "trigger": { "name": "use-item" },
         │                "item": { "name": "thunder-stone" }
         │              }]
         │            }
         │          ]
         │        }
         │      ]
         │    }
         │  }
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ RESULT: Pichu → (happiness) → Pikachu → (Thunder Stone) → Raichu            │
└─────────────────────────────────────────────────────────────────────────────┘
         │
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ STEP 5: Handle Forms/Varieties (Regional variants, Megas, etc.)             │
│ Purpose: Show all the different forms of a Pokemon                          │
└─────────────────────────────────────────────────────────────────────────────┘
         │
         │  From pokemon-species, we got the varieties array.
         │  Each variety is a SEPARATE /pokemon endpoint with different stats.
         │
         │  Example: GET /pokemon-species/meowth
         │  Response: {
         │    "varieties": [
         │      { "is_default": true, "pokemon": { "name": "meowth", "url": "/pokemon/52/" } },
         │      { "is_default": false, "pokemon": { "name": "meowth-alola", "url": "/pokemon/10091/" } },
         │      { "is_default": false, "pokemon": { "name": "meowth-galar", "url": "/pokemon/10161/" } }
         │    ]
         │  }
         │
         │  Each of these has DIFFERENT:
         │  - Stats (Alolan Meowth has different base stats)
         │  - Types (Galarian is Steel-type, not Normal)
         │  - Abilities
         │  - Moves
         │  - Sprites
         │
         │  GET /pokemon/meowth-alola:
         │  - types: [{ "type": { "name": "dark" } }]
         │  - Different stats, different sprites
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ RESULT: Species "Meowth" has 3 varieties you might want to display          │
│ (But check if each variety EXISTS in the selected game!)                    │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## COMPLETE RESOURCE RELATIONSHIP MAP

```
                                    ┌──────────────┐
                                    │   VERSION    │
                                    │   "black"    │
                                    └──────┬───────┘
                                           │
                                           │ version_group
                                           ▼
┌──────────────┐                    ┌──────────────────┐
│  GENERATION  │◄───── generation ──│  VERSION-GROUP   │
│"generation-v"│                    │  "black-white"   │
└──────────────┘                    └────────┬─────────┘
       │                                     │
       │ pokemon_species                     │ pokedexes[]
       │ (introduced in this gen)            ▼
       │                            ┌──────────────────┐
       │                            │    POKEDEX       │
       │                            │"original-unova"  │
       │                            └────────┬─────────┘
       │                                     │
       │                                     │ pokemon_entries[]
       │                                     ▼
       │                            ┌──────────────────┐
       └───────────────────────────►│ POKEMON-SPECIES  │◄──────┐
                                    │   "pikachu"      │       │
                                    └────────┬─────────┘       │
                                             │                 │
                          ┌──────────────────┼──────────────┐  │
                          │                  │              │  │
                          │ varieties[]      │ evolution_   │  │
                          ▼                  │ chain        │  │
                    ┌──────────┐             │              │  │
                    │ POKEMON  │             ▼              │  │
                    │"pikachu" │      ┌─────────────┐       │  │
                    └────┬─────┘      │ EVOLUTION   │       │  │
                         │            │   CHAIN     │───────┘  │
    ┌────────────────────┼────────────┴─────────────┘          │
    │                    │                                     │
    │ stats[]            │ types[]        abilities[]          │
    │ moves[]            │ sprites        forms[]              │
    │                    │                                     │
    │                    │                                     │
    ▼                    ▼                                     │
┌────────┐         ┌──────────┐      ┌───────────────┐         │
│  STAT  │         │   TYPE   │      │    ABILITY    │         │
│ "hp"   │         │"electric"│      │   "static"    │         │
└────────┘         └──────────┘      └───────────────┘         │
                                                               │
    ┌──────────────────────────────────────────────────────────┘
    │
    │  moves[] array in /pokemon contains:
    │
    │  ┌─────────────────────────────────────────────────────────┐
    │  │ {                                                       │
    │  │   "move": { "name": "thunderbolt" },                    │
    │  │   "version_group_details": [                            │
    │  │     {                                                   │
    │  │       "level_learned_at": 26,                           │
    │  │       "move_learn_method": { "name": "level-up" },      │
    │  │       "version_group": { "name": "black-white" } ◄──────┼── Filter by this!
    │  │     }                                                   │
    │  │   ]                                                     │
    │  │ }                                                       │
    │  └─────────────────────────────────────────────────────────┘
    │
    │  Then fetch move details:
    ▼
┌──────────┐
│   MOVE   │
│"thunder- │
│  bolt"   │
├──────────┤
│power: 90 │
│acc: 100  │
│type:elec │
│class:spec│
└──────────┘
```

---

## SUMMARY: Which Endpoints You Need and Why

### Required Endpoints (you MUST fetch these)

| Endpoint | Purpose | When to Fetch |
|----------|---------|---------------|
| `/version/{name}` | Get version-group for a game | Once per game, at startup |
| `/version-group/{name}` | Get pokedexes and move-learn-methods | Once per game, at startup |
| `/pokedex/{name}` | Get list of available pokemon-species | Once per game, at startup |
| `/pokemon-species/{id}` | Get evolution chain, varieties, flavor text | Once per species |
| `/pokemon/{id}` | Get stats, types, abilities, moves, sprites | Once per pokemon variety |
| `/evolution-chain/{id}` | Get evolution tree | Once per chain (shared by family) |
| `/move/{id}` | Get move power, accuracy, type, effect | Once per unique move |

### Optional Endpoints (nice to have)

| Endpoint | Purpose | When to Fetch |
|----------|---------|---------------|
| `/ability/{id}` | Get ability description and effect | When displaying ability details |
| `/type/{id}` | Get type effectiveness chart | Once at startup |
| `/pokemon-form/{id}` | Get cosmetic form sprites | Only for cosmetic variants |

### NOT Needed Endpoints

| Endpoint | Why Not |
|----------|---------|
| `/generation/{id}` | Only tells you what was INTRODUCED, not what's AVAILABLE |
| `/location-area/{id}` | Encounter data - complex and not needed for basic Pokedex |

---

## PRACTICAL EXAMPLE: Building the "Pokemon Black" Dataset

```
STEP 1: Bootstrap game data
─────────────────────────────
GET /version/black
  → version_group = "black-white"

GET /version-group/black-white
  → pokedexes = ["original-unova"]
  → generation = "generation-v"

GET /pokedex/original-unova
  → pokemon_entries = [
      { entry_number: 1, pokemon_species: "victini" },
      { entry_number: 2, pokemon_species: "snivy" },
      ... 156 total
    ]

STEP 2: For each species, fetch Pokemon data
─────────────────────────────────────────────
For species "snivy" (entry #2):

GET /pokemon-species/snivy
  → evolution_chain = "/evolution-chain/125/"
  → varieties = [{ pokemon: "snivy", is_default: true }]
  → flavor_text_entries (filter by version="black")

GET /pokemon/snivy
  → stats: { hp: 45, attack: 45, defense: 55, sp_atk: 45, sp_def: 55, speed: 63 }
  → types: ["grass"]
  → abilities: ["overgrow", "contrary (hidden)"]
  → sprites: { front_default: "...", ... }
  → moves: [... filter where version_group == "black-white" ...]

STEP 3: Get evolution chain
───────────────────────────
GET /evolution-chain/125
  → Snivy → (level 17) → Servine → (level 36) → Serperior

STEP 4: Get move details
────────────────────────
For each move Snivy can learn in Black/White:
GET /move/tackle
GET /move/vine-whip
GET /move/leaf-blade
... etc.
```

---

## KEY GOTCHAS TO REMEMBER

1. **Pokedex ≠ All obtainable Pokemon**
   - Pokedex shows Pokemon with regional dex numbers
   - Some Pokemon are obtainable via transfer but not in the regional dex
   - For most use cases, the regional pokedex is sufficient

2. **Moves MUST be filtered by version-group**
   - The moves array in /pokemon contains moves across ALL games
   - Same move might be learned at different levels in different games
   - Always filter: `version_group_details.version_group.name == "your-game"`

3. **Species vs Pokemon**
   - Species = the concept (Pikachu as a creature)
   - Pokemon = a specific form with stats (regular Pikachu, Gigantamax Pikachu)
   - Start with species, then fetch the `is_default: true` pokemon for basic use

4. **Regional forms are separate Pokemon IDs**
   - Alolan Meowth is `/pokemon/10091`, not `/pokemon/52`
   - They appear in the species' `varieties` array
   - Each has completely different stats, types, and moves

5. **Version-group is the key for game-specific data**
   - NOT version (individual game)
   - NOT generation (introduction era)
   - Use version-group to filter moves, TMs, and match pokedexes
