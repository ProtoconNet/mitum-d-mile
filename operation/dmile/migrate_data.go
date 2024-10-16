package dmile

import (
	"fmt"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

type DmileItem interface {
	util.Byter
	util.IsValider
	Currency() currencytypes.CurrencyID
}

var MigrateDataItems uint = 100

var (
	MigrateDataFactHint = hint.MustNewHint("mitum-d-mile-migrate-data-operation-fact-v0.0.1")
	MigrateDataHint     = hint.MustNewHint("mitum-d-mile-migrate-data-operation-v0.0.1")
)

type MigrateDataFact struct {
	mitumbase.BaseFact
	sender mitumbase.Address
	items  []MigrateDataItem
}

func NewMigrateDataFact(
	token []byte, sender mitumbase.Address, items []MigrateDataItem) MigrateDataFact {
	bf := mitumbase.NewBaseFact(MigrateDataFactHint, token)
	fact := MigrateDataFact{
		BaseFact: bf,
		sender:   sender,
		items:    items,
	}

	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact MigrateDataFact) IsValid(b []byte) error {
	if n := len(fact.items); n < 1 {
		return common.ErrFactInvalid.Wrap(common.ErrArrayLen.Wrap(errors.Errorf("empty items")))
	} else if n > int(MigrateDataItems) {
		return common.ErrFactInvalid.Wrap(common.ErrArrayLen.Wrap(errors.Errorf("items, %d over max, %d", n, MigrateDataItems)))
	}

	if err := util.CheckIsValiders(nil, false,
		fact.sender,
	); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	founds := map[string]struct{}{}
	for _, it := range fact.items {
		if err := it.IsValid(nil); err != nil {
			return common.ErrFactInvalid.Wrap(err)
		}

		if it.contract.Equal(fact.sender) {
			return common.ErrFactInvalid.Wrap(common.ErrSelfTarget.Wrap(errors.Errorf("sender %v is same with contract account", fact.sender)))
		}

		k := fmt.Sprintf("%s-%s-merkleRoot", it.Contract(), it.MerkleRoot())

		if _, found := founds[k]; found {
			return common.ErrFactInvalid.Wrap(common.ErrDupVal.Wrap(errors.Errorf("merkleRoot %v for contract account %v", it.MerkleRoot(), it.Contract())))
		}

		k = fmt.Sprintf("%s-%s-txhash", it.Contract(), it.TxID())

		if _, found := founds[k]; found {
			return common.ErrFactInvalid.Wrap(common.ErrDupVal.Wrap(errors.Errorf("txhash %v for contract account %v", it.TxID(), it.Contract())))
		}

		founds[k] = struct{}{}
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	return nil
}

func (fact MigrateDataFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact MigrateDataFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact MigrateDataFact) Bytes() []byte {
	is := make([][]byte, len(fact.items))
	for i := range fact.items {
		is[i] = fact.items[i].Bytes()
	}

	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		util.ConcatBytesSlice(is...),
	)
}

func (fact MigrateDataFact) Token() mitumbase.Token {
	return fact.BaseFact.Token()
}

func (fact MigrateDataFact) Sender() mitumbase.Address {
	return fact.sender
}

func (fact MigrateDataFact) Items() []MigrateDataItem {
	return fact.items
}

func (fact MigrateDataFact) Addresses() ([]mitumbase.Address, error) {
	var as []mitumbase.Address

	adrMap := make(map[string]struct{})
	for i := range fact.items {
		for j := range fact.items[i].Addresses() {
			if _, found := adrMap[fact.items[i].Addresses()[j].String()]; !found {
				adrMap[fact.items[i].Addresses()[j].String()] = struct{}{}
				as = append(as, fact.items[i].Addresses()[j])
			}
		}
	}
	as = append(as, fact.sender)

	return as, nil
}

type MigrateData struct {
	common.BaseOperation
}

func NewMigrateData(fact MigrateDataFact) (MigrateData, error) {
	return MigrateData{BaseOperation: common.NewBaseOperation(MigrateDataHint, fact)}, nil
}
