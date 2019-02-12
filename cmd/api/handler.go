package main

import (
	"encoding/json"
	"log"
	"net/http"

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
func handler(log *log.Logger) http.Handler {
	router := httprouter.New()

	// register slot machines in a map by their ids.
	slotMachines := map[string]slotMachine{
		"atkins-diet": &atkinsdiet.Machine{},
	}
	router.POST("/api/machines/:id/spins", spinHandler(slotMachines))

	return router
}

func spinHandler(slotMachines map[string]slotMachine) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		sm, ok := slotMachines[ps.ByName("id")]
		if !ok {
			// no game with such id
			http.Error(w, "game not found", http.StatusNotFound)
			return
		}

		// parse jwt here
		// decrease chips by bet size, if chips are negative -- return an error

		total, spins := sm.Spin(1000)

		// create new jwt here
		jwt := "new jwt"

		res := struct {
			Total int64         `json:"total"`
			Spins []interface{} `json:"spins"`
			JWT   string        `json:"jwt"`
		}{
			Total: total,
			Spins: spins,
			JWT:   jwt,
		}

		js, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

