package handlers

import (
	"fmt"
	"net/http"

	"github.com/ishmulyan/slotssvc/cmd/api/middleware"
	"github.com/ishmulyan/slotssvc/pkg/atkinsdiet"
	"github.com/julienschmidt/httprouter"
)

// AtkinsDiet atkins-diet slot machine API handler set.
type AtkinsDiet struct {
	machine *atkinsdiet.Machine
}

// NewAtkinsDiet creates atkins-diet slot machine API handler set.
func NewAtkinsDiet() *AtkinsDiet {
	return &AtkinsDiet{
		machine: atkinsdiet.New(),
	}
}

// Spin handles spin actions.
func (h *AtkinsDiet) Spin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	bet := r.Context().Value(middleware.CtxBetKey).(int64)
	res := r.Context().Value(middleware.CtxSpinResponse).(*middleware.SpinResponse)

	lines := 20
	betPerLine := int(bet / 4)
	st, sr, err := h.machine.Spin(betPerLine, lines)
	if err != nil {
		res.Error = fmt.Errorf("spin failed: %v", err)
		return
	}

	res.Total = st.Total
	res.Spins = sr
}
