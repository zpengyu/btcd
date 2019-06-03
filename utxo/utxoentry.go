package utxo

// txoFlags is a bitmask defining additional information and state for a
// transaction output in a utxo view.
type TxoFlags uint8

const (
	// TfCoinBase indicates that a txout was contained in a coinbase tx.
	TfCoinBase TxoFlags = 1 << iota

	// TfSpent indicates that a txout is spent.
	TfSpent

	// TfModified indicates that a txout has been modified since it was
	// loaded.
	TfModified

	// TfExpired indicates that a txout has been expired
	TfExpired
)

// UtxoEntry houses details about an individual transaction output in a utxo
// view such as whether or not it was contained in a coinbase tx, the height of
// the block that contains the tx, whether or not it is spent, its public key
// script, and how much it pays.
type UtxoEntry struct {
	// NOTE: Additions, deletions, or modifications to the order of the
	// definitions in this struct should not be changed without considering
	// how it affects alignment on 64-bit platforms.  The current order is
	// specifically crafted to result in minimal padding.  There will be a
	// lot of these in memory, so a few extra bytes of padding adds up.

	Amount      int64
	PkScript    []byte // The public key script for the output.
	BlockHeight int32  // Height of block containing tx.

	// packedFlags contains additional info about output such as whether it
	// is a coinbase, whether it is spent, and whether it has been modified
	// since it was loaded.  This approach is used in order to reduce memory
	// usage since there will be a lot of these in memory.
	PackedFlags TxoFlags
}

// IsModified returns whether or not the output has been modified since it was
// loaded.
func (entry *UtxoEntry) IsModified() bool {
	return entry.PackedFlags&TfModified == TfModified
}

// IsCoinBase returns whether or not the output was contained in a coinbase
// transaction.
func (entry *UtxoEntry) IsCoinBase() bool {
	return entry.PackedFlags&TfCoinBase == TfCoinBase
}

// IsSpent returns whether or not the output has been spent based upon the
// current state of the unspent transaction output view it was obtained from.
func (entry *UtxoEntry) IsSpent() bool {
	return entry.PackedFlags&TfSpent == TfSpent
}

// CheckExpired returns if utxo has expired or not by giving current height.
// Active utxo means that it's existed in the active blockchain.
// Expired utxo means that it's not existed in the active blockchain.
// Active blockchain keeps the latest 368208 blocks.
// 368208 = (7y x 365d x 24h + 2d x 24h) x 6
func (entry *UtxoEntry) CheckExpired(txHeight int32) bool {
	return txHeight-entry.BlockHeight > 368208 // TODO: CHANGE IT
}

// IsExpired returns if utxo has expired from the packedFlags information.
func (entry *UtxoEntry) IsExpired() bool {
	return entry.PackedFlags&TfExpired == TfExpired
}

// Expired marks the output as expired
func (entry *UtxoEntry) Expired() {
	// Mark the output as expired
	entry.PackedFlags |= TfExpired
}

// Spend marks the output as spent.  Spending an output that is already spent
// has no effect.
func (entry *UtxoEntry) Spend() {
	// Nothing to do if the output is already spent.
	if entry.IsSpent() {
		return
	}

	// Mark the output as spent and modified.
	entry.PackedFlags |= TfSpent | TfModified
}

// Clone returns a shallow copy of the utxo entry.
func (entry *UtxoEntry) Clone() *UtxoEntry {
	if entry == nil {
		return nil
	}

	return &UtxoEntry{
		Amount:      entry.Amount,
		PkScript:    entry.PkScript,
		BlockHeight: entry.BlockHeight,
		PackedFlags: entry.PackedFlags,
	}
}
