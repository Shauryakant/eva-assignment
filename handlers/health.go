package handlers

import (
	"net/http"

	"ticket-system/utils"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.JSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
