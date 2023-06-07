package utils

import "math/big"

// ToSkh number of SKH to Wei
func ToSkh(skh uint64) *big.Int {
	return new(big.Int).Mul(new(big.Int).SetUint64(skh), big.NewInt(1e18))
}
