package main

import (
	"fmt"
	"math"
)

// types{{{
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

func (g *game) remainingCandidates() []chars {
	candidates := product(g.colors, g.slots)
	total := len(candidates)
	for _, round := range g.rounds {
		temp := make([]chars, 0, len(candidates))
		for _, combo := range candidates {
			black, white := calcPips(round.guess, combo)
			if black == round.black && white == round.white {
				temp = append(temp, combo)
			}
		}
		candidates = temp

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

func pipsHash(black, white int) int {
	return black*100 + white
}

func scoreGuess(guess chars, candidates []chars) int {
	pipsCount := map[int]int{}
	max := 0
	for _, code := range candidates {
		b, w := calcPips(code, guess)
		h := pipsHash(b, w)
		if val, ok := pipsCount[h]; ok {
			pipsCount[h] = val + 1
		} else {
			pipsCount[h] = 1
		}
		if pipsCount[h] > max {
			max = pipsCount[h]
		}
	}
	return max
}

func (g *game) nextGuesses() (allOptimal []chars) {
	allCombos := product(g.colors, g.slots)
	candidates := g.remainingCandidates()

	guesses := map[string]int{}
	minScore := 1000
	for _, guess := range allCombos {
		score := scoreGuess(guess, candidates)
		guesses[string(guess)] = score
		if score < minScore {
			minScore = score
		}
	}

	ret := make([]chars, 0, len(guesses))
	for k, v := range guesses {
		if v == minScore {
			ret = append(ret, chars(k))
		}
	}
	return ret
}

func main() {

	var g = game{
		colors: chars("bdglpry"),
		slots:  6,
		rounds: []round{}}

	allOptimal := g.nextGuesses()
	fmt.Printf("\nThere are %d optimal guesses: ", len(allOptimal))
	for _, v := range allOptimal {
		fmt.Printf("%s ", v)
	}
	fmt.Println()
}
