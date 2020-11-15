package dec

import (
	"math/big"

	"github.com/shopspring/decimal"
)

var DECIMAL_ZERO = NewDecimalFromInt64(0)

func NewDecimalFromInt64(val int64) decimal.Decimal {
	return decimal.NewFromBigInt(big.NewInt(val), 0)
}
