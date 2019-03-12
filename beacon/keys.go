package beacon

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"github.com/cloudflare/bn256"
)

// Identity is public key
type Identity struct {
	Key  *bn256.G2
	Addr string
}

// Pair contains public and private key
type Pair struct {
	Private *big.Int
	Public  *Identity
}

// NewKeyPair creates new keyPair
func NewKeyPair(address string) *Pair {
	key, _ := rand.Int(rand.Reader, bn256.Order)
	pubKey := new(bn256.G2).ScalarBaseMult(key)
	pub := &Identity{
		Key:  pubKey,
		Addr: address,
	}
	return &Pair{
		Private: key,
		Public:  pub,
	}
}

// Sign creates a BLS signature S = x * H(m) on a message m using the private key x.
func Sign(x *big.Int, msg []byte) ([]byte, error) {
	h := new(bn256.G1)
	_, err := h.Unmarshal(msg)
	if err != nil {
		print("bls: could not sign")
		return nil, err
	}
	hx := new(bn256.G1).ScalarMult(h, x)
	return hx.Marshal(), nil
}

// Verify checks that e(H(m), X) == e(S, B2)
func Verify(X *Identity, msg []byte, sig []byte) error {
	hx := new(bn256.G1)
	_, err := hx.Unmarshal(sig)
	if err != nil {
		print("bls: bad sig")
		return err
	}
	u := bn256.Pair(hx, new(bn256.G2).ScalarBaseMult(big.NewInt(1)))

	h := new(bn256.G1)
	_, err = h.Unmarshal(msg)
	if err != nil {
		print("bls: bad sig")
		return err
	}
	p := bn256.Pair(h, X.Key)
	if p == u {
		return nil
	}
	return errors.New("bls: bad signature")
}

func printSig(round uint64, sig []byte) {
	hx := new(bn256.G1)
	_, err := hx.Unmarshal(sig)
	if err != nil {
		print("bls: bad sig")
	}
	fmt.Printf("Round : %d\nSignature : %s\n", round, hx.String())
}
