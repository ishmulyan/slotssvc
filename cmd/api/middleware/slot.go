package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ishmulyan/slotssvc/internal/jwt"
	"github.com/julienschmidt/httprouter"
)

type contextKey string

const (
	// CtxBetKey is a key under which bet from the request (jwt) put in the context.
	CtxBetKey = contextKey("bet")
	// CtxSpinResponse is a key under which SpinResponse put in the context.
	CtxSpinResponse = contextKey("spinResponse")
)

// SpinResponse is a response for slot spin.
type SpinResponse struct {
	Total int64       `json:"total"`
	Spins interface{} `json:"spins"`
	JWT   string      `json:"jwt"`
	Error error       `json:"-"`
}

// Slot is slot machine middleware set.
type Slot struct {
	JWTSvc jwt.Service
}

// SpinHandler returns slot machine middleware.
// Parses JWT token from response before handling spin.
// Creates JWT token after spin handled and repsponses with SpinResponse.
func (m *Slot) SpinHandler(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("can't read request body: %v", err), http.StatusBadRequest)
			return
		}

		// parse JWT
		tokenClaims, err := m.JWTSvc.Decode(string(body))
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

		ctx := r.Context()
		ctx = context.WithValue(ctx, CtxBetKey, tokenClaims.Bet)
		spinResponse := SpinResponse{}
		ctx = context.WithValue(ctx, CtxSpinResponse, &spinResponse)

		//
		next(w, r.WithContext(ctx), ps)
		//

		if spinResponse.Error != nil {
			http.Error(w, spinResponse.Error.Error(), http.StatusInternalServerError)
			return
		}

		tokenClaims.Chips += spinResponse.Total

		// create new JWT
		tokenString, err := m.JWTSvc.Encode(tokenClaims)
		if err != nil {
			http.Error(w, fmt.Sprintf("can't encode jwt token: %v", err), http.StatusInternalServerError)
			return
		}

		spinResponse.JWT = tokenString

		js, err := json.Marshal(spinResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
