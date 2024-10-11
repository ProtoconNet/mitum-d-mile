package digest

import (
	mongodb "github.com/ProtoconNet/mitum-currency/v3/digest/mongodb"
	bsonutil "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	cstate "github.com/ProtoconNet/mitum-currency/v3/state"
	"github.com/ProtoconNet/mitum-d-mile/state"
	"github.com/ProtoconNet/mitum-d-mile/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type DmileDesignDoc struct {
	mongodb.BaseDoc
	st     base.State
	design types.Design
}

// NewDmileDesignDoc get the State of Dmile Design
func NewDmileDesignDoc(st base.State, enc encoder.Encoder) (DmileDesignDoc, error) {
	design, err := state.GetDesignFromState(st)

	if err != nil {
		return DmileDesignDoc{}, err
	}

	b, err := mongodb.NewBaseDoc(nil, st, enc)
	if err != nil {
		return DmileDesignDoc{}, err
	}

	return DmileDesignDoc{
		BaseDoc: b,
		st:      st,
		design:  design,
	}, nil
}

func (doc DmileDesignDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := cstate.ParseStateKey(doc.st.Key(), state.DmileStateKeyPrefix, 3)

	m["contract"] = parsedKey[1]
	m["height"] = doc.st.Height()

	return bsonutil.Marshal(m)
}

type DmileDataDoc struct {
	mongodb.BaseDoc
	st   base.State
	data types.Data
}

func NewDmileDataDoc(st base.State, enc encoder.Encoder) (DmileDataDoc, error) {
	data, err := state.GetDataFromState(st)
	if err != nil {
		return DmileDataDoc{}, err
	}

	b, err := mongodb.NewBaseDoc(nil, st, enc)
	if err != nil {
		return DmileDataDoc{}, err
	}

	return DmileDataDoc{
		BaseDoc: b,
		st:      st,
		data:    data,
	}, nil
}

func (doc DmileDataDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := cstate.ParseStateKey(doc.st.Key(), state.DmileStateKeyPrefix, 4)
	if err != nil {
		return nil, err
	}

	m["contract"] = parsedKey[1]
	m["merkle_root"] = doc.data.MerkleRoot()
	m["tx_hash"] = doc.data.TxID()
	m["height"] = doc.st.Height()

	return bsonutil.Marshal(m)
}
