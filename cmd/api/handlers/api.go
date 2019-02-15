package handlers

import (
	"net/http"

	"github.com/ishmulyan/slotssvc/cmd/api/middleware"
	"github.com/ishmulyan/slotssvc/internal/jwt"
	"github.com/julienschmidt/httprouter"
)

// API returns an API handler.
func API(jwtSecret []byte) http.Handler {
	router := httprouter.New()

	jwtSvc := jwt.NewService(jwtSecret)
	slotMiddleware := middleware.Slot{jwtSvc}

	atkinsdiet := NewAtkinsDiet()
	router.POST("/api/machines/atkins-diet/spins", slotMiddleware.SpinHandler(atkinsdiet.Spin))

	return router
}
