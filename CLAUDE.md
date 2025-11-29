# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Golang backend application that provides a secure REST API for syncing pokemon data. It syncs Pokemon information from the external PokeAPI (https://pokeapi.co/docs/v2) and stores it in MongoDB.

## Common Commands

### Development
```bash
# Start development server with hot reload
go run main.go
```

### Testing
```bash
# Run all unit tests
go test
```

### MongoDB
```bash
# Start MongoDB via Docker Compose
docker-compose up -d

# Stop MongoDB
docker-compose down
```

## Environment Variables

Required environment variables (see `.env.example`):

- `MONGODB_URI`: MongoDB connection string
- `PORT`: Server port (default: 3000)
- `POKE_API_BASE_URL`: PokeAPI base URL for Pokemon
- `POKE_API_ABILITY_BASE_URL`: PokeAPI URL for abilities
- `POKE_API_MOVE_BASE_URL`: PokeAPI URL for moves
- `POKE_API_SPECIES_URL`: PokeAPI URL for species (if needed)
- `POKE_API_POKEDEX_URL`: PokeAPI URL for Pokedex (if needed)