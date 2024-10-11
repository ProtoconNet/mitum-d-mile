package digest

import (
	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	utilc "github.com/ProtoconNet/mitum-currency/v3/digest/util"
	"github.com/ProtoconNet/mitum-d-mile/state"
	"github.com/ProtoconNet/mitum-d-mile/types"
	"github.com/ProtoconNet/mitum2/base"
	utilm "github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var maxLimit int64 = 50

var (
	defaultColNameAccount         = "digest_ac"
	defaultColNameContractAccount = "digest_ca"
	defaultColNameBalance         = "digest_bl"
	defaultColNameCurrency        = "digest_cr"
	defaultColNameOperation       = "digest_op"
	defaultColNameDmile           = "digest_dmile_service"
	defaultColNameDmileData       = "digest_dmile_data"
	defaultColNameBlock           = "digest_bm"
)

func DmileDesign(st *cdigest.Database, contract string) (types.Design, base.State, error) {
	filter := utilc.NewBSONFilter("contract", contract)
	q := filter.D()

	opt := options.FindOne().SetSort(
		utilc.NewBSONFilter("height", -1).D(),
	)
	var sta base.State
	if err := st.MongoClient().GetByFilter(
		defaultColNameDmile,
		q,
		func(res *mongo.SingleResult) error {
			i, err := cdigest.LoadState(res.Decode, st.Encoders())
			if err != nil {
				return err
			}
			sta = i
			return nil
		},
		opt,
	); err != nil {
		return types.Design{}, nil, utilm.ErrNotFound.WithMessage(err, "Dmile design by contract account %v", contract)
	}

	if sta != nil {
		de, err := state.GetDesignFromState(sta)
		if err != nil {
			return types.Design{}, nil, err
		}
		return de, sta, nil
	} else {
		return types.Design{}, nil, errors.Errorf("state is nil")
	}
}

func DmileDataByTxID(db *cdigest.Database, contract, key string) (*types.Data, base.State, error) {
	filter := utilc.NewBSONFilter("contract", contract)
	filter = filter.Add("tx_hash", key)
	q := filter.D()

	opt := options.FindOne().SetSort(
		utilc.NewBSONFilter("height", -1).D(),
	)
	var data *types.Data
	var sta base.State
	var err error
	if err := db.MongoClient().GetByFilter(
		defaultColNameDmileData,
		q,
		func(res *mongo.SingleResult) error {
			sta, err = cdigest.LoadState(res.Decode, db.Encoders())
			if err != nil {
				return err
			}
			d, err := state.GetDataFromState(sta)
			if err != nil {
				return err
			}
			data = &d
			return nil
		},
		opt,
	); err != nil {
		return nil, nil, utilm.ErrNotFound.WithMessage(
			err, "Dmile data for txHash %s in contract account %s", key, contract)
	}

	if data != nil {
		return data, sta, nil
	} else {
		return nil, nil, errors.Errorf("data is nil")
	}
}

func DmileDataByMerkleRoot(db *cdigest.Database, contract, key string) (*types.Data, base.State, error) {
	filter := utilc.NewBSONFilter("contract", contract)
	filter = filter.Add("merkle_root", key)
	q := filter.D()

	opt := options.FindOne().SetSort(
		utilc.NewBSONFilter("height", -1).D(),
	)
	var data *types.Data
	var sta base.State
	var err error
	if err := db.MongoClient().GetByFilter(
		defaultColNameDmileData,
		q,
		func(res *mongo.SingleResult) error {
			sta, err = cdigest.LoadState(res.Decode, db.Encoders())
			if err != nil {
				return err
			}
			d, err := state.GetDataFromState(sta)
			if err != nil {
				return err
			}
			data = &d
			return nil
		},
		opt,
	); err != nil {
		return nil, nil, utilm.ErrNotFound.WithMessage(
			err, "Dmile data for merkleRoot %s in contract account %s", key, contract)
	}

	if data != nil {
		return data, sta, nil
	} else {
		return nil, nil, errors.Errorf("data is nil")
	}
}
