package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
)

type DonkeyKong struct {
	PlayerName string `json:"player_name"`
	Score      uint32 `json:"score"`
}

func draw_donkey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	HighScores := []DonkeyKong{}

	P1 := DonkeyKong{PlayerName: "RRD", Score: 123123}
	HighScores = append(HighScores, DonkeyKong{PlayerName: "RRD", Score: 234234})
	HighScores = append(HighScores, DonkeyKong{PlayerName: "PHX", Score: 323234})
	HighScores = append(HighScores, DonkeyKong{PlayerName: "ALX", Score: 3323234})
	HighScores = append(HighScores, P1)

	sort.Slice(HighScores, func(i, j int) bool {
		return HighScores[i].Score > HighScores[j].Score
	})

	err := json.NewEncoder(w).Encode(HighScores)
	if err != nil {
		_, err = io.WriteString(w, err.Error())
		if err != nil {
			fmt.Println("Oops!!! {s}", err.Error())
		}
	}
}

func main() {
	fmt.Println("Running web serivce")
	http.HandleFunc("/api/v1/donkey", draw_donkey)
	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		fmt.Println("Oops {s}", err.Error())
	}
}
