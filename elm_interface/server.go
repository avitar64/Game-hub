package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/avitar64/Boost-bot/games"
)

func marshalMastermindGHame(game *games.MastermindGame) (string, error) {
	b, err := json.Marshal(game)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

type MastermindGuessPOSTRequest struct {
	Game  *games.MastermindGame `json:"game"`
	Guess [4]games.MasterColor  `json:"guess"`
}

func handlePOST(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var guessStruct = &MastermindGuessPOSTRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(guessStruct)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	guessStruct.Game.Guess(guessStruct.Guess)

	jsonData, err := marshalMastermindGHame(guessStruct.Game)
	if err != nil {
		log.Println("!! ERROR: ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Fprint(w, jsonData)
}

func handleMastermindRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	if r.Method == "GET" {
		log.Println("get")
		handleGET(w, r)
	} else if r.Method == "POST" {
		log.Println("post")
		handlePOST(w, r)
	}
}

func handleGET(w http.ResponseWriter, r *http.Request) {
	game := games.NewMastermindGame()

	jsonData, err := marshalMastermindGHame(game)
	if err != nil {
		log.Println("!! ERROR: ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Fprint(w, jsonData)
}

func main() {
	http.HandleFunc("/mastermind", handleMastermindRequests)
	http.Handle("/", http.FileServer(http.Dir(".")))

	log.Println("Running...")

	http.ListenAndServe(":8080", nil)
}
