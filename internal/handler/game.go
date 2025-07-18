package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"scrable3/internal/dto"
	"scrable3/internal/model"
	"scrable3/internal/svc"
	"time"

	"github.com/google/uuid"
)

type gameHandler struct {
	gameService   svc.GameService
	playerService svc.PlayerService
	fieldService  svc.FieldService
}

func NewGameHandler(gameService svc.GameService, playerService svc.PlayerService, fieldService svc.FieldService) http.Handler {
	return &gameHandler{
		gameService:   gameService,
		playerService: playerService,
		fieldService:  fieldService,
	}
}

// |PRIVATE| //

func (h *gameHandler) setCookieWithPlayerUUID(
	w http.ResponseWriter, game *model.Game, player *model.Player,
) {
	http.SetCookie(w, &http.Cookie{
		Name:     "player-uuid-" + game.UUID.String(),
		Value:    player.UUID.String(),
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})
}

func (h *gameHandler) getGame(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/game/game.html")
	if err != nil {
		fmt.Printf(": %v", err)
		http.Error(w, err.Error(), 500)
		return
	}
	gameUUID, err := uuid.Parse(r.PathValue("gameUUID"))
	if err != nil {
		fmt.Printf(": %v", err)
		http.Error(w, err.Error(), 500)
		return
	}
	game, err := h.gameService.GetWithUUID(gameUUID)
	if err != nil {
		fmt.Printf("game: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var player *model.Player
	playerUUIDCookie, err := r.Cookie("player-uuid-" + gameUUID.String())
	if err != nil {
		player, err = h.playerService.Create(game)
		if err != nil {
			fmt.Printf("player: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		h.setCookieWithPlayerUUID(w, game, player)
	} else {
		playerUUID, err := uuid.Parse(playerUUIDCookie.Value)
		if err != nil {
			fmt.Printf(": %v", err)
			http.Error(w, err.Error(), 500)
			return
		}

		player, err = h.playerService.GetWithUUID(playerUUID)
		if err != nil {
			fmt.Printf("player: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	if player.GameUUID != game.UUID {
		http.Error(w, "player is not linked to this game", http.StatusUnauthorized)
		return
	}

	data := dto.GamePageData{Title: "Game", GameUUID: game.UUID.String()}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (h *gameHandler) createGameInitialData(gameUUID uuid.UUID, playerUUID uuid.UUID) error {
	_, err := h.fieldService.Create(gameUUID, playerUUID, 0, "A", [3]int{7, 7, 7})
	return err
}

func (h *gameHandler) createGame(w http.ResponseWriter) {
	game, err := h.gameService.Create()
	if err != nil {
		fmt.Printf("h.GameService.Create(): %v", err)
		http.Error(w, err.Error(), 500)
		return
	}
	player, err := h.playerService.Create(game)
	if err != nil {
		fmt.Printf("h.PlayerService.Create(game): %v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	err = h.createGameInitialData(game.UUID, player.UUID)
	if err != nil {
		fmt.Printf("h.createGameInitialData: %v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	h.setCookieWithPlayerUUID(w, game, player)
	redirectURL := "/game/" + game.UUID.String()
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusOK)
}

// |PUBLIC| //

func (h *gameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.PathValue("gameUUID") != "" {
			h.getGame(w, r)
		}
	case http.MethodPost:
		h.createGame(w)
	}
}
