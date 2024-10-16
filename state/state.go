package state

import (
	"fmt"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-d-mile/types"
	"strings"

	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	DesignStateValueHint = hint.MustNewHint("mitum-d-mile-design-state-value-v0.0.1")
	DmileStateKeyPrefix  = "dmile"
	DesignStateKeySuffix = "design"
)

func DmileStateKey(addr mitumbase.Address) string {
	return fmt.Sprintf("%s:%s", DmileStateKeyPrefix, addr.String())
}

type DesignStateValue struct {
	hint.BaseHinter
	Design types.Design
}

func NewDesignStateValue(design types.Design) DesignStateValue {
	return DesignStateValue{
		BaseHinter: hint.NewBaseHinter(DesignStateValueHint),
		Design:     design,
	}
}

func (sv DesignStateValue) Hint() hint.Hint {
	return sv.BaseHinter.Hint()
}

func (sv DesignStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid DesignStateValue")

	if err := sv.BaseHinter.IsValid(DesignStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := sv.Design.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (sv DesignStateValue) HashBytes() []byte {
	return sv.Design.Bytes()
}

func GetDesignFromState(st mitumbase.State) (types.Design, error) {
	v := st.Value()
	if v == nil {
		return types.Design{}, errors.Errorf("state value is nil")
	}

	d, ok := v.(DesignStateValue)
	if !ok {
		return types.Design{}, errors.Errorf("expected DesignStateValue but %T", v)
	}

	return d.Design, nil
}

func IsDesignStateKey(key string) bool {
	return strings.HasPrefix(key, DmileStateKeyPrefix) && strings.HasSuffix(key, DesignStateKeySuffix)
}

func DesignStateKey(addr mitumbase.Address) string {
	return fmt.Sprintf("%s:%s", DmileStateKey(addr), DesignStateKeySuffix)
}

var (
	DataStateValueHint     = hint.MustNewHint("mitum-d-mile-data-state-value-v0.0.1")
	DataStateKeySuffix     = "data"
	DataStateTxIDKeySuffix = "txid"
)

type DataStateValue struct {
	hint.BaseHinter
	Data types.Data
}

func NewDataStateValue(data types.Data) DataStateValue {
	return DataStateValue{
		BaseHinter: hint.NewBaseHinter(DataStateValueHint),
		Data:       data,
	}
}

func (sv DataStateValue) Hint() hint.Hint {
	return sv.BaseHinter.Hint()
}

func (sv DataStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid DataStateValue")

	if err := sv.BaseHinter.IsValid(DataStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := sv.Data.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (sv DataStateValue) HashBytes() []byte {
	return sv.Data.Bytes()
}

func GetDataFromState(st mitumbase.State) (types.Data, error) {
	v := st.Value()
	if v == nil {
		return types.Data{}, errors.Errorf("State value is nil")
	}

	ts, ok := v.(DataStateValue)
	if !ok {
		return types.Data{}, common.ErrTypeMismatch.Wrap(errors.Errorf("expected DataStateValue found, %T", v))
	}

	return ts.Data, nil
}

func IsDataStateKey(key string) bool {
	return strings.HasPrefix(key, DmileStateKeyPrefix) && strings.HasSuffix(key, DataStateKeySuffix)
}

func DataStateMerkleRootKey(addr mitumbase.Address, key string) string {
	return fmt.Sprintf("%s:%s:%s", DmileStateKey(addr), key, DataStateKeySuffix)
}

func DataStateTxIDKey(addr mitumbase.Address, key string) string {
	return fmt.Sprintf("%s:%s:%s", DmileStateKey(addr), key, DataStateTxIDKeySuffix)
}
