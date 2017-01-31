package getall1

func GetAll1(number int64) int64 {
	var count int64 = 0
	if number == 0 {
		return 0
	}
	var high, cur, low int64 = 0, 0, 0
	var base int64 = 1
	for ; number/base > 0; base *= 10 {
		high = number / 10 / base
		cur = (number / base) % 10
		low = number - number/base*base
		switch cur {
		case 0:
			count += high * base
		case 1:
			count += high*base + low + 1
		default:
			count += (high + 1) * base
		}
	}
	return count
}
