package generator

var (
	idContinue = func() []rune {
		runes := []rune{
			'_',
			'$',
		}

		for c := 'a'; c <= 'z'; c++ {
			runes = append(runes, c)
		}

		for c := 'A'; c <= 'Z'; c++ {
			runes = append(runes, c)
		}

		for c := '0'; c <= '9'; c++ {
			runes = append(runes, c)
		}

		return runes
	}()
	maxIdContinueIdx = len(idContinue) - 1

	idStart = func() []rune {
		runes := []rune{
			'_',
			'$',
		}

		for c := 'a'; c <= 'z'; c++ {
			runes = append(runes, c)
		}

		for c := 'A'; c <= 'Z'; c++ {
			runes = append(runes, c)
		}

		return runes
	}()
	maxIdStartIdx = len(idStart) - 1
)

type IdGenerator struct {
	bitIdx int
	bits   []int
}

func newIdGenerator() *IdGenerator {
	generator := &IdGenerator{
		bitIdx: 0,
		bits:   []int{0},
	}

	return generator
}

func getListForIdx(idx int) (idList []rune, maxIdListIdx int) {
	if idx == 0 {
		idList = idStart
		maxIdListIdx = maxIdStartIdx
	} else {
		idList = idContinue
		maxIdListIdx = maxIdContinueIdx
	}

	return
}

func (i *IdGenerator) grow() {
	growToLen := len(i.bits) + 1

	i.bits = make([]int, growToLen)
	i.bitIdx = growToLen - 1
}

func (i *IdGenerator) Next() (id string) {
	maxBitIdx := len(i.bits) - 1
	idList, maxIdListIdx := getListForIdx(i.bitIdx)

	readyBits := make([]rune, len(i.bits))
	for index, bit := range i.bits {
		readyBits[index] = idList[bit]
	}
	id = string(readyBits)

	// current bit is maxed
	if i.bits[i.bitIdx] == maxIdListIdx {
		// it's the only bit
		if maxBitIdx == 0 {
			i.grow()
			return
		}

		// now we need to increment first previous bit that isn't maxed
		for j := i.bitIdx; j >= 0; j-- {
			_, maxIdListIdx = getListForIdx(j)

			if i.bits[j] == maxIdListIdx {
				if j == 0 {
					i.grow()
					return
				}

				i.bits[j] = 0
			} else {
				// found it, increment it and move bitIdx to the end
				i.bits[j]++
				i.bitIdx = maxBitIdx
				return
			}
		}
	} else {
		i.bits[i.bitIdx]++
	}

	return
}
