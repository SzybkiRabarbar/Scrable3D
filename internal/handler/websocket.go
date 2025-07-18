package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"scrable3/internal/ctrl"
	"scrable3/internal/dto"
	"scrable3/internal/svc"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type websocketHandler struct {
	gameService    svc.GameService
	playerService  svc.PlayerService
	gameController ctrl.GameController
	upgrader       websocket.Upgrader
	// hashmap of websocket connections grouped with sessionUUID (gameUUID) and
	// connUUID (playerUUID)
	//  connections[sessionUUID][connUUID] = conn
	connections map[string]map[string]*websocket.Conn
}

func NewWebsocketHandler(
	gameService svc.GameService,
	playerService svc.PlayerService,
	gameController ctrl.GameController,

) http.Handler {
	return &websocketHandler{
		gameService:    gameService,
		playerService:  playerService,
		gameController: gameController,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		connections: make(map[string]map[string]*websocket.Conn),
	}
}

// |PRIVATE| //

func (h *websocketHandler) unmarshalAndValidate(
	data []byte, v dto.Validatable,
) error {
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("error unmarshaling data: %w", err)
	}
	return v.Validate()
}

func (h *websocketHandler) buildErrPopup(e error) ([]byte, error) {
	var htmlContent []byte
	tmpl, err := template.ParseFiles("views/game/error-popup.html")
	if err != nil {
		return htmlContent, err
	}
	type errorMessage struct{ Error string }
	eMsg := errorMessage{e.Error()}

	var buff bytes.Buffer
	err = tmpl.Execute(&buff, eMsg)
	if err != nil {
		return htmlContent, err
	}
	htmlContent = buff.Bytes()
	return htmlContent, nil
}

func (h *websocketHandler) createContext(
	r *http.Request,
) (*dto.WsContext, error) {
	var ctx dto.WsContext
	var err error

	// Get Game
	gameUUID, err := uuid.Parse(r.PathValue("gameUUID"))
	if err != nil {
		return &ctx, err
	}
	game, err := h.gameService.GetWithUUID(gameUUID)
	if err != nil {
		return &ctx, err
	}
	ctx.Game = game

	// Get Player
	playerUUIDCookie, err := r.Cookie(
		"player-uuid-" + ctx.Game.UUID.String())
	if err != nil {
		return &ctx, err
	}
	playerUUID, err := uuid.Parse(playerUUIDCookie.Value)
	if err != nil {
		return &ctx, err
	}
	player, err := h.playerService.GetWithUUID(playerUUID)
	if err != nil {
		return &ctx, err
	}
	ctx.Player = player

	return &ctx, nil
}

func (h *websocketHandler) refreshContext(ctx *dto.WsContext) error {
	if err := h.gameService.Refresh(ctx.Game); err != nil {
		return err
	}
	if err := h.playerService.Refresh(ctx.Player); err != nil {
		return err
	}
	return nil
}

func (h *websocketHandler) sendInitialData(
	ctx *dto.WsContext,
	conn *websocket.Conn,
) error {
	resultChars, err := h.gameController.GetAvaibleChars(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	resultFields, err := h.gameController.GetCurrentFields(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	initialResult := append(resultChars, resultFields...)
	fmt.Println(string(initialResult))
	if err := conn.WriteMessage(websocket.TextMessage, initialResult); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// |PUBLIC| //

func (h *websocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Create connection
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusUpgradeRequired) // 426
		return
	}
	defer conn.Close()

	ctx, err := h.createContext(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound) // 404
		return
	}

	sessionUUID := ctx.Game.UUID.String()
	connUUID := ctx.Player.UUID.String()
	_, ok := h.connections[sessionUUID]
	if !ok {
		h.connections[sessionUUID] = map[string]*websocket.Conn{}
	}
	h.connections[sessionUUID][connUUID] = conn

	h.sendInitialData(ctx, conn)

	for {
		messageType, p, err := conn.ReadMessage()
		log.Println(string(p))
		if err != nil {
			log.Println(err)
			delete(h.connections[sessionUUID], connUUID)
			return
		}
		err = h.refreshContext(ctx)
		if err != nil {
			log.Println(err)
			continue
		}
		action := &dto.ActionData{}
		err = h.unmarshalAndValidate(p, action)
		if err != nil {
			// TODO add here msg or smth
			log.Println(err)
			continue
		}

		var broadcastResponse []byte
		var senderResponse []byte
		switch action.Type {
		case "addTest":
			broadcastResponse, err = ctrl.GetRandomField()
		case "raiseExampleError":
			senderResponse, err = ctrl.GetExampleError()
		case "dismissError":
			senderResponse = []byte(`<div id="error-dialog"></div>`)
		case "getChars":
			senderResponse, err = h.gameController.GetAvaibleChars(ctx)
		case "makePlay":
			playData := &dto.PlayData{}
			err = h.unmarshalAndValidate(p, playData)
			if err != nil {
				break
			}
			broadcastResponse, senderResponse, err = h.gameController.ReceiveChars(ctx, playData)
			if err != nil {
				break
			}
			senderResponse, err = h.gameController.GetAvaibleChars(ctx)
		}

		if err != nil {
			broadcastResponse = nil
			senderResponse, err = h.buildErrPopup(err)
			if err != nil {
				log.Println(err)
				continue
			}
		}

		// fmt.Println(string(broadcastResponse))
		// fmt.Println(string(senderResponse))
		if broadcastResponse != nil {
			for _, c := range h.connections[sessionUUID] {
				if err := c.WriteMessage(messageType, broadcastResponse); err != nil {
					log.Println(err)
				}
			}
		}
		if senderResponse != nil {
			if err := conn.WriteMessage(messageType, senderResponse); err != nil {
				log.Println(err)
			}
		}
	}
}
