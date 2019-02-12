package atkinsdiet

// Machine is an Atkins Diet Slot Machine implementation.
// https://wizardofodds.com/games/slots/atkins-diet/
type Machine struct {
}

// SpinResult is a spin result structure.
type SpinResult struct {
	Type  string  `json:"type"`
	Total int64   `json:"total"`
	Stops [3]int8 `json:"stops"`
}

// Spin does a spin with bet and returns its total win and spins.
// Spins are []SpinResult but returnes as []interface{} to satisfy a generic interface.
func (m *Machine) Spin(bet int64) (int64, []interface{}) {
	return 0, []interface{}{SpinResult{}}
}
