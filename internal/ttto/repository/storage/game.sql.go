// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: game.sql

package storage

import (
	"context"
	"database/sql"
	"encoding/json"
)

const createGame = `-- name: CreateGame :one
INSERT INTO games (owner_id)
VALUES ($1)
RETURNING id
`

func (q *Queries) CreateGame(ctx context.Context, ownerID int32) (int32, error) {
	row := q.queryRow(ctx, q.createGameStmt, createGame, ownerID)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const createGameHistoryRecord = `-- name: CreateGameHistoryRecord :exec
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
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
`

type CreateGameHistoryRecordParams struct {
	ID              int32
	OwnerID         int32
	OpponentID      sql.NullInt32
	CurrentPlayerID sql.NullInt32
	StepCount       int32
	WinnerID        sql.NullInt32
	Field           json.RawMessage
	CurrentState    string
}

func (q *Queries) CreateGameHistoryRecord(ctx context.Context, arg CreateGameHistoryRecordParams) error {
	_, err := q.exec(ctx, q.createGameHistoryRecordStmt, createGameHistoryRecord,
		arg.ID,
		arg.OwnerID,
		arg.OpponentID,
		arg.CurrentPlayerID,
		arg.StepCount,
		arg.WinnerID,
		arg.Field,
		arg.CurrentState,
	)
	return err
}

const getGame = `-- name: GetGame :one
SELECT id, owner_id, opponent_id, current_player_id, step_count, winner_id, field, current_state, owner_name, opponent_name, current_player_name, winner_name
FROM games_with_usernames
WHERE id = $1
`

func (q *Queries) GetGame(ctx context.Context, id int32) (GamesWithUsername, error) {
	row := q.queryRow(ctx, q.getGameStmt, getGame, id)
	var i GamesWithUsername
	err := row.Scan(
		&i.ID,
		&i.OwnerID,
		&i.OpponentID,
		&i.CurrentPlayerID,
		&i.StepCount,
		&i.WinnerID,
		&i.Field,
		&i.CurrentState,
		&i.OwnerName,
		&i.OpponentName,
		&i.CurrentPlayerName,
		&i.WinnerName,
	)
	return i, err
}

const getGameHistory = `-- name: GetGameHistory :many
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
ORDER BY gh.id
`

type GetGameHistoryRow struct {
	ID                int32
	OwnerID           int32
	OpponentID        sql.NullInt32
	CurrentPlayerID   sql.NullInt32
	StepCount         int32
	WinnerID          sql.NullInt32
	Field             json.RawMessage
	CurrentState      string
	OwnerName         string
	OpponentName      string
	CurrentPlayerName string
	WinnerName        sql.NullString
}

func (q *Queries) GetGameHistory(ctx context.Context, gameID int32) ([]GetGameHistoryRow, error) {
	rows, err := q.query(ctx, q.getGameHistoryStmt, getGameHistory, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetGameHistoryRow
	for rows.Next() {
		var i GetGameHistoryRow
		if err := rows.Scan(
			&i.ID,
			&i.OwnerID,
			&i.OpponentID,
			&i.CurrentPlayerID,
			&i.StepCount,
			&i.WinnerID,
			&i.Field,
			&i.CurrentState,
			&i.OwnerName,
			&i.OpponentName,
			&i.CurrentPlayerName,
			&i.WinnerName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateGame = `-- name: UpdateGame :exec
UPDATE games
SET owner_id=$2,
    opponent_id=$3,
    current_player_id=$4,
    step_count=$5,
    winner_id=$6,
    field=$7,
    current_state=$8
WHERE id = $1
`

type UpdateGameParams struct {
	ID              int32
	OwnerID         int32
	OpponentID      sql.NullInt32
	CurrentPlayerID sql.NullInt32
	StepCount       int32
	WinnerID        sql.NullInt32
	Field           json.RawMessage
	CurrentState    string
}

func (q *Queries) UpdateGame(ctx context.Context, arg UpdateGameParams) error {
	_, err := q.exec(ctx, q.updateGameStmt, updateGame,
		arg.ID,
		arg.OwnerID,
		arg.OpponentID,
		arg.CurrentPlayerID,
		arg.StepCount,
		arg.WinnerID,
		arg.Field,
		arg.CurrentState,
	)
	return err
}
