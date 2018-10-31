package acting

func to3BinArray(number int) [3]bool {
	return [3]bool{number&1 > 0, number&2 > 0, number&4 > 0}
}

func CreateByte(x []bool) byte {
	result := byte(0)
	for i, elem := range x {
		if elem {
			result += byte(1) << uint8(i)
		}
	}
	return result
}
