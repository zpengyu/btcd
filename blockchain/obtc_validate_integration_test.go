package blockchain

import (
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

func loadBlkFileToChain(filename string, chain *BlockChain, skipGenesis bool) (int, int, error) {
	// Extract all blocks from the given file
	var blocks []*btcutil.Block
	blocks, err := loadBlocks(filename)
	if err != nil {
		return -1, -1, err
	}
	fmt.Printf("Load blocks from file: %s, blockTmp size: %+v\n", filename, len(blocks))

	// Add blocks to the chain
	beginIndex := 0
	if skipGenesis {
		beginIndex = 1
	}

	for i := beginIndex; i < len(blocks); i++ {
		_, _, err := chain.ProcessBlock(blocks[i], BFNone)
		if err != nil {
			return -1, -1, err
		}
	}
	fmt.Printf("Chain height: %+v\n", len(chain.bestChain.nodes))
	fmt.Printf("total orphans: %+v\n", len(chain.orphans))

	return len(chain.bestChain.nodes), len(chain.orphans), nil
}

// Download the blk00000.dat from bitcoin official site before testing this
func TestFetchUtxosByHeight(t *testing.T) {
	// Create a new database and chain instance to run tests against.
	// The chain has already contained the genisis block.
	chain, teardownFunc, err := chainSetup("FetchUtxosByHeight", &chaincfg.MainNetParams)
	if err != nil {
		t.Errorf("Failed to setup chain instance: %v", err)
		return
	}
	defer teardownFunc()

	// Since we're not dealing with the real block chain, set the coinbase
	// maturity to 1.
	chain.TstSetCoinbaseMaturity(1)

	testFiles := []string{
		"blk00000.dat",
	}

	for _, filename := range testFiles {
		skipGenesis := false
		if filename == "blk00000.dat" {
			skipGenesis = true
		}
		height, orphanSize, err := loadBlkFileToChain(filename, chain, skipGenesis)
		if err != nil {
			t.Errorf("Error loading blocks from file %s: %v", filename, err)
			break
		}
		fmt.Printf("Load blocks from file %s. chain size: %d orphan size: %d\n", filename, height, orphanSize)
	}

	// File blk00000.dat contains 119878 blocks
	for i := 1; i < 119878; i++ {
		utxos, err := chain.fetchUtxosByHeight(int32(i))
		if err != nil {
			t.Errorf("Cannot fetch utxos by block height, reason is: %v", err)
			return
		}
		if utxos == nil {
			t.Errorf("utxos is nil")
			return
		}
	}
}
