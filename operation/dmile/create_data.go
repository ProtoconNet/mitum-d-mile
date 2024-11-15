package dmile

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var LenMerkleRoot = 64

var (
	CreateDataFactHint = hint.MustNewHint("mitum-d-mile-create-data-operation-fact-v0.0.1")
	CreateDataHint     = hint.MustNewHint("mitum-d-mile-create-data-operation-v0.0.1")
)

type CreateDataFact struct {
	mitumbase.BaseFact
	sender     mitumbase.Address
	contract   mitumbase.Address
	merkleRoot string
	currency   currencytypes.CurrencyID
}

func NewCreateDataFact(
	token []byte, sender, contract mitumbase.Address,
	merkleRoot string, currency currencytypes.CurrencyID) CreateDataFact {
	bf := mitumbase.NewBaseFact(CreateDataFactHint, token)
	fact := CreateDataFact{
		BaseFact:   bf,
		sender:     sender,
		contract:   contract,
		merkleRoot: merkleRoot,
		currency:   currency,
	}

	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact CreateDataFact) IsValid(b []byte) error {
	if fact.sender.Equal(fact.contract) {
		return common.ErrFactInvalid.Wrap(
			common.ErrSelfTarget.Wrap(errors.Errorf("sender %v is same with contract account", fact.sender)))
	}

	if !currencytypes.ReValidSpcecialCh.Match([]byte(fact.merkleRoot)) {
		return common.ErrValueInvalid.Wrap(errors.Errorf("merkleRoot %v, must match regex `^[^\\s:/?#\\[\\]$@]*$`", fact.merkleRoot))
	}

	if len(fact.merkleRoot) != 64 {
		return common.ErrValOOR.Wrap(errors.Errorf("merkleRoot length must be %d but %d", LenMerkleRoot, len(fact.merkleRoot)))
	}

	if err := util.CheckIsValiders(nil, false,
		fact.BaseHinter,
		fact.sender,
		fact.contract,
		fact.currency,
	); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	return nil
}

func (fact CreateDataFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact CreateDataFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact CreateDataFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		[]byte(fact.merkleRoot),
		fact.currency.Bytes(),
	)
}

func (fact CreateDataFact) Token() mitumbase.Token {
	return fact.BaseFact.Token()
}

func (fact CreateDataFact) Sender() mitumbase.Address {
	return fact.sender
}

func (fact CreateDataFact) Contract() mitumbase.Address {
	return fact.contract
}

func (fact CreateDataFact) MerkleRoot() string {
	return fact.merkleRoot
}

func (fact CreateDataFact) Currency() currencytypes.CurrencyID {
	return fact.currency
}

func (fact CreateDataFact) Addresses() ([]mitumbase.Address, error) {
	as := []mitumbase.Address{fact.sender}

	return as, nil
}

type CreateData struct {
	common.BaseOperation
}

func NewCreateData(fact CreateDataFact) (CreateData, error) {
	return CreateData{BaseOperation: common.NewBaseOperation(CreateDataHint, fact)}, nil
}
