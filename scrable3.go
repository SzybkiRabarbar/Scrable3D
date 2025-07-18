package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"scrable3/internal/ctrl"
	"scrable3/internal/handler"
	"scrable3/internal/repo"
	"scrable3/internal/svc"
)

const port string = ":8080"
const databaseFile = "sqlite.db"

func main() {
	wordsController, err := ctrl.NewWordsController("words/words_alpha.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	os.Remove(databaseFile)
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	repo, err := repo.NewSqlite3Connection(db)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer repo.CloseConn()

	mux := http.NewServeMux()
	static := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	mux.Handle("/static/", static)
	styles := http.StripPrefix("/styles/", http.FileServer(http.Dir("styles")))
	mux.Handle("/styles/", styles)
	scripts := http.StripPrefix("/scripts/", http.FileServer(http.Dir("scripts")))
	mux.Handle("/scripts/", scripts)

	gameService := svc.NewGameService(repo)
	playerService := svc.NewPlayerService(repo)
	fieldService := svc.NewFieldService(repo)
	avCharService := svc.NewAvCharService(repo)

	gameController := ctrl.NewGameController(
		wordsController,
		playerService,
		fieldService,
		avCharService,
	)

	homeHandler := handler.NewHomeHandler()
	gameHandler := handler.NewGameHandler(gameService, playerService, fieldService)
	websocketHandler := handler.NewWebsocketHandler(
		gameService,
		playerService,
		gameController,
	)

	mux.Handle("/", homeHandler)
	mux.Handle("/game", gameHandler)
	mux.Handle("/game/{gameUUID}", gameHandler)
	mux.Handle("/ws/{gameUUID}", websocketHandler)

	fmt.Printf("Start server, port %v\n", port)
	http.ListenAndServe(port, mux)
}
