package dmile

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type MigrateDataItemJSONMarshaler struct {
	hint.BaseHinter
	Contract   base.Address             `json:"contract"`
	MerkleRoot string                   `json:"merkleRoot"`
	TxID       string                   `json:"txid"`
	Currency   currencytypes.CurrencyID `json:"currency"`
}

func (it MigrateDataItem) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(MigrateDataItemJSONMarshaler{
		BaseHinter: it.BaseHinter,
		Contract:   it.contract,
		MerkleRoot: it.merkleRoot,
		TxID:       it.txID,
		Currency:   it.currency,
	})
}

type MigrateDataItemJSONUnMarshaler struct {
	Hint       hint.Hint `json:"_hint"`
	Contract   string    `json:"contract"`
	MerkleRoot string    `json:"merkleRoot"`
	TxID       string    `json:"txid"`
	Currency   string    `json:"currency"`
}

func (it *MigrateDataItem) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var uit MigrateDataItemJSONUnMarshaler
	if err := enc.Unmarshal(b, &uit); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *it)
	}

	if err := it.unpack(enc,
		uit.Hint,
		uit.Contract,
		uit.MerkleRoot,
		uit.TxID,
		uit.Currency,
	); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *it)
	}

	return nil
}
