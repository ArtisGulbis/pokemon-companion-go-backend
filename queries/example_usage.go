package queries

// This file shows examples of how to use the embedded queries
// You can delete this file - it's just for reference

import (
	"database/sql"
	"fmt"
)

// ExampleSavePokemon shows how to use the InsertPokemon query
func ExampleSavePokemon(db *sql.DB, id int, name string, height, weight float64) error {
	// Step 1: Prepare the statement
	// Preparing compiles the SQL and makes it ready for execution
	// Using Prepare is good practice because:
	// - It prevents SQL injection attacks
	// - It's more efficient if you run the query multiple times
	stmt, err := db.Prepare(InsertPokemon)
	if err != nil {
		// %w wraps the error so you can use errors.Is() and errors.As() later
		return fmt.Errorf("failed to prepare pokemon insert: %w", err)
	}
	// defer means "run this when the function exits"
	// We always want to close prepared statements to free resources
	defer stmt.Close()

	// Step 2: Execute the query with parameters
	// The ? placeholders in the SQL are replaced with these values in order
	result, err := stmt.Exec(id, name, height, weight)
	if err != nil {
		return fmt.Errorf("failed to execute pokemon insert: %w", err)
	}

	// Optional: Check how many rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	fmt.Printf("Inserted/updated Pokemon, rows affected: %d\n", rowsAffected)
	return nil
}

// ExampleGetPokemon shows how to use the GetPokemonByID query
func ExampleGetPokemon(db *sql.DB, id int) error {
	// For queries that return a single row, use QueryRow
	// QueryRow returns a *sql.Row that you can Scan into variables
	row := db.QueryRow(GetPokemonByID, id)

	// Declare variables to hold the results
	var pokemonID int
	var name string
	var height, weight float64

	// Scan reads the row data into your variables
	// The order must match the SELECT column order
	// The & symbol means "address of" - we're giving Scan pointers
	// so it can modify our variables directly
	err := row.Scan(&pokemonID, &name, &height, &weight)
	if err != nil {
		// sql.ErrNoRows is a special error when no rows are found
		// You might want to handle this differently than other errors
		if err == sql.ErrNoRows {
			return fmt.Errorf("pokemon with id %d not found", id)
		}
		return fmt.Errorf("failed to scan pokemon: %w", err)
	}

	// Use your data
	fmt.Printf("Found Pokemon: ID=%d, Name=%s, Height=%.1f, Weight=%.1f\n",
		pokemonID, name, height, weight)
	return nil
}
