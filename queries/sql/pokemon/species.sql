INSERT INTO species (
    id, 
    name,
    evolution_chain_id,
    gender_rate,       
    capture_rate,
    base_happiness,
    is_baby,
    is_legendary,
    is_mythical,
    growth_rate_name,   
    generation_name
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);