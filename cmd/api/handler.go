package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ishmulyan/slotssvc/internal/jwt"
	"github.com/ishmulyan/slotssvc/pkg/atkinsdiet"
	"github.com/julienschmidt/httprouter"
)

// slotMachine is an interface for slot machine.
// Spin method is designed to return spins as []interface{} to support different games,
// possibly with different spins structure. e.g. free spins can hold a multiplier.
type slotMachine interface {
	Spin(bet int64) (int64, []interface{})
}

// handler returns an API handler.
func handler(jwtSecret []byte) http.Handler {
	router := httprouter.New()

	jwtSvc := jwt.NewService(jwtSecret)

	// register slot machines in a map by their ids.
	slotMachines := map[string]slotMachine{
		"atkins-diet": &atkinsdiet.Machine{},
	}
	slotMachineH := slotMachineHandler{
		slotMachines: slotMachines,
		jwtSvc:       jwtSvc,
	}
	router.POST("/api/machines/:id/spins", slotMachineH.spin)

	return router
}

type slotMachineHandler struct {
	slotMachines map[string]slotMachine
	jwtSvc       jwt.Service
}

func (h *slotMachineHandler) spin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sm, ok := h.slotMachines[ps.ByName("id")]
	if !ok {
		// no game with such id
		http.Error(w, "game not found", http.StatusNotFound)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("can't read request body: %v", err), http.StatusBadRequest)
		return
	}

	// parse jwt
	tokenClaims, err := h.jwtSvc.Decode(string(body))
	if err != nil {
		http.Error(w, fmt.Sprintf("can't parse jwt token: %v", err), http.StatusBadRequest)
		return
	}

	if tokenClaims.Bet <= 0 {
		http.Error(w, "bet should be greater 0", http.StatusBadRequest)
		return
	}

	// decrease chips by bet size, if chips are negative -- return an error
	tokenClaims.Chips -= tokenClaims.Bet
	if tokenClaims.Chips < 0 {
		http.Error(w, "not enough chips to make a spin", http.StatusBadRequest)
		return
	}

	// make a spin
	win, spins := sm.Spin(tokenClaims.Bet)
	tokenClaims.Chips += win

	tokenString, err := h.jwtSvc.Encode(tokenClaims)
	if err != nil {
		http.Error(w, fmt.Sprintf("can't encode jwt token: %v", err), http.StatusInternalServerError)
		return
	}

	res := struct {
		Total int64         `json:"total"`
		Spins []interface{} `json:"spins"`
		JWT   string        `json:"jwt"`
	}{
		Total: win,
		Spins: spins,
		JWT:   tokenString,
	}

	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
