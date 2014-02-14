package product

import "math"

// Return all possible combinations, given colors list and length of string.
func Product(colors []byte, length int) [][]byte {

	// how many of these are we going to make?
	comboCount := int(math.Pow(float64(len(colors)), float64(length)))

	// initialize slice to return
	ret := make([][]byte, comboCount)
	for i := 0; i < comboCount; i++ {
		ret[i] = make([]byte, length)
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
