package types

import (
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

const MinKeyLen = 128
const Nid = "0000"

var DataHint = hint.MustNewHint("mitum-d-mile-data-v0.0.1")
var DocumentHint = hint.MustNewHint("mitum-d-mile-document-v0.0.1")

type Data struct {
	hint.BaseHinter
	merkleRoot string
	txID       string
}

func NewData(
	merkleRoot, txID string,
) Data {
	data := Data{
		BaseHinter: hint.NewBaseHinter(DataHint),
		merkleRoot: merkleRoot,
		txID:       txID,
	}
	return data
}

func (d Data) IsValid([]byte) error {
	return nil
}

func (d Data) Bytes() []byte {
	return util.ConcatBytesSlice(
		[]byte(d.merkleRoot),
		[]byte(d.txID),
	)
}

func (d Data) MerkleRoot() string {
	return d.merkleRoot
}

func (d Data) TxID() string {
	return d.txID
}

func (d Data) Equal(ct Data) bool {
	if d.merkleRoot != ct.merkleRoot {
		return false
	}
	if d.txID != ct.txID {
		return false
	}

	return true
}
