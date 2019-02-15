package main

import (
	"log"
	"sync"

	"github.com/ishmulyan/slotssvc/pkg/atkinsdiet"
)

func main() {
	const (
		n           = 1e7
		bet         = 1
		lines       = 20
		concurrency = 20
	)

	log.Println("Simulations started...")
	log.Printf("concurrency:	%v", concurrency)
	log.Printf("spins:		%v", n)
	log.Printf("bet:		%v", bet)
	log.Printf("lines:		%v", lines)

	m := atkinsdiet.New()

	spinStats := make(chan atkinsdiet.SpinStats)

	go spin(m, bet, lines, n, concurrency, spinStats)

	var stats atkinsdiet.SpinStats
	for s := range spinStats {
		stats.LinePays += s.LinePays
		stats.Scatter += s.Scatter
		stats.Bonus += s.Bonus
		stats.Total += s.Total
	}

	totalBets := bet * lines * int64(n)

	log.Println("Simulation finished")
	log.Printf("Total Bets:		%d", totalBets)
	log.Printf("Total Wins:		%d", stats.Total)
	log.Printf("Total RTP:		%f", float64(stats.Total)/float64(totalBets))

	linePays := float64(stats.LinePays) / float64(totalBets)
	scatter := float64(stats.Scatter) / float64(totalBets)
	bonus := float64(stats.Bonus) / float64(totalBets)

	log.Printf("Line Pays RTP:	%f", linePays)
	log.Printf("Scatter RTP: 	%f", scatter)
	log.Printf("Bonus RTP: 		%f", bonus)
	log.Printf("Sum RTP: 		%f", linePays+scatter+bonus)
}

func spin(m *atkinsdiet.Machine, bet, lines, n, concurrency int, stats chan atkinsdiet.SpinStats) {
	spins := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for range spins {
				st, _, _ := m.Spin(bet, lines)
				stats <- st
			}
		}()
	}

	for i := 0; i < n; i++ {
		spins <- struct{}{}
	}
	close(spins)

	wg.Wait()
	close(stats)
}
