package types

import (
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DataJSONMarshaler struct {
	hint.BaseHinter
	MerkleRoot string `json:"merkleRoot"`
	TxID       string `json:"txid"`
}

func (d Data) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DataJSONMarshaler{
		BaseHinter: d.BaseHinter,
		MerkleRoot: d.merkleRoot,
		TxID:       d.txID,
	})
}

type DataJSONUnmarshaler struct {
	Hint       hint.Hint `json:"_hint"`
	MerkleRoot string    `json:"merkleRoot"`
	TxID       string    `json:"txid"`
}

func (d *Data) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of Data")

	var u DataJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	return d.unmarshal(u.Hint, u.MerkleRoot, u.TxID)
}
