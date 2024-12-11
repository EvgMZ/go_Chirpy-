package handler

import (
	"chirpy/reponse"
	"encoding/json"
	"net/http"
	"regexp"
)

func replaceProfaneWords(msg string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	replacement := "****"
	for _, word := range profaneWords {
		re := regexp.MustCompile(`\b(?i)` + word + `\b`)
		msg = re.ReplaceAllString(msg, replacement)
	}
	return msg
}
func NewValidateChripHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Req string `json:"body"`
	}
	type response struct {
		CleanedBody string `json:"cleaned_body"`
	}
	decoder := json.NewDecoder(r.Body)
	res := request{}
	err := decoder.Decode(&res)
	if len(res.Req) > 140 {
		reponse.RespondWithError(w, 400, "Chirp is too long")
		return
	}
	if err != nil {
		reponse.RespondWithError(w, 500, "Something went wrong")
		return
	}
	cleanedBody := replaceProfaneWords(res.Req)
	reponse.RespondWithJSON(w, 200, response{CleanedBody: cleanedBody})
}
