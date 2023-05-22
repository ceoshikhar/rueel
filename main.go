package main

import (
	"fmt"
	"math/rand"
	"time"
)

type color int

const (
	// 1
	YELLOW color = iota
	// 3
	GREEN
	// 5
	BLUE
	// 10
	PURPLE
	// 20
	RED
)

// TOTAL_COLORS_AVAIL is the total kinds of colors that exists in the Rueel.
// The count comes from the number of variants in the `Color` enum.
const TOTAL_COLORS_AVAIL = 5

// WHEEL reprents each slot with a Color. There are total 25 slots.
// YELLOW => 12, GREEN => 6, BLUE => 4, PURPLE => 2, RED => 1
var WHEEL = [25]color{
	RED,
	YELLOW,
	GREEN,
	YELLOW,
	BLUE,
	YELLOW,
	GREEN,
	YELLOW,
	PURPLE,
	YELLOW,
	GREEN,
	YELLOW,
	BLUE,
	YELLOW,
	BLUE,
	GREEN,
	YELLOW,
	PURPLE,
	YELLOW,
	GREEN,
	YELLOW,
	BLUE,
	YELLOW,
	GREEN,
	YELLOW,
}

func (c color) int() int {
	switch c {
	case YELLOW:
		return 1
	case GREEN:
		return 3
	case BLUE:
		return 5
	case PURPLE:
		return 10
	case RED:
		return 20
	default:
		return 0
	}
}

type bet map[color]int

func NewBet() bet {
	bet := bet{}

	bet[YELLOW] = 0
	bet[GREEN] = 0
	bet[BLUE] = 0
	bet[PURPLE] = 0
	bet[RED] = 0

	return bet
}

func (b bet) total() int {
	if len(b) != TOTAL_COLORS_AVAIL {
		panic("Bet.Total: Bet object have some Color variants missing")
	}

	var total int = 0

	for _, scraps := range b {
		total += scraps
	}

	return total
}

type Strategy func(scraps int) bet

// defaultStrategy bets all on YELLOW.
func defaultStrategy(scraps int) bet {
	bet := NewBet()
	bet[YELLOW] = scraps
	return bet
}

// halfOnYellowStrategy strategy will keep betting half the scraps on YELLOW.
func halfOnYellowStrategy(scraps int) bet {
	bet := NewBet()
	bet[YELLOW] = scraps / 2
	return bet
}

// rueel represents the Rust's wheel famous for gambling.
type rueel struct {
	// scraps is the amount of scraps available to bet.
	scraps int
	// scrapsGoal is the amount of scraps at which we stop the simulation.
	scrapsGoal int
	// strategy is what strategy you want to use for a simulation `Rueel.Simulation`.
	strategy Strategy
	// maxIteration is the iteration count when we want to stop the simulation.
	// 0 means it will keep spinning until the scraps reach 0 or `scrapsGoal`.
	maxIteration int

	// nIteration is the last iteration that happened in simulation.
	nIteration int
	// mostScraps is the maxinum number of scraps that we had at one point.
	mostScraps int
}

func (r *rueel) simulate() {

	for {
		if r.scraps > r.mostScraps {
			r.mostScraps = r.scraps
		}

		if r.maxIteration > 0 && r.nIteration >= r.maxIteration {
			break
		}

		if r.scrapsGoal > 0 && r.scraps >= r.scrapsGoal {
			break
		}

		scrapsWonOrLost, stopIterating := r.spin(r.strategy)
		if stopIterating {
			break
		}

		r.scraps += scrapsWonOrLost

		// fmt.Printf("Total %d \n", r.scraps)

		if r.scraps <= 0 {
			break
		}

		r.nIteration += 1
	}

	r.report()
}

func (r rueel) spin(s Strategy) (scrapsEarned int, stopIterating bool) {
	if r.scraps <= 1 {
		return 0, true
	}

	bet := s(r.scraps)

	wheelSlotThatWon := rand.Intn(len(WHEEL)-0) + 0
	colorThatWon := WHEEL[wheelSlotThatWon]
	betMadeOnColorThatWon := bet[colorThatWon]

	scrapsWon := betMadeOnColorThatWon + (betMadeOnColorThatWon * colorThatWon.int())

	// fmt.Printf("Wagered %d Won %d ", bet.total(), scrapsWon)

	return scrapsWon - (bet.total() - betMadeOnColorThatWon), false
}

func (r rueel) report() {
	fmt.Println("############# Rueel simulation completed ############# ")
	fmt.Printf("%+v\n", r)
}

type rueelBuilder struct {
	rueel rueel
}

func newRueelBuilder() rueelBuilder {
	return rueelBuilder{
		rueel: rueel{
			scraps:       1000,
			scrapsGoal:   100000,
			strategy:     defaultStrategy,
			maxIteration: 100,
			nIteration:   0,
			mostScraps:   1000,
		},
	}
}

func (rb rueelBuilder) startWith(scraps int) rueelBuilder {
	rb.rueel.scraps = scraps
	return rb
}

func (rb rueelBuilder) use(strategy Strategy) rueelBuilder {
	rb.rueel.strategy = strategy
	return rb
}

func (rb rueelBuilder) stopWhen(scrapsGoal int, maxIteration int) rueelBuilder {
	rb.rueel.scrapsGoal = scrapsGoal
	rb.rueel.maxIteration = maxIteration
	return rb
}

func (rb rueelBuilder) build() rueel {
	return rb.rueel
}

func main() {
	// Seed the RNG.
	rand.Seed(time.Now().UnixNano())

	rueel := newRueelBuilder().
		startWith(4096).
		use(halfOnYellowStrategy).
		stopWhen(1_000_00, 1000).build()

	rueel.simulate()
}
