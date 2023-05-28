package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/qerdcv/ttto/internal/domain"
	"github.com/qerdcv/ttto/internal/ttto/service"
)

func (s *server) setupGameHandlers(g *gin.RouterGroup) {
	gameG := g.Group("/games")
	{
		gameG.POST("", s.authRequired(), s.createGame)
		gameG.GET("/:gameID", s.getGame)
		gameG.PATCH("/:gameID", s.authRequired(), s.makeStep)
		gameG.PATCH("/:gameID/login", s.authRequired(), s.loginGame)
		gameG.GET("/:gameID/history", s.getGameHistory)
		gameG.GET("/:gameID/subscribe", s.subscribeToGameUpdates)
	}
}

// createGame godoc
//
//	@Summary		create new game
//	@Description	create new game
//	@Tags			game
//	@Produce		json
//	@Success		201	{object}	http.Response
//	@Failure		500	{object}	http.Response
//	@Router			/api/games [post]
func (s *server) createGame(c *gin.Context) {
	gID, err := s.service.CreateGame(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": gID,
	})
}

// getGame godoc
//
//	@Summary		get game by id
//	@Description	get game by id
//	@Tags			game
//	@Produce		json
//	@Param			gameID	path		int	true	"id of the game"
//	@Success		200		{object}	domain.Game
//	@Failure		400		{object}	http.Response
//	@Failure		404		{object}	http.Response
//	@Failure		500		{object}	http.Response
//	@Router			/api/games/{gameID} [get]
func (s *server) getGame(c *gin.Context) {
	rawGID := c.Param("gameID")
	gID, _ := strconv.Atoi(rawGID)
	g, err := s.service.GetGame(c.Request.Context(), int32(gID))
	if err != nil {
		if errors.Is(err, service.ErrGameNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": http.StatusText(http.StatusNotFound),
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	c.JSON(http.StatusOK, g)
}

// loginGame godoc
//
//	@Summary		login into the game
//	@Description	login into the game
//	@Tags			game
//	@Produce		json
//	@Param			gameID	path		int	true	"id of the game"
//	@Success		200		{object}	http.Response
//	@Failure		401		{object}	http.Response	"Need to be authorized"
//	@Failure		409		{object}	http.Response	"User already in game"
//	@Failure		400		{object}	http.Response	"Invalid game state"
//	@Failure		404		{object}	http.Response	"Game with provided id is not found"
//	@Failure		500		{object}	http.Response	"Something went wrong"
//	@Router			/api/games/{gameID}/login [patch]
func (s *server) loginGame(c *gin.Context) {
	rawGID := c.Param("gameID")
	gID, _ := strconv.Atoi(rawGID)
	if err := s.service.LoginGame(c.Request.Context(), int32(gID)); err != nil {
		switch {
		case errors.Is(err, service.ErrUnauthorized):
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": http.StatusText(http.StatusUnauthorized),
			})
		case errors.Is(err, service.ErrUserAlreadyInGame):
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"message": http.StatusText(http.StatusConflict),
			})
		case errors.Is(err, service.ErrInvalidGameState):
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": http.StatusText(http.StatusBadRequest),
			})
		case errors.Is(err, service.ErrGameNotFound):
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": http.StatusText(http.StatusNotFound),
			})
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": http.StatusText(http.StatusInternalServerError),
			})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": http.StatusText(http.StatusOK),
	})
}

