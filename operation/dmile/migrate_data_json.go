package dmile

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-currency/v3/common"

	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type MigrateDataFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Sender base.Address      `json:"sender"`
	Items  []MigrateDataItem `json:"items"`
}

func (fact MigrateDataFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(MigrateDataFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Items:                 fact.items,
	})
}

type MIgrateDataFactJSONUnMarshaler struct {
	base.BaseFactJSONUnmarshaler
	Sender string          `json:"sender"`
	Items  json.RawMessage `json:"items"`
}

func (fact *MigrateDataFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var uf MIgrateDataFactJSONUnMarshaler
	if err := enc.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	fact.BaseFact.SetJSONUnmarshaler(uf.BaseFactJSONUnmarshaler)

	if err := fact.unpack(enc, uf.Sender, uf.Items); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

type MigrateDataMarshaler struct {
	common.BaseOperationJSONMarshaler
}

func (op MigrateData) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(MigrateDataMarshaler{
		BaseOperationJSONMarshaler: op.BaseOperation.JSONMarshaler(),
	})
}

func (op *MigrateData) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperation = ubo

	return nil
}
