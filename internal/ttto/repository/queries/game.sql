-- name: CreateGame :one
INSERT INTO games (owner_id)
VALUES ($1)
RETURNING id;

-- name: GetGame :one
SELECT *
FROM games_with_usernames
WHERE id = $1;

-- name: UpdateGame :exec
UPDATE games
SET owner_id=$2,
    opponent_id=$3,
    current_player_id=$4,
    step_count=$5,
    winner_id=$6,
    field=$7,
    current_state=$8
WHERE id = $1;

-- name: CreateGameHistoryRecord :exec
INSERT INTO games_history(
    game_id,
    owner_id,
    opponent_id,
    current_player_id,
    step_count,
    winner_id,
    field,
    current_state
)
VALUES(
    @id,
    @owner_id,
    @opponent_id,
    @current_player_id,
    @step_count,
    @winner_id,
    @field,
    @current_state
);

-- name: GetGameHistory :many
SELECT
    gh.game_id as id,
    gh.owner_id, gh.opponent_id, gh.current_player_id,
    gh.step_count, gh.winner_id, gh.field, gh.current_state,
    owner.username as owner_name,
    opponent.username as opponent_name,
    current_player.username as current_player_name,
    winner.username as winner_name
FROM games_history gh
         JOIN users owner on owner.id = gh.owner_id
         JOIN users opponent on opponent.id = gh.opponent_id
         JOIN users current_player on current_player.id = gh.current_player_id
         LEFT JOIN users winner on winner.id = gh.winner_id
WHERE gh.game_id=$1
ORDER BY gh.id;