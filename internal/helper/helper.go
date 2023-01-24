package helper

import "fmt"

func FormatInt(number int) string {
	if number >= 1e6 {
		return fmt.Sprintf("%.1fm", float64(number)/1e6)
	} else if number >= 1e3 {
		return fmt.Sprintf("%.1fk", float64(number)/1e3)
	}

	return fmt.Sprintf("%d", number)
}
