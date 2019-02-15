package atkinsdiet

import "errors"

const (
	spinTypeMain = "main"
	spinTypeFree = "free"
)

var (
	// ErrNonPositiveBet is an error value returned by Machine.Spin method when bet is non positive.
	ErrNonPositiveBet = errors.New("bet is not positive")
	// ErrBadNLines is an error value returned by Machine.Spin method when lines isn't in (0,20]).
	ErrBadNLines = errors.New("number of lines must be in (0,20]")
)

// Machine is an Atkins Diet Slot Machine implementation.
// https://wizardofodds.com/games/slots/atkins-diet/
type Machine struct {
	rand *randomizer
}

// SpinStats represents statistics. Used mostly for simulation purpose.
type SpinStats struct {
	LinePays int64
	Scatter  int64
	Bonus    int64
	Total    int64
}

// SpinResult is a spin result structure.
type SpinResult struct {
	Type  string  `json:"type"`
	Total int64   `json:"total"`
	Stops [5]int8 `json:"stops"`
}

// New creates atkins diet Machine.
func New() *Machine {
	return &Machine{
		rand: newRandomizer(),
	}
}

// Spin does a spin with bet per line 'bet' and number of lines 'nline', returns its total win and spins.
func (m *Machine) Spin(bet, lines int) (SpinStats, []SpinResult, error) {
	var st SpinStats

	if bet < 0 {
		return st, nil, ErrNonPositiveBet
	}

	if lines < 0 {
		return st, nil, ErrBadNLines
	}
	if len(paylines) < lines {
		return st, nil, ErrBadNLines
	}

	st, sr := m.spin(bet, lines)
	return st, sr, nil
}

// spin does spin loop in case there're free spins.
func (m *Machine) spin(bet, lines int) (SpinStats, []SpinResult) {
	// main spin
	stops := m.reelStops()
	symbols := mapSymbolsToStops(stops)
	baseLineWin := caclculateLinesWin(symbols, bet, lines)
	// scatterMultiplier := int64(0)
	scatterMultiplier := determineScattersMultiplier(symbols)
	freeSpins := 0
	var scattersWin int64
	if 0 < scatterMultiplier {
		scattersWin = int64(bet) * int64(lines) * scatterMultiplier
		freeSpins = freeSpinsAmount
	}

	stats := SpinStats{
		LinePays: baseLineWin,
		Scatter:  scattersWin,
		Total:    baseLineWin + scattersWin,
	}
	spins := []SpinResult{
		{
			Type:  spinTypeMain,
			Total: stats.Total,
			Stops: stops[1],
		},
	}

	// free spins
	for i := 0; i < freeSpins; i++ {
		stops = m.reelStops()
		symbols = mapSymbolsToStops(stops)
		baseLineWin = caclculateLinesWin(symbols, bet, lines) * freeSpinsMultiplier

		scattersWin = 0
		scatterMultiplier = determineScattersMultiplier(symbols)
		if 0 < scatterMultiplier {
			scattersWin = int64(bet) * int64(lines) * scatterMultiplier * freeSpinsMultiplier
			freeSpins += freeSpinsAmount
		}

		spinWin := baseLineWin + scattersWin
		spins = append(spins, SpinResult{
			Type:  spinTypeFree,
			Total: spinWin,
			Stops: stops[1],
		})

		stats.Bonus += spinWin
		stats.Total += spinWin
	}
	return stats, spins
}

// reelStops randomly chooses and returns reel stops.
func (m *Machine) reelStops() [3][5]int8 {
	var stops [3][5]int8

	// randomly choose reel stops for the top row
	for i := range stops[0] {
		stops[0][i] = int8(m.rand.Intn(32))
	}

	// populate two other rows respectfully
	for i := 1; i < 3; i++ {
		for j := range stops[i] {
			if stops[i-1][j] == 31 {
				stops[i][j] = 0
				continue
			}
			stops[i][j] = stops[i-1][j] + 1
		}
	}
	return stops
}

// mapSymbolsToStops maps symbols to reel stops.
func mapSymbolsToStops(stops [3][5]int8) [3][5]symbol {
	var symbols [3][5]symbol
	for reelID := 0; reelID < 5; reelID++ {
		for rowID := 0; rowID < 3; rowID++ {
			stop := stops[rowID][reelID]
			symbols[rowID][reelID] = reelStrips[reelID][stop]
		}
	}
	return symbols
}

// caclculateLinesWin calculates win amount for symbols with bet and lines.
func caclculateLinesWin(symbols [3][5]symbol, bet, lines int) int64 {
	var (
		win  int64
		line [5]symbol
	)
	for i := 0; i < lines; i++ {
		payline := paylines[i][:]
		for reelID := 0; reelID < 5; reelID++ {
			line[reelID] = symbols[payline[reelID]][reelID]
		}
		win += int64(bet) * int64(winForLine(line))
	}
	return win
}

// winForLine retuns a single win for the line.
func winForLine(line [5]symbol) int16 {
	numPays := len(paytable)
	for payID := 0; payID < numPays; payID++ {
		if checkLine(line, paytable[payID].n, paytable[payID].symbol) {
			return paytable[payID].win
		}
	}
	return 0
}

// checkLine check whether the line has n sym in a row.
func checkLine(line [5]symbol, n int8, sym symbol) bool {
	for i := int8(0); i < n; i++ {
		// symAtkins is wild so it is allowed to be in the line
		if line[i] != sym && line[i] != symAtkins {
			return false
		}
	}
	return true
}

// determineScattersMultiplier returns scatter multiplier.
func determineScattersMultiplier(symbols [3][5]symbol) int64 {
	var n int64
	for _, row := range symbols {
		for _, s := range row {
			if s == symScale {
				n++
			}
		}
	}
	switch n {
	case 3:
		return 5
	case 4:
		return 25
	case 5:
		return 100
	default:
		return 0
	}
}
