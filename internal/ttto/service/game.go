package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"

	"github.com/qerdcv/ttto/internal/domain"
	"github.com/qerdcv/ttto/internal/ttto/repository"
	"github.com/qerdcv/ttto/internal/xctx"
)

var (
	ErrUnauthorized      = errors.New("unauthorized")
	ErrGameNotFound      = errors.New("game not found")
	ErrUserAlreadyInGame = errors.New("user already in game")
	ErrInvalidGameState  = errors.New("invalid game state")
	ErrCellOccupied      = errors.New("cell occupied")
	ErrNotUsersTurn      = errors.New("not user's turn")
)

func (s *Service) CreateGame(ctx context.Context) (int32, error) {
	u := xctx.UserFromContext(ctx)
	if u == nil {
		return 0, ErrUnauthorized
	}

	gID, err := s.repo.CreateGame(ctx, u)
	if err != nil {
		return 0, fmt.Errorf("repo create game: %w", err)
	}

	return gID, nil
}

func (s *Service) GetGame(ctx context.Context, gID int32) (*domain.Game, error) {
	g, err := s.repo.GetGame(ctx, gID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrGameNotFound
		}

		return nil, fmt.Errorf("repo get game: %w", err)
	}

	return g, nil
}

func (s *Service) UpdateGame(ctx context.Context, g *domain.Game) (*domain.Game, error) {
	if err := s.repo.UpdateGame(ctx, g); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrGameNotFound
		}

		return nil, fmt.Errorf("repo update game: %w", err)
	}

	s.es.SendEvent(g.ID, g)

	return nil, nil
}

func (s *Service) LoginGame(ctx context.Context, gID int32) error {
	u := xctx.UserFromContext(ctx)
	if u == nil {
		return ErrUnauthorized
	}

	g, err := s.GetGame(ctx, gID)
	if err != nil {
		return fmt.Errorf("get game: %w", err)
	}

	if g.CurrentState != domain.PendingState {
		return ErrInvalidGameState
	}

	if u.ID == g.Owner.ID {
		return ErrUserAlreadyInGame
	}

	g.Opponent = u
	g.CurrentPlayer = []*domain.User{g.Owner, g.Opponent}[rand.Intn(2)]
	g.CurrentState = domain.InGameState

	if _, err = s.UpdateGame(ctx, g); err != nil {
		return err
	}

	return nil
}

func (s *Service) MakeStep(ctx context.Context, gID int32, step *domain.Step) error {
	u := xctx.UserFromContext(ctx)
	if u == nil {
		return ErrUnauthorized
	}

	if err := step.Validate(); err != nil {
		return newErrValidation(err)
	}

	g, err := s.GetGame(ctx, gID)
	if err != nil {
		return err
	}

	if g.CurrentState != domain.InGameState {
		return ErrInvalidGameState
	}

	if g.CurrentPlayer.ID != u.ID {
		return ErrNotUsersTurn
	}

	if g.Field[step.Row][step.Col] != "" {
		return ErrCellOccupied
	}

	fieldSize := len(g.Field) * len(g.Field)
	g.StepCount += 1
	g.Field[step.Row][step.Col] = g.CurrentPlayerMark()
	if g.StepCount >= int32(math.Ceil(float64(fieldSize)/2)) && g.IsWin() {
		g.Winner = g.CurrentPlayer
		g.CurrentState = domain.DoneState
	} else {
		if g.CurrentPlayer.ID == g.Owner.ID {
			g.CurrentPlayer = g.Opponent
		} else {
			g.CurrentPlayer = g.Owner
		}
	}

	if g.StepCount == int32(fieldSize) {
		g.CurrentState = domain.DoneState
	}

	if _, err = s.UpdateGame(ctx, g); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetGameHistory(ctx context.Context, gID int32) ([]*domain.Game, error) {
	gHist, err := s.repo.GetGameHistory(ctx, gID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrGameNotFound
		}

		return nil, fmt.Errorf("repo get game history: %w", err)
	}

	return gHist, nil
}
