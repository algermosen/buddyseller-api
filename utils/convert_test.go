package utils

import (
	"math/big"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func newNumeric(integer *big.Int, exp int32) pgtype.Numeric {
	numeric := pgtype.Numeric{}
	numeric.Int = integer
	numeric.Exp = exp
	numeric.Valid = true
	return numeric
}

func TestFloatToNumeric(t *testing.T) {
	// Happy path
	var actual, expected pgtype.Numeric
	happyPathValidations := []struct {
		Float      float64
		NumericInt int
		NumericExp int
	}{
		{Float: 0.0, NumericInt: 0, NumericExp: 0},
		{Float: 5.3, NumericInt: 53, NumericExp: -1},
		{Float: 0.045, NumericInt: 45, NumericExp: -3},
		{Float: 38900, NumericInt: 389, NumericExp: 2},
		{Float: 56, NumericInt: 56, NumericExp: 0},
		{Float: 239040, NumericInt: 23904, NumericExp: 1},
	}

	for _, eval := range happyPathValidations {
		actual = FloatToNumeric(eval.Float)
		expected = newNumeric(big.NewInt(int64(eval.NumericInt)), int32(eval.NumericExp))

		if !actual.Valid {
			t.Errorf("Invalid: %v", eval.Float)
		}
		if actual.Int.Cmp(expected.Int) != 0 {
			t.Errorf("NumericInt not equal on %v. \nActual: %v \nExpected: %v", eval.Float, actual.Int, expected.Int)
		}
		if actual.Exp != expected.Exp {
			t.Errorf("NumericExp not equal on %v. \nActual: %v \nExpected: %v", eval.Float, actual.Exp, expected.Exp)
		}
	}

}
