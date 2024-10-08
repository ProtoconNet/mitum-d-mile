package dmile

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (it *MigrateDataItem) unpack(enc encoder.Encoder, ht hint.Hint,
	cAdr, merkleRoot, txid, cid string,
) error {
	it.BaseHinter = hint.NewBaseHinter(ht)

	switch a, err := base.DecodeAddress(cAdr, enc); {
	case err != nil:
		return err
	default:
		it.contract = a
	}

	it.merkleRoot = merkleRoot
	it.txID = txid
	it.currency = currencytypes.CurrencyID(cid)

	return nil
}
