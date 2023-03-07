package airdrop

import "math/big"

type TokensToBlocksChunk struct {
	IntervalFUD    *big.Int // FUD Tokens that a user has deposited in a block interval
	BlockNumber    int64    // The block number that the user had the IntervalFUD deposited
	BlocksInterval int64    // The amount of blocks that the user has the IntervalFUD deposited
}

type AirdropData struct {
	Depositor              string // The account of the depositor
	TokenToBlocks          []TokensToBlocksChunk
	CurrentBlockNumber     int64 // Block number in the block interval (i.e. 1-100)
	LastAirdropBlockNumber int64 // Block number that last airdrop happened
}

type AirdropDataRepository interface {
	CreateChunk(user string, intervalBlockNumber int64, amount *big.Int) (AirdropData, error)
	UpdateCurrentBlockNumber(blockNumber int64) int64
	GetLastAirdropBlockNumber() int64
	CalculateWinTokens(user string, proportion int64, blocksInterval int64, blockNumber int64) (*big.Int, error)
	GetDepositors() []string
	CleanAirdropDataByUser(user string, currentBlockNumber int64, blockInterval int64) error
	UpdateLastAirdropNumber(currentBlockNumber int64) error
}
