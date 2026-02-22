package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func handleChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}
	dec := json.NewDecoder(r.Body)
	params := parameters{}

	if err := dec.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	cleaned := cleanMessage(params.Body)

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleaned,
	})
}

func cleanMessage(msg string) string {
	words := strings.Split(msg, " ")
	for i, word := range words {
		word = strings.ToLower(word)
		if word == "kerfuffle" || word == "sharbert" || word == "fornax" {
			words[i] = "****"
		}
	}
	msg = strings.Join(words, " ")
	return msg
}
