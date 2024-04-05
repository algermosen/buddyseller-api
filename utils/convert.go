package utils

import (
	"math"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
)

func NumericToFloat(n pgtype.Numeric) float64 {
	// Calculate the floating-point value using the integer part and exponent of the numeric value
	intValue := float64(n.Int.Int64())
	exponent := int(n.Exp)

	// Apply the exponent to the integer value to adjust the decimal point
	for i := 0; i < exponent; i++ {
		intValue /= 10
	}

	return float64(intValue)
}

func FloatToNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	n.Exp = 0

	var shift int8 = 0
	var maxIterations = 100

	for {
		if maxIterations == 0 {
			panic("Infinite loop")
		}

		residual := math.Mod(f, 1)
		if residual == f || residual != 0 {
			shift = 1 // Move point to the right 0.45 * 10 -> 4.5
		} else if math.Mod(f, 10) == 0 {
			shift = -1 // Move point to the left 0.45 / 10 -> 0.045
		} else {
			shift = 0 // No need to move point 45
		}

		if shift == 0 {
			n.Int = big.NewInt(int64(f))
			n.Valid = true
			break
		}

		n.Exp -= int32(shift)
		mult := math.Pow10(int(shift))
		f = f * mult
		maxIterations--
	}

	return n
}

func StringToText(s string) pgtype.Text {
	return pgtype.Text{
		Valid:  true,
		String: s,
	}
}

func TextToString(t pgtype.Text) string {
	if t.Valid {
		return t.String
	}
	return ""
}
