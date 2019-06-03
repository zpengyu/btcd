package blockchain

import (
	"testing"

	"github.com/btcsuite/btcd/utxo"
)

func TestCheckExpired(t *testing.T) {
	utxoEntry := &utxo.UtxoEntry{
		Amount:      1000000,
		BlockHeight: 100001,
	}

	if utxoEntry.CheckExpired(200000) {
		t.Error("Entry should not expired")
	}
}

func TestIsExpired(t *testing.T) {
	utxoEntryExpired := &utxo.UtxoEntry{
		Amount:      1000000,
		BlockHeight: 10000,
		PackedFlags: utxo.TfExpired,
	}
	if !utxoEntryExpired.IsExpired() {
		t.Error("Entry should expired")
	}

	utxoEntryActive := &utxo.UtxoEntry{
		Amount:      1000000,
		BlockHeight: 10000,
	}
	if utxoEntryActive.IsExpired() {
		t.Error("Entry should not expired")
	}
}

func TestExpired(t *testing.T) {
	utxoEntryActive := &utxo.UtxoEntry{
		Amount:      1000000,
		BlockHeight: 10000,
	}
	if utxoEntryActive.IsExpired() {
		t.Error("Entry should not expired")
	}
	utxoEntryActive.Expired()
	if !utxoEntryActive.IsExpired() {
		t.Error("Entry should expired after set tfExpired flag")
	}
}
