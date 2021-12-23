package openapi

import (
	"encoding/json"
	"math/big"
	"strconv"
)

// A Number represents a JSON / YAML number literal.
type Number json.Number

// String returns the literal text of the number.
func (n Number) String() string { return string(n) }

// Float64 returns the number as a float64.
func (n Number) Float64() (float64, error) {
	return strconv.ParseFloat(string(n), 64)
}

// Int64 returns the number as an int64.
func (n Number) Int64() (int64, error) {
	return strconv.ParseInt(string(n), 10, 64)
}

// BigRat returns a *big.Rat representation of n
func (n Number) BigRat() (*big.Rat, bool) {
	return new(big.Rat).SetString(string(n))
}

// BigInt returns a new *big.Int from n
func (n Number) BigInt() (*big.Int, bool) {
	return new(big.Int).SetString(string(n), 10)
}

// BigFloat returns a *big.Float
func (n Number) BigFloat(m big.RoundingMode) (*big.Float, error) {
	f, _, err := big.ParseFloat(string(n), 10, 256, m)
	return f, err
}

// MarshalJSON marshals json
func (n Number) MarshalJSON() ([]byte, error) { return []byte(n), nil }

// UnmarshalJSON unmarshals json
func (n *Number) UnmarshalJSON(data []byte) error {
	var jn json.Number
	err := json.Unmarshal(data, &jn)
	if err != nil {
		return err
	}
	*n = Number(jn)
	return nil
}
