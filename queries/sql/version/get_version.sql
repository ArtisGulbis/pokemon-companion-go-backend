SELECT
    id,
    name,
    display_name,
    version_group_id
FROM versions
WHERE id = ?
