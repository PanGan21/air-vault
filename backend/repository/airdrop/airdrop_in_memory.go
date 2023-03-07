package airdrop

import (
	"fmt"
	"math/big"
	"sync"
)

// airdrop => current block number in interval (1-100)
// user starts with 0
// at block number 1 deposits 10  => TotalFUD = 10 => airvault=(10*nil)
// at block number 40 user deposits 10 => airvault=(10*40)+((10+10)*nil)
// at block number 100 => airvault=(10*(0+40))+((10+10)*(100-40)) => do airdrop
// per user => [(fud tokens) * (deposit/withdraw block of next operation)]
type airdropRepository struct {
	depositors             []string // All the depositors addresses. Store them in an array so we can iterate
	currentBlockNumber     int64    // Block number
	lastAirdropBlockNumber int64    // Block number that last airdrop happened

	mu                         *sync.Mutex
	tokensToBlocksChunksByUser map[string][]TokensToBlocksChunk
}

func NewAirdropRepository() *airdropRepository {
	return &airdropRepository{
		depositors:                 []string{},
		currentBlockNumber:         0,
		lastAirdropBlockNumber:     0, // No airdrop has happened
		mu:                         &sync.Mutex{},
		tokensToBlocksChunksByUser: make(map[string][]TokensToBlocksChunk),
	}
}

func (a *airdropRepository) CreateChunk(user string, blockNumber int64, amount *big.Int) (AirdropData, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	tokensToBlocksChunks, ok := a.tokensToBlocksChunksByUser[user]
	if !ok {
		// That's because users are not persisted in the database
		if !containsString(a.depositors, user) {
			a.depositors = append(a.depositors, user)
		}

		newTokensToBlocksChunk := TokensToBlocksChunk{
			IntervalFUD:    amount,
			BlockNumber:    blockNumber,
			BlocksInterval: 0,
		}
		a.tokensToBlocksChunksByUser[user] = []TokensToBlocksChunk{newTokensToBlocksChunk}
	} else {
		lastTokensToBlocksChunk := tokensToBlocksChunks[len(tokensToBlocksChunks)-1]
		newTokensToBlocksChunk := TokensToBlocksChunk{
			IntervalFUD:    amount,
			BlockNumber:    blockNumber,
			BlocksInterval: 0,
		}
		tokensToBlocksChunks[len(tokensToBlocksChunks)-1].BlocksInterval = blockNumber - lastTokensToBlocksChunk.BlockNumber
		if tokensToBlocksChunks[len(tokensToBlocksChunks)-1].BlockNumber == blockNumber {
			// Override IntervalFUD of last element
			tokensToBlocksChunks[len(tokensToBlocksChunks)-1].IntervalFUD = amount
		} else {
			a.tokensToBlocksChunksByUser[user] = append(tokensToBlocksChunks, newTokensToBlocksChunk)
		}

	}

	airdropData := AirdropData{
		Depositor:              user,
		TokenToBlocks:          a.tokensToBlocksChunksByUser[user],
		CurrentBlockNumber:     a.currentBlockNumber,
		LastAirdropBlockNumber: a.lastAirdropBlockNumber,
	}

	return airdropData, nil
}

func (a *airdropRepository) UpdateCurrentBlockNumber(blockNumber int64) int64 {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.currentBlockNumber = blockNumber

	return a.currentBlockNumber
}

func (a *airdropRepository) GetLastAirdropBlockNumber() int64 {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.lastAirdropBlockNumber
}

// WIN tokens in airdrop = proportion * (# FUD tokens deposited) * (# blocks deposited) / ( total # blocks)
func (a *airdropRepository) CalculateWinTokens(user string, proportion int64, blocksInterval int64, blockNumber int64) (*big.Int, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	tokensToBlocksChunks, ok := a.tokensToBlocksChunksByUser[user]
	if !ok {
		return new(big.Int), fmt.Errorf("user not initialized")
	}

	tokensToBlocksChunks[len(tokensToBlocksChunks)-1].BlocksInterval = blockNumber - tokensToBlocksChunks[len(tokensToBlocksChunks)-1].BlockNumber

	var amount = new(big.Int)
	for _, tokensToBlocksChunk := range tokensToBlocksChunks {
		chunkAmount := new(big.Int).Mul(tokensToBlocksChunk.IntervalFUD, big.NewInt(tokensToBlocksChunk.BlocksInterval))
		amount = new(big.Int).Add(amount, chunkAmount)
	}

	mintProportion := float64(proportion) / 100.0
	blocksRatio := mintProportion / float64(blocksInterval)
	winTokensAmount := new(big.Float).Mul(new(big.Float).SetInt(amount), new(big.Float).SetFloat64(blocksRatio))

	// convert the result to a big.Int and round down
	winTokenAmountInt := new(big.Int)
	winTokensAmount.Int(winTokenAmountInt)

	return winTokenAmountInt, nil
}

func (a *airdropRepository) GetDepositors() []string {
	return a.depositors
}

func (a *airdropRepository) CleanAirdropDataByUser(user string, currentBlockNumber int64, blockInterval int64) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	tokensToBlocksChunks, ok := a.tokensToBlocksChunksByUser[user]
	if !ok {
		return fmt.Errorf("user not initialized")
	}

	lastTokensToBlocksChunks := tokensToBlocksChunks[len(tokensToBlocksChunks)-1]

	delete(a.tokensToBlocksChunksByUser, user)

	remainingChunkNumber := currentBlockNumber - a.lastAirdropBlockNumber
	if remainingChunkNumber > blockInterval {
		remainingChunk := TokensToBlocksChunk{
			IntervalFUD:    lastTokensToBlocksChunks.IntervalFUD,
			BlockNumber:    currentBlockNumber,
			BlocksInterval: currentBlockNumber - blockInterval,
		}
		a.tokensToBlocksChunksByUser[user] = append(a.tokensToBlocksChunksByUser[user], remainingChunk)
	}

	return nil
}

func (a *airdropRepository) UpdateLastAirdropNumber(currentBlockNumber int64) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.lastAirdropBlockNumber = currentBlockNumber
	return nil
}

func containsString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
