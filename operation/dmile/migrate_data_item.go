package dmile

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	crcytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var MigrateDataItemHint = hint.MustNewHint("mitum-d-mile-migrate-data-item-v0.0.1")

type MigrateDataItem struct {
	hint.BaseHinter
	contract   base.Address
	merkleRoot string
	txID       string
	currency   crcytypes.CurrencyID
}

func NewMigrateDataItem(
	contract base.Address,
	merkleRoot string,
	txID string,
	currency crcytypes.CurrencyID,
) MigrateDataItem {
	return MigrateDataItem{
		BaseHinter: hint.NewBaseHinter(MigrateDataItemHint),
		contract:   contract,
		merkleRoot: merkleRoot,
		txID:       txID,
		currency:   currency,
	}
}

func (it MigrateDataItem) Bytes() []byte {
	return util.ConcatBytesSlice(
		it.contract.Bytes(),
		[]byte(it.merkleRoot),
		[]byte(it.txID),
		it.currency.Bytes(),
	)
}

func (it MigrateDataItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		it.BaseHinter,
		it.contract,
	); err != nil {
		return common.ErrItemInvalid.Wrap(err)
	}

	if !crcytypes.ReValidSpcecialCh.Match([]byte(it.merkleRoot)) {
		return common.ErrValueInvalid.Wrap(errors.Errorf("merkleRoot %v, must match regex `^[^\\s:/?#\\[\\]$@]*$`", it.merkleRoot))
	}

	if len(it.merkleRoot) != 64 {
		return common.ErrValOOR.Wrap(errors.Errorf("merkleRoot length must be %d but %d", LenMerkleRoot, len(it.merkleRoot)))
	}

	return nil
}

func (it MigrateDataItem) Contract() base.Address {
	return it.contract
}

func (it MigrateDataItem) MerkleRoot() string {
	return it.merkleRoot
}

func (it MigrateDataItem) TxID() string {
	return it.txID
}

func (it MigrateDataItem) Currency() crcytypes.CurrencyID {
	return it.currency
}

func (it MigrateDataItem) Addresses() []base.Address {
	ad := make([]base.Address, 1)

	ad[0] = it.contract

	return ad
}
