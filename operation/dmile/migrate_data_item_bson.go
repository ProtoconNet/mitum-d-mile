package dmile

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (it MigrateDataItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":      it.Hint().String(),
			"contract":   it.contract,
			"merkleRoot": it.merkleRoot,
			"txid":       it.txID,
			"currency":   it.currency,
		},
	)
}

type MigrateDataItemBSONUnmarshaler struct {
	Hint       string `bson:"_hint"`
	Contract   string `bson:"contract"`
	MerkleRoot string `bson:"merkleRoot"`
	TxID       string `bson:"txid"`
	Currency   string `bson:"currency"`
}

func (it *MigrateDataItem) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	var uit MigrateDataItemBSONUnmarshaler
	if err := bson.Unmarshal(b, &uit); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}

	ht, err := hint.ParseHint(uit.Hint)
	if err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}

	if err := it.unpack(enc, ht,
		uit.Contract,
		uit.MerkleRoot,
		uit.TxID,
		uit.Currency,
	); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *it)
	}

	return nil
}
