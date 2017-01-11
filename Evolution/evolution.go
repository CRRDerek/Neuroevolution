package evolution

import (
	"github.com/CRRDerek/Neuroevolution/Games"
)

// Given a game, a player factory function, (specific to that game) the number of
// generations, maximum number of games, and an initial population of agents
//(which must all be the same type!) run an evolutionary algorithm and return
// the best agent.
func EvolveAgents(g games.Game, playerMaker games.PlayerMaker, generations int,
	max_games int, pop []games.Agent) games.Agent {

	// Initialize an array of channels for each member of the population
	fitness_channels := make([]chan int, len(pop))
	fitness_values := make([]int, len(pop))

	// Initialize an array for the new population
	var new_pop []games.Agent

	// Initialize variables for max fitness and max agent
	max_fitness := -9999999999
	var max_agent games.Agent

	//Initialize each channel
	for i := 0; i < len(pop); i++ {
		fitness_channels[i] = make(chan int)
	}

	// Loop the algorithm for as many iterations are specified in the number of
	// generations.
	i := 0
	for {
		// Start a goroutine to test each member of the population.
		for j := 0; j < len(pop); j++ {
			go func() {
				score := 0
				player := playerMaker(pop[j])
				// Keep testing this player until the maximum number of games is
				// reached.
				for k := 0; k < max_games; k++ {
					// Play the game against a random opponent
					switch games.PlayerTrial(g, player) {
					// If the agent player wins, reward it
					case 1:
						score += 1
					// Reward draws too.
					case 0:
						score += 1
					// If they lose, break out of the loop.
					case -1:
						k = max_games
					}
				}

				// Send the score over the appropriate channel
				fitness_channels[j] <- score
			}()
		}

		// Receive fitness values from channels and find the maximum fitness
		max_fitness = -9999999
		for j := 0; j < generations; j++ {
			fitness_values[j] = <-fitness_channels[j]
			if fitness_values[j] > max_fitness {
				max_fitness = fitness_values[j]
				max_agent = pop[j]
			}
		}

		// Iterate the generation number and return if the algorithm is complete.
		i++
		if i >= generations {
			return max_agent
		}

		// Create a new array for the new population
		new_pop = make([]games.Agent, len(pop))
		new_pop[0] = max_agent

		// Create the next generation by mating based on fitness
		for j := 1; j < len(pop); j++ {
			p1 := weighted_selection(pop, fitness_values)
			p2 := weighted_selection(pop, fitness_values)
			new_pop[j] = p1.(games.Agent).Mate(p2.(games.Agent))
		}

		pop = new_pop
	}

}

func weighted_selection(items []games.Agent, weights []int) games.Agent {
	// Unimplemented
	return nil
}
