package health

import "net/http"

// Read godoc
//
// @Summary Health Check
// @Description Health Check
// @Tags Health
// @Success 200
// @Router /health [get]
func Read(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server is healthy"))
}
