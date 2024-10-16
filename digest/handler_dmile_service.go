package digest

import (
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-d-mile/types"
	"github.com/ProtoconNet/mitum2/base"
	"net/http"
	"strings"
	"time"
)

func (hd *Handlers) handleDmileDesign(w http.ResponseWriter, r *http.Request) {
	cacheKey := currencydigest.CacheKeyPath(r)
	if err := currencydigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := currencydigest.ParseRequest(w, r, "contract")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleDmileDesignInGroup(contract)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)

		if !shared {
			currencydigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleDmileDesignInGroup(contract string) ([]byte, error) {
	var de types.Design
	var st base.State

	de, st, err := DmileDesign(hd.database, contract)
	if err != nil {
		return nil, err
	}

	i, err := hd.buildDmileDesign(contract, de, st)
	if err != nil {
		return nil, err
	}
	return hd.encoder.Marshal(i)
}

func (hd *Handlers) buildDmileDesign(contract string, de types.Design, st base.State) (currencydigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathDmileDesign, "contract", contract)
	if err != nil {
		return nil, err
	}

	var hal currencydigest.Hal
	hal = currencydigest.NewBaseHal(de, currencydigest.NewHalLink(h, nil))

	h, err = hd.combineURL(currencydigest.HandlerPathBlockByHeight, "height", st.Height().String())
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("block", currencydigest.NewHalLink(h, nil))

	for i := range st.Operations() {
		h, err := hd.combineURL(currencydigest.HandlerPathOperation, "hash", st.Operations()[i].String())
		if err != nil {
			return nil, err
		}
		hal = hal.AddLink("operations", currencydigest.NewHalLink(h, nil))
	}

	return hal, nil
}

func (hd *Handlers) handleDmileDataByTxID(w http.ResponseWriter, r *http.Request) {
	cacheKey := currencydigest.CacheKeyPath(r)
	if err := currencydigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := currencydigest.ParseRequest(w, r, "contract")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	txHash, err, status := currencydigest.ParseRequest(w, r, "tx_hash")
	key := strings.TrimPrefix(txHash, "0x")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleDmileDataByTxIDInGroup(contract, key)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)

		if !shared {
			currencydigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleDmileDataByTxIDInGroup(contract, key string) ([]byte, error) {
	data, st, err := DmileDataByTxID(hd.database, contract, key)
	if err != nil {
		return nil, err
	}

	i, err := hd.buildDmileDataByTxIDHal(contract, *data, st)
	if err != nil {
		return nil, err
	}
	return hd.encoder.Marshal(i)
}

func (hd *Handlers) buildDmileDataByTxIDHal(
	contract string, data types.Data, st base.State) (currencydigest.Hal, error) {
	h, err := hd.combineURL(
		HandlerPathDmileDataByTxID,
		"contract", contract, "tx_hash", data.TxID())
	if err != nil {
		return nil, err
	}

	var hal currencydigest.Hal
	nTxID := "0x" + data.TxID()
	nData := types.NewData(data.MerkleRoot(), nTxID)
	hal = currencydigest.NewBaseHal(nData, currencydigest.NewHalLink(h, nil))
	h, err = hd.combineURL(currencydigest.HandlerPathBlockByHeight, "height", st.Height().String())
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("block", currencydigest.NewHalLink(h, nil))

	for i := range st.Operations() {
		h, err := hd.combineURL(currencydigest.HandlerPathOperation, "hash", st.Operations()[i].String())
		if err != nil {
			return nil, err
		}
		hal = hal.AddLink("operations", currencydigest.NewHalLink(h, nil))
	}

	return hal, nil
}

func (hd *Handlers) handleDmileDataByMerkleRoot(w http.ResponseWriter, r *http.Request) {
	cacheKey := currencydigest.CacheKeyPath(r)
	if err := currencydigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := currencydigest.ParseRequest(w, r, "contract")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	key, err, status := currencydigest.ParseRequest(w, r, "merkle_root")
	if err != nil {
		currencydigest.HTTP2ProblemWithError(w, err, status)
		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleDmileDataByMerkleRootInGroup(contract, key)
	}); err != nil {
		currencydigest.HTTP2HandleError(w, err)
	} else {
		currencydigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)

		if !shared {
			currencydigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleDmileDataByMerkleRootInGroup(contract, key string) ([]byte, error) {
	data, st, err := DmileDataByMerkleRoot(hd.database, contract, key)
	if err != nil {
		return nil, err
	}

	i, err := hd.buildDmileDataByMerkleRootHal(contract, *data, st)
	if err != nil {
		return nil, err
	}
	return hd.encoder.Marshal(i)
}

func (hd *Handlers) buildDmileDataByMerkleRootHal(
	contract string, data types.Data, st base.State) (currencydigest.Hal, error) {
	h, err := hd.combineURL(
		HandlerPathDmileDataByMerkleRoot,
		"contract", contract, "merkle_root", data.MerkleRoot())
	if err != nil {
		return nil, err
	}

	var hal currencydigest.Hal
	nTxID := "0x" + data.TxID()
	nData := types.NewData(data.MerkleRoot(), nTxID)
	hal = currencydigest.NewBaseHal(nData, currencydigest.NewHalLink(h, nil))
	h, err = hd.combineURL(currencydigest.HandlerPathBlockByHeight, "height", st.Height().String())
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("block", currencydigest.NewHalLink(h, nil))

	for i := range st.Operations() {
		h, err := hd.combineURL(currencydigest.HandlerPathOperation, "hash", st.Operations()[i].String())
		if err != nil {
			return nil, err
		}
		hal = hal.AddLink("operations", currencydigest.NewHalLink(h, nil))
	}

	return hal, nil
}
