package logic

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

var (
	p          = big.NewInt(123)
	g          = big.NewInt(321)
	privateKey = big.NewInt(456)
	publicKey  = new(big.Int).Exp(g, privateKey, p)
)

func CreateChalleng(r1 int64, user string) (int64, string) {
	authID := fmt.Sprintf("%s_%d", user, time.Now().UnixNano())
	r2, _ := rand.Int(rand.Reader, p)
	// c := new(big.Int).Add(r2, new(big.Int).Mul(big.NewInt(r1), privateKey))
	c := new(big.Int).Add(r2, new(big.Int).Mul(big.NewInt(r1), privateKey))
	// c = c.Mod(c, p)
	return c.Int64(), authID
}

func CreateVerify(s, r2 int64) bool {
	challenge := r2 // сохраненное значение redis
	sVal := new(big.Int).Exp(g, big.NewInt(s), p)
	sInv := new(big.Int).ModInverse(sVal, p)
	result := new(big.Int).Mul(sInv, new(big.Int).Exp(g, big.NewInt(challenge), p))
	result.Mod(result, p)
	if result.Cmp(publicKey) == 0 {
		return true
	}
	return false
}
