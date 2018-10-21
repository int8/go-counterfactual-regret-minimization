package gopoker

import "math"

func to3BinArray(number int) [3]bool {
	return [3]bool{number&1 > 0, number&2 > 0, number&4 > 0}
}

func to4BinArray(number int) [4]bool {
	return [4]bool{number&1 > 0, number&2 > 0, number&4 > 0, number&8 > 0}
}

func cardNameCompare(name1 CardName, name2 CardName) int8 {
	for i := 3; i > -1; i-- {
		if name1[i] != name2[i] {
			if name1[i] {
				return 1 // card 1 is higher
			}
			return -1 // card 2 is higher
		}
	}
	return 0 // cards are equal
}

func CardNameInt(name1 CardName) int8 {
	result := int8(0)
	for i := 0; i < 4; i++ {
		if name1[i] {
			result += int8(math.Pow(2, float64(i)))
		}
	}
	return result
}
