package beacon

import "math/big"

type Store struct {
	store map[uint64]int
}

func (s *Store) Len() int {
	return 0
}

func (s *Store) Last() int {
	return 0
}

func (s *Store) Get(round uint64) (*big.Int, []byte) {
	return nil, nil
}

func (s *Store) Close() {
	return
}