// makeStep godoc
//
//	@Summary		place mark into the field
//	@Description	place mark into the field
//	@Tags			game
//	@Accept			json
//	@Produce		json
//	@Param			gameID	path		int			true	"id of the game"
//	@Param			step	body		domain.Step	true	"Step coords"
//	@Success		200		{object}	http.Response
//	@Failure		401		{object}	http.Response	"Need to be authorized"
//	@Failure		409		{object}	http.Response	"Cell already occupied"
//	@Failure		400		{object}	http.Response	"Invalid game state"
//	@Failure		404		{object}	http.Response	"Game with provided id is not found"
//	@Failure		500		{object}	http.Response	"Something went wrong"
//	@Router			/api/games/{gameID} [patch]
func (s *server) makeStep(c *gin.Context) {
	var step domain.Step
	if err := c.ShouldBindJSON(&step); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": http.StatusText(http.StatusBadRequest),
		})
		return
	}

	rawGID := c.Param("gameID")
	gID, _ := strconv.Atoi(rawGID)
	if err := s.service.MakeStep(c.Request.Context(), int32(gID), &step); err != nil {
		var valErr *service.ErrValidation
		switch {
		case errors.As(err, &valErr):
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"message": valErr,
			})
		case errors.Is(err, service.ErrUnauthorized):
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": http.StatusText(http.StatusUnauthorized),
			})
		case errors.Is(err, service.ErrUserAlreadyInGame):
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"message": http.StatusText(http.StatusConflict),
			})
		case errors.Is(err, service.ErrInvalidGameState), errors.Is(err, service.ErrNotUsersTurn):
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": http.StatusText(http.StatusBadRequest),
			})
		case errors.Is(err, service.ErrGameNotFound):
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": http.StatusText(http.StatusNotFound),
			})
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": http.StatusText(http.StatusInternalServerError),
			})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": http.StatusText(http.StatusOK),
	})
}

// getGameHistory godoc
//
//	@Summary		get game history
//	@Description	get game history
//	@Tags			game
//	@Produce		json
//	@Param			gameID	path		int	true	"id of the game"
//	@Success		200		{object}	http.Response
//	@Failure		404		{object}	http.Response	"Game with provided id is not found"
//	@Failure		500		{object}	http.Response	"Something went wrong"
//	@Router			/api/games/{gameID}/history [get]
func (s *server) getGameHistory(c *gin.Context) {
	rawGID := c.Param("gameID")
	gID, _ := strconv.Atoi(rawGID)
	gHist, err := s.service.GetGameHistory(c.Request.Context(), int32(gID))
	if err != nil {
		if errors.Is(err, service.ErrGameNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": http.StatusText(http.StatusNotFound),
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	c.JSON(http.StatusOK, gHist)
}

type Client struct {
	name   string
	events chan *DashBoard
}
type DashBoard struct {
	User uint
}

// subscribeToGameUpdates godoc
//
//	@Summary		subscribe to the game updates
//	@Description	subscribe to the game updates
//	@Tags			game
//	@Produce		json
//	@Param			gameID	path		int	true	"id of the game"
//	@Success		200		{object}	http.Response
//	@Failure		404		{object}	http.Response	"Game with provided id is not found"
//	@Failure		500		{object}	http.Response	"Something went wrong"
//	@Router			/api/games/{gameID}/subscribe [get]
func (s *server) subscribeToGameUpdates(c *gin.Context) {

	rawGID := c.Param("gameID")
	gID, _ := strconv.Atoi(rawGID)
	g, err := s.service.GetGame(c.Request.Context(), int32(gID))
	if err != nil {
		if errors.Is(err, service.ErrGameNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": http.StatusText(http.StatusNotFound),
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	w := c.Writer

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	sendGameUpdate := func(g *domain.Game) {
		b, innErr := json.Marshal(g)
		if innErr != nil {
			log.Println("ERROR: ", innErr)
			return
		}

		if _, innErr = fmt.Fprintf(w, "data: %s\n\n", string(b)); innErr != nil {
			log.Println("ERROR: ", innErr)
			return
		}

		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}

	ch := make(chan *domain.Game)
	subscriberID := uuid.NewString()

	defer s.es.Unsubscribe(int32(gID), subscriberID)

	s.es.Subscribe(int32(gID), subscriberID, ch)

	sendGameUpdate(g)
	for {
		select {
		case ev := <-ch:
			sendGameUpdate(ev)
		case <-c.Request.Context().Done():
			return
		}

	}
}
