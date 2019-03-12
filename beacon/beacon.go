package beacon

import (
	"crypto/rand"
	"sync"
	"time"

	"github.com/cloudflare/bn256"
)

// Beacon holds info for a specific beacon
type Beacon struct {
	id      int32
	keyPair *Pair
	ticker  *time.Ticker
	round   uint64
	close   chan bool //to stop the beacon
	store   *Store
	sync.Mutex
}

// NewBeacon creates a Beacon
func NewBeacon(address string) *Beacon {
	return &Beacon{
		id:      1000,
		keyPair: NewKeyPair(address),
		close:   make(chan bool),
		round:   0,
		store:   &Store{store: make(map[uint64]int)},
	}
}

// Loop makes a beacon call run every x seconds
func (b *Beacon) Loop(period time.Duration) {
	b.Lock()
	b.ticker = time.NewTicker(period)
	b.Unlock()
	var goToNextRound = true
	var currentRoundFinished bool
	doneCh := make(chan uint64)
	closingCh := make(chan bool)

	for {
		if goToNextRound {
			close(closingCh)
			closingCh = make(chan bool)
			round := b.nextRound()
			go b.run(round, b.keyPair, doneCh, closingCh)
			goToNextRound = false
			currentRoundFinished = false
		}
		select {
		case <-b.close:
			return
		case <-b.ticker.C:
			if !currentRoundFinished {
				close(closingCh)
			}
			goToNextRound = true
			continue
		case roundCh := <-doneCh:
			if roundCh != b.round {
				continue
			}
			currentRoundFinished = true
		}
	}
}

// nextRound increase the round counter
func (b *Beacon) nextRound() uint64 {
	b.Lock()
	b.round++
	b.Unlock()
	return b.round
}

// run creates committements and sign them
func (b *Beacon) run(round uint64, keyPair *Pair, doneCh chan uint64, closingCh chan bool) {
	select {
	case <-closingCh:
		return
	default:
		k, _ := rand.Int(rand.Reader, bn256.Order)
		H := new(bn256.G1).ScalarBaseMult(k)
		msg := H.Marshal()
		sig, err := Sign(b.keyPair.Private, msg)
		if err != nil {
			return
		}
		printSig(round, sig)
		doneCh <- round
	}
}

// Stop stops the beacon
func (b *Beacon) Stop() {
	// TODO: does not work
	b.Lock()
	defer b.Unlock()
	if b.ticker != nil {
		b.ticker.Stop()
	}
	close(b.close)
}
