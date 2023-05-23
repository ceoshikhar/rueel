package main

import (
	"fmt"
	"math/rand"
	"strings"
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

func (c color) worth() int {
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

func (c color) string() string {
	switch c {
	case YELLOW:
		return "Yellow"
	case GREEN:
		return "Green"
	case BLUE:
		return "Blue"
	case PURPLE:
		return "Purple"
	case RED:
		return "Red"
	default:
		return ""
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

type analytics struct {
	// higestScraps is the maxinum number of scraps that we had at one point.
	higestScraps int
	// outcomeDistribution is a map where key is a color and value is the number
	// of times that color was the outcome of a spin.
	outcomeDistribution map[color]int
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

	// analytics is the analytics for a given simulation.
	analytics analytics
}

func (r *rueel) simulate() {

	for {
		if r.scraps > r.analytics.higestScraps {
			r.analytics.higestScraps = r.scraps
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
		r.nIteration += 1
	}

}

func (r rueel) spin(s Strategy) (scrapsWon int, stopIterating bool) {
	if r.scraps <= 1 {
		stopIterating = true
		return scrapsWon, stopIterating
	}
	bet := s(r.scraps)

	wheelSlotThatWon := rand.Intn(len(WHEEL)-0) + 0
	colorThatWon := WHEEL[wheelSlotThatWon]

	count, ok := r.analytics.outcomeDistribution[colorThatWon]
	if !ok {
		r.analytics.outcomeDistribution[colorThatWon] = 1
	} else {
		r.analytics.outcomeDistribution[colorThatWon] = count + 1
	}

	wagerOnColorThatWon := bet[colorThatWon]
	rewardOnColorThatWon := wagerOnColorThatWon + (wagerOnColorThatWon * colorThatWon.worth())
	scrapsWon = rewardOnColorThatWon - (bet.total() - wagerOnColorThatWon)

	// fmt.Printf("%d Bet %+v - Color %s - Wagered %d - Won/Loss %d - Total %d \n", r.nIteration, bet, colorThatWon.string(), bet.total(), scrapsWon, r.scraps+scrapsWon)

	return scrapsWon, stopIterating
}

func (r rueel) report() string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("Scraps: %d \n", r.scraps))
	sb.WriteString(fmt.Sprintf("Total spins: %d \n", r.nIteration))
	sb.WriteString(fmt.Sprintf("Highest scraps reached: %d \n", r.analytics.higestScraps))

	od := r.analytics.outcomeDistribution

	sb.WriteString("Sping outocme distribution: {")
	colorSequence := [TOTAL_COLORS_AVAIL]color{YELLOW, GREEN, BLUE, PURPLE, RED}
	for i, color := range colorSequence {
		count := od[color]
		percentage := float64(count) / float64(r.nIteration) * 100
		sb.WriteString(fmt.Sprintf(" %s: %d (%.2f%%)", color.string(), count, percentage))

		if i != TOTAL_COLORS_AVAIL-1 {
			sb.WriteString(",")
		}
	}

	return sb.String()
}

type rueelBuilder struct {
	rueel rueel
}

func newRueelBuilder() rueelBuilder {
	return rueelBuilder{
		rueel: rueel{
			scraps:       1000,
			strategy:     defaultStrategy,
			maxIteration: 10,
			analytics: analytics{
				outcomeDistribution: map[color]int{YELLOW: 0, GREEN: 0, BLUE: 0, PURPLE: 0, RED: 0},
			},
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

func (rb rueelBuilder) stopIfScrapsReach(scrapsGoal int) rueelBuilder {
	rb.rueel.scrapsGoal = scrapsGoal
	return rb
}

func (rb rueelBuilder) stopIfNIterationIs(maxIteration int) rueelBuilder {
	rb.rueel.maxIteration = maxIteration
	return rb
}

func (rb rueelBuilder) build() rueel {
	return rb.rueel
}

func doSimulation() {
	// rueel := newRueelBuilder().
	// 	startWith(1000).
	// 	use(halfOnYellowStrategy).
	// 	stopIfScrapsReach(1_000_000).
	// 	stopIfNIterationIs(100).
	// 	build()

	rueel := newRueelBuilder().
		startWith(1000).
		use(halfOnYellowStrategy).
		stopIfScrapsReach(1_000_000_000_000_000_000).
		stopIfNIterationIs(10_000_000).
		build()

	rueel.simulate()

	fmt.Println(rueel.report())
}

func main() {
	// Seed the RNG.
	rand.Seed(time.Now().UnixNano())
	doSimulation()
}
