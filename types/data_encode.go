package types

import (
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (d *Data) unmarshal(
	ht hint.Hint,
	merkleRoot, txID string,
) error {
	d.BaseHinter = hint.NewBaseHinter(ht)
	d.merkleRoot = merkleRoot
	d.txID = txID

	return nil
}
