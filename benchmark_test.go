package main

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkSimulation(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	rueel := newRueelBuilder().
		startWith(100_000).
		use(halfOnYellowStrategy).
		stopIfScrapsReach(1_000_000_000_000_000_000).
		stopIfNIterationIs(1_000_000).
		build()

	// rueel := rueel{
	// 	scraps:       1000,
	// 	scrapsGoal:   1_000_000,
	// 	strategy:     halfOnYellowStrategy,
	// 	maxIteration: 1,
	// 	nIteration:   0,
	// 	analytics: analytics{
	// 		outcomeDistribution: map[color]int{YELLOW: 0, GREEN: 0, BLUE: 0, PURPLE: 0, RED: 0},
	// 	},
	// }

	for i := 0; i < b.N; i++ {
		// rueel := newRueelBuilder().
		// 	startWith(1000).
		// 	use(halfOnYellowStrategy).
		// 	stopIfScrapsReach(1_000_000).
		// 	stopIfNIterationIs(1).
		// 	build()

		// rueel := rueel{
		// 	scraps:       1000,
		// 	scrapsGoal:   1_000_000,
		// 	strategy:     halfOnYellowStrategy,
		// 	maxIteration: 1,
		// 	nIteration:   0,
		// 	analytics: analytics{
		// 		outcomeDistribution: map[color]int{YELLOW: 0, GREEN: 0, BLUE: 0, PURPLE: 0, RED: 0},
		// 	},
		// }

		rueel.simulate()
		// fmt.Println(rueel.report())
	}
}
