package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/qerdcv/ttto/internal/domain"
	"github.com/qerdcv/ttto/internal/ttto/repository/storage"
)

func (r *Repository) CreateGame(ctx context.Context, owner *domain.User) (int32, error) {
	gID, err := r.q.CreateGame(ctx, owner.ID)
	if err != nil {
		return 0, fmt.Errorf("queries create game: %w", err)
	}

	return gID, nil
}

func (r *Repository) GetGame(ctx context.Context, gID int32) (*domain.Game, error) {
	dbG, err := r.q.GetGame(ctx, gID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("queries get game: %w", err)
	}

	game, err := domainGameFromDB(dbG)
	if err != nil {
		return nil, fmt.Errorf("domain game from db: %w", err)
	}

	return game, nil
}

func (r *Repository) UpdateGame(ctx context.Context, g *domain.Game) error {
	field, err := json.Marshal(g.Field)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	params := storage.UpdateGameParams{
		ID:           g.ID,
		OwnerID:      g.Owner.ID,
		StepCount:    g.StepCount,
		Field:        field,
		CurrentState: g.CurrentState.String(),
	}

	if g.Opponent != nil {
		params.OpponentID = sql.NullInt32{
			Int32: g.Opponent.ID,
			Valid: true,
		}
	}

	if g.CurrentPlayer != nil {
		params.CurrentPlayerID = sql.NullInt32{
			Int32: g.CurrentPlayer.ID,
			Valid: true,
		}
	}

	if g.Winner != nil {
		params.WinnerID = sql.NullInt32{
			Int32: g.Winner.ID,
			Valid: true,
		}
	}

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("db begin: %w", err)
	}

	defer tx.Rollback()

	q := r.q.WithTx(tx)
	if err = q.UpdateGame(ctx, params); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}

		return fmt.Errorf("queries update game: %w", err)
	}

	if err = q.CreateGameHistoryRecord(ctx, storage.CreateGameHistoryRecordParams(params)); err != nil {
		return fmt.Errorf("create game history record ")
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("tx commit: %w", err)
	}

	return nil
}

func domainGameFromDB(dbG storage.GamesWithUsername) (*domain.Game, error) {
	game := &domain.Game{
		ID: dbG.ID,
		Owner: &domain.User{
			ID:       dbG.OwnerID,
			Username: dbG.OwnerName.String,
		},
		StepCount:    dbG.StepCount,
		CurrentState: domain.State(dbG.CurrentState),
	}

	if err := json.Unmarshal(dbG.Field, &game.Field); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	if dbG.OpponentName.Valid {
		game.Opponent = &domain.User{
			ID:       dbG.OpponentID.Int32,
			Username: dbG.OpponentName.String,
		}
	}

	if dbG.CurrentPlayerName.Valid {
		game.CurrentPlayer = &domain.User{
			ID:       dbG.CurrentPlayerID.Int32,
			Username: dbG.CurrentPlayerName.String,
		}
	}

	if dbG.WinnerName.Valid {
		game.Opponent = &domain.User{
			ID:       dbG.WinnerID.Int32,
			Username: dbG.WinnerName.String,
		}
	}

	return game, nil
}

func (r *Repository) GetGameHistory(ctx context.Context, gID int32) ([]*domain.Game, error) {
	dbGHistory, err := r.q.GetGameHistory(ctx, gID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("queries get game history: %w", err)
	}

	gHistory := make([]*domain.Game, len(dbGHistory))
	for i, dbGHist := range dbGHistory {
		var field domain.Field

		if unmarshalErr := json.Unmarshal(dbGHist.Field, &field); err != nil {
			return nil, fmt.Errorf("json unmarshal: %w", unmarshalErr)
		}

		gHistory[i] = &domain.Game{
			ID: dbGHist.ID,
			Owner: &domain.User{
				ID:       dbGHist.OwnerID,
				Username: dbGHist.OwnerName,
			},
			StepCount:    dbGHist.StepCount,
			Field:        field,
			CurrentState: domain.State(dbGHist.CurrentState),
			Opponent: &domain.User{
				ID:       dbGHist.OpponentID.Int32,
				Username: dbGHist.OpponentName,
			},
			CurrentPlayer: &domain.User{
				ID:       dbGHist.CurrentPlayerID.Int32,
				Username: dbGHist.CurrentPlayerName,
			},
		}

		if dbGHist.WinnerID.Valid {
			gHistory[i].Winner = &domain.User{
				ID:       dbGHist.WinnerID.Int32,
				Username: dbGHist.WinnerName.String,
			}
		}
	}

	return gHistory, nil
}
