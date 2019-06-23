package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Board Game Struct
type BoardGame struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author Struct
type Author struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

// Init var as a slice struct
var boardGames []BoardGame

// Get all board games
func getBoardGames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(boardGames)
}

// Get single board game
func getBoardGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // get params
	// Find id
	for _, item := range boardGames {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&BoardGame{})
}

func createBoardGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var boardGame BoardGame
	_ = json.NewDecoder(r.Body).Decode(&boardGame)
	boardGame.ID = strconv.Itoa(rand.Intn(10000000)) // not safe, just mocking
	boardGames = append(boardGames, boardGame)
	json.NewEncoder(w).Encode(boardGame)
}

func updateBoardGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // get params
	for index, item := range boardGames {
		if item.ID == params["id"] {
			var boardGame BoardGame
			_ = json.NewDecoder(r.Body).Decode(&boardGame)
			boardGame.ID = params["id"]

			boardGames = append(boardGames[:index], boardGames[index+1:]...) // Delete old one
			boardGames = append(boardGames, boardGame)                       // Add new one
			json.NewEncoder(w).Encode(boardGame)
			return
		}
	}
	json.NewEncoder(w).Encode(boardGames)
}

func deleteBoardgame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // get params
	for index, item := range boardGames {
		if item.ID == params["id"] {
			boardGames = append(boardGames[:index], boardGames[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(boardGames)
}

func main() {
	router := mux.NewRouter()

	// Mock Data
	boardGames = append(boardGames, BoardGame{ID: "1", Isbn: "12345", Title: "Rock Paper Wizard", Author: &Author{FirstName: "Dungeon", LastName: "Dragons"}})
	boardGames = append(boardGames, BoardGame{ID: "2", Isbn: "22345", Title: "Blokus", Author: &Author{FirstName: "Matter", LastName: "Games"}})

	// Route Handlers
	router.HandleFunc("/api/boardgames", getBoardGames).Methods("GET")
	router.HandleFunc("/api/boardgames/{id}", getBoardGame).Methods("GET")
	router.HandleFunc("/api/boardgames", createBoardGame).Methods("POST")
	router.HandleFunc("/api/boardgames/{id}", updateBoardGame).Methods("PUT")
	router.HandleFunc("/api/boardgames/{id}", deleteBoardgame).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}
