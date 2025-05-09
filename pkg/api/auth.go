package api

import (
	"encoding/json"
	"net/http"

	"github.com/langowen/go_final_project/pkg/auth"
	"github.com/langowen/go_final_project/pkg/config"
)

type SignInRequest struct {
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

func SignInHandler(cfg *config.Config, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
		return
	}

	var req SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "неверный JSON формат")
		return
	}

	if cfg.Token == "" {
		respondWithError(w, http.StatusInternalServerError, "отсутствует токен")
		return
	}

	if req.Password != cfg.Token {
		respondWithError(w, http.StatusUnauthorized, "неверный пароль")
		return
	}

	token, err := auth.GenerateToken(cfg)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "не удалось сгенерировать токен")
		return
	}

	respondWithJSON(w, http.StatusOK, SignInResponse{Token: token})
}
