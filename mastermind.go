// This program is designed to answer one question: given the current state of
// a mastermind game (number of colors in pool, length of code, guesses/replies
// made), what is/are the optimal next guess/es?
package main

import (
	"fmt"
	"math"
)

// types{{{

// I'm representing sequences of colors -- "bbdd" -- as slices of bytes.
// Apparently you're supposed to use strings for this, even to represent chars.
// Well, I got neck-deep in before I learned that, and I don't want to fix it
// now. Besides, I don't think it matters in this case because I know I'll only
// ever be using letters from the alphabet.
type chars []byte

type round struct {
	guess chars
	black int
	white int
}

type game struct {
	colors chars
	slots  int
	rounds []round
}

// }}}
// product{{{

// all possible combinations, given colors list and length of string
func product(colors chars, length int) []chars {

	// how many of these are we going to make?
	comboCount := int(math.Pow(float64(len(colors)), float64(length)))

	// initialize slice to return
	ret := make([]chars, comboCount)
	for i := 0; i < comboCount; i++ {
		ret[i] = make(chars, length)
	}

	incEvery := 1
	comboIndex := 0 // index of invisible for-loop

	// build up answer place by place
	for place := length - 1; place >= 0; place-- {
	NextPlace:
		// repeat through all combos
		for {
			// for each char...
			for _, char := range colors {
				// ...repeat it incEvery times
				for r := 0; r < incEvery; r++ {

					ret[comboIndex][place] = char

					comboIndex++
					if comboIndex == comboCount {
						break NextPlace
					}
				}
			}
		}

		comboIndex = 0
		incEvery *= len(colors)
	}
	return ret
}

// }}}

// apparently go doesn't get a lib func for this
func indexOf(ls chars, val byte) int {
	for i, v := range ls {
		if v == val {
			return i
		}
	}
	return -1
}

func calcPips(code, guess chars) (black, white int) {

	c := make(chars, len(code))
	copy(c, code)

	g := make(chars, len(guess))
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
func (g *game) remainingCandidates() []chars {
	candidates := product(g.colors, g.slots)
	total := len(candidates)
	for _, round := range g.rounds {

		plausible := make([]chars, 0, len(candidates))
		for _, candidate := range candidates {
			black, white := calcPips(round.guess, candidate)
			if black == round.black && white == round.white {
				plausible = append(plausible, candidate)
			}
		}
		candidates = plausible

		// logging
		fmt.Printf("After guess \"%s\" %d/%d combinations remain.\n",
			round.guess, len(candidates), total)
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
func scoreGuess(guess chars, candidates []chars) int {
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
func (g *game) nextGuesses() (allOptimal []chars) {

	// sadly, product() is called twice; passing it seems sloppy, though
	allCombos := product(g.colors, g.slots)
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
	ret := make([]chars, 0, len(guesses))
	for guess, score := range guesses {
		if score == minScore {
			ret = append(ret, chars(guess))
		}
	}
	return ret
}

func main() {

	// setup the game we'll analyze
	g := game{
		colors: chars("bdglpry"),
		slots:  6,
		rounds: []round{
			round{chars("bbbddd"), 0, 1},
			round{chars("ggglll"), 2, 1},
			round{chars("ppprrr"), 0, 2}}}

	allOptimal := g.nextGuesses()
	fmt.Printf("\nThere are %d optimal guesses: ", len(allOptimal))
	for _, guess := range allOptimal {
		fmt.Printf("%s ", guess)
	}
	fmt.Println()
}
