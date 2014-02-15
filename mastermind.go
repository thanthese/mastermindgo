// This program is designed to answer one question: given the current state of
// a mastermind game (number of colors in pool, length of code, guesses/replies
// made), what is/are the optimal next guess/es?
package main

import (
	"fmt"
	"github.com/thanthese/mastermind/product"
)

type round struct {
	guess []byte
	black int
	white int
}

type game struct {
	colors []byte
	slots  int
	rounds []round
}

// apparently go doesn't get a lib func for this
func indexOf(ls []byte, val byte) int {
	for i, v := range ls {
		if v == val {
			return i
		}
	}
	return -1
}

func calcPips(code, guess []byte) (black, white int) {

	c := make([]byte, len(code))
	copy(c, code)

	g := make([]byte, len(guess))
	copy(g, guess)

	for i := 0; i < len(c); i++ {
		if c[i] == g[i] {
			black++
			c[i] = ' '
			g[i] = ' '
		}
	}

	for i := 0; i < len(c); i++ {
		if c[i] != ' ' {
			if j := indexOf(g, c[i]); j != -1 {
				white++
				c[i] = ' '
				g[j] = ' '
			}
		}
	}
	return
}

// Given the current game state, what codes are still possible? I think the
// itermediate steps are interesting, so they're logged to the stdout.
func (g *game) remainingCandidates() [][]byte {
	candidates := product.Product(g.colors, g.slots)
	total := len(candidates)
	for _, round := range g.rounds {

		plausible := make([][]byte, 0, len(candidates))
		for _, candidate := range candidates {
			black, white := calcPips(round.guess, candidate)
			if black == round.black && white == round.white {
				plausible = append(plausible, candidate)
			}
		}
		candidates = plausible

		// logging
		fmt.Printf("After guess \"%s\" %6d/%d (%.2f%%) combinations remain.\n",
			round.guess, len(candidates), total,
			float64(len(candidates))/float64(total)*100.0)
		if len(candidates) < 10 {
			fmt.Printf("  ")
			for _, c := range candidates {
				fmt.Printf("%s ", c)
			}
			fmt.Println()
		}

	}
	return candidates
}

// I guess I could have used a throwaway struct instead. But that's not what I
// did in javascript or python, and maybe this is faster anyway.
func pipsHash(black, white int) int {
	return black*100 + white
}

// How well does a guess measure up against a pool of plausible candidates?
// Lower score is better.
func scoreGuess(guess []byte, candidates [][]byte) int {
	pipsCount := map[int]int{}
	max := -1
	for _, code := range candidates {
		b, w := calcPips(code, guess)
		h := pipsHash(b, w)
		if val, ok := pipsCount[h]; ok {
			pipsCount[h] = val + 1
		} else {
			pipsCount[h] = 1
		}
		if max == -1 || pipsCount[h] > max {
			max = pipsCount[h]
		}
	}
	return max
}

// Answer the fundamental question of this program: given a game state, what
// are the optimal guesses?
func (g *game) nextGuesses() (allOptimal [][]byte) {

	// sadly, product() is called twice; passing it seems sloppy, though
	allCombos := product.Product(g.colors, g.slots)
	candidates := g.remainingCandidates()

	// score all possible guesses, and find best score
	guesses := map[string]int{}
	minScore := -1
	for _, guess := range allCombos {
		score := scoreGuess(guess, candidates)
		guesses[string(guess)] = score // needs string as key
		if minScore == -1 || score < minScore {
			minScore = score
		}
	}

	// return all guesses that match the best score
	ret := make([][]byte, 0, len(guesses))
	for guess, score := range guesses {
		if score == minScore {
			ret = append(ret, []byte(guess))
		}
	}
	return ret
}

func main() {

	// setup the game we'll analyze
	g := game{
		colors: []byte("bdglpry"),
		slots:  6,
		rounds: []round{
			round{[]byte("bbbddd"), 0, 1},
			round{[]byte("ggglll"), 2, 1},
			round{[]byte("ppprrr"), 0, 2}}}

	allOptimal := g.nextGuesses()
	fmt.Printf("\nThere are %d optimal guesses: ", len(allOptimal))
	for _, guess := range allOptimal {
		fmt.Printf("%s ", guess)
	}
	fmt.Println()
}
