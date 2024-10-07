package digest

import (
	"github.com/ProtoconNet/mitum-d-mile/state"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/mongo"
)

func (bs *BlockSession) prepareDmile() error {
	if len(bs.sts) < 1 {
		return nil
	}

	var dMileModels []mongo.WriteModel
	var dMileDataModels []mongo.WriteModel
	for i := range bs.sts {
		st := bs.sts[i]
		switch {
		case state.IsDesignStateKey(st.Key()):
			j, err := bs.handleDmileDesignState(st)
			if err != nil {
				return err
			}
			dMileModels = append(dMileModels, j...)
		case state.IsDataStateKey(st.Key()):
			j, err := bs.handleDmileDataState(st)
			if err != nil {
				return err
			}
			dMileDataModels = append(dMileDataModels, j...)
		default:
			continue
		}
	}

	bs.dMileModels = dMileModels
	bs.dMileDataModels = dMileDataModels

	return nil
}

func (bs *BlockSession) handleDmileDesignState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if DmileDesignDoc, err := NewDmileDesignDoc(st, bs.st.Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(DmileDesignDoc),
		}, nil
	}
}

func (bs *BlockSession) handleDmileDataState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if DmileDataDoc, err := NewDmileDataDoc(st, bs.st.Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(DmileDataDoc),
		}, nil
	}
}
