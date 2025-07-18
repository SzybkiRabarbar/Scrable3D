package dto

import "scrable3/internal/model"

type WsContext struct {
	Game   *model.Game
	Player *model.Player
}
