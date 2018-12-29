package blockchain

import "testing"

func TestCheckExpired(t *testing.T) {
	utxoEntry := &UtxoEntry{
		amount:      1000000,
		blockHeight: 100001,
	}

	if utxoEntry.CheckExpired(200000) {
		t.Error("Entry should not expired")
	}
}

func TestIsExpired(t *testing.T) {
	utxoEntryExpired := &UtxoEntry{
		amount:      1000000,
		blockHeight: 10000,
		packedFlags: tfExpired,
	}
	if !utxoEntryExpired.IsExpired() {
		t.Error("Entry should expired")
	}

	utxoEntryActive := &UtxoEntry{
		amount:      1000000,
		blockHeight: 10000,
	}
	if utxoEntryActive.IsExpired() {
		t.Error("Entry should not expired")
	}
}

func TestExpired(t *testing.T) {
	utxoEntryActive := &UtxoEntry{
		amount:      1000000,
		blockHeight: 10000,
	}
	if utxoEntryActive.IsExpired() {
		t.Error("Entry should not expired")
	}
	utxoEntryActive.Expired()
	if !utxoEntryActive.IsExpired() {
		t.Error("Entry should expired after set tfExpired flag")
	}
}
