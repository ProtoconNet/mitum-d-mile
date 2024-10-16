package dmile

import (
	"context"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/state"
	crtypes "github.com/ProtoconNet/mitum-currency/v3/types"
	dmilestate "github.com/ProtoconNet/mitum-d-mile/state"
	"github.com/ProtoconNet/mitum-d-mile/types"
	"github.com/pkg/errors"
	"sync"

	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	stateextension "github.com/ProtoconNet/mitum-currency/v3/state/extension"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

var createDataProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(CreateDataProcessor)
	},
}

func (CreateData) Process(
	_ context.Context, _ mitumbase.GetStateFunc,
) ([]mitumbase.StateMergeValue, mitumbase.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type CreateDataProcessor struct {
	*mitumbase.BaseOperationProcessor
}

func NewCreateDataProcessor() crtypes.GetNewProcessor {
	return func(
		height mitumbase.Height,
		getStateFunc mitumbase.GetStateFunc,
		newPreProcessConstraintFunc mitumbase.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc mitumbase.NewOperationProcessorProcessFunc,
	) (mitumbase.OperationProcessor, error) {
		e := util.StringError("failed to create new CreateDataProcessor")

		nopp := createDataProcessorPool.Get()
		opp, ok := nopp.(*CreateDataProcessor)
		if !ok {
			return nil, e.Errorf("expected %T, not %T", CreateDataProcessor{}, nopp)
		}

		b, err := mitumbase.NewBaseOperationProcessor(
			height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
		if err != nil {
			return nil, e.Wrap(err)
		}

		opp.BaseOperationProcessor = b

		return opp, nil
	}
}

func (opp *CreateDataProcessor) PreProcess(
	ctx context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc,
) (context.Context, mitumbase.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(CreateDataFact)
	if !ok {
		return ctx, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", CreateDataFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	if err := state.CheckExistsState(statecurrency.DesignStateKey(fact.Currency()), getStateFunc); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCurrencyNF).Errorf("currency id %v", fact.Currency())), nil
	}

	if _, _, aErr, cErr := state.ExistsCAccount(fact.Sender(), "sender", true, false, getStateFunc); aErr != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", aErr)), nil
	} else if cErr != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.Wrap(common.ErrMCAccountNA).
				Errorf("%v", cErr)), nil
	}

	if err := state.CheckFactSignsByState(fact.Sender(), op.Signs(), getStateFunc); err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMSignInvalid).
				Errorf("%v", err)), nil
	}

	_, cSt, aErr, cErr := state.ExistsCAccount(fact.Contract(), "contract", true, true, getStateFunc)
	if aErr != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", aErr)), nil
	} else if cErr != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", cErr)), nil
	}

	_, err := stateextension.CheckCAAuthFromState(cSt, fact.Sender())
	if err != nil {
		return ctx, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	if err := state.CheckExistsState(dmilestate.DesignStateKey(fact.Contract()), getStateFunc); err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMServiceNF).Errorf("d-mile service in contract account %v",
				fact.Contract(),
			)), nil
	}

	if found, _ := state.CheckNotExistsState(dmilestate.DataStateMerkleRootKey(fact.Contract(), fact.MerkleRoot()), getStateFunc); found {
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMStateE).Errorf("d-mile data with merkleRoot %v for contract account %v",
				fact.MerkleRoot(), fact.Contract(),
			)), nil
	}

	if found, _ := state.CheckNotExistsState(dmilestate.DataStateTxIDKey(fact.Contract(), fact.Hash().String()), getStateFunc); found {
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMStateE).Errorf("d-mileage data with txHash %v for contract account %v",
				fact.Hash().String(), fact.Contract(),
			)), nil
	}

	return ctx, nil, nil
}

func (opp *CreateDataProcessor) Process( // nolint:dupl
	_ context.Context, op mitumbase.Operation, getStateFunc mitumbase.GetStateFunc) (
	[]mitumbase.StateMergeValue, mitumbase.OperationProcessReasonError, error,
) {
	fact, _ := op.Fact().(CreateDataFact)

	data := types.NewData(
		fact.MerkleRoot(), fact.Hash().String(),
	)
	if err := data.IsValid(nil); err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("invalid d-mile data; %w", err), nil
	}

	var sts []mitumbase.StateMergeValue // nolint:prealloc
	sts = append(sts, state.NewStateMergeValue(
		dmilestate.DataStateMerkleRootKey(fact.Contract(), fact.MerkleRoot()),
		dmilestate.NewDataStateValue(data),
	))
	sts = append(sts, state.NewStateMergeValue(
		dmilestate.DataStateTxIDKey(fact.Contract(), fact.Hash().String()),
		dmilestate.NewDataStateValue(data),
	))

	currencyPolicy, err := state.ExistsCurrencyPolicy(fact.Currency(), getStateFunc)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError("currency not found, %v; %w", fact.Currency(), err), nil
	}

	if currencyPolicy.Feeer().Receiver() == nil {
		return sts, nil, nil
	}

	fee, err := currencyPolicy.Feeer().Fee(common.ZeroBig)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			"failed to check fee of currency, %q; %w",
			fact.Currency(),
			err,
		), nil
	}

	senderBalSt, err := state.ExistsState(
		statecurrency.BalanceStateKey(fact.Sender(), fact.Currency()),
		"sender balance",
		getStateFunc,
	)
	if err != nil {
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			"sender %v balance not found; %w",
			fact.Sender(),
			err,
		), nil
	}

	switch senderBal, err := statecurrency.StateBalanceValue(senderBalSt); {
	case err != nil:
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			"failed to get balance value, %q; %w",
			statecurrency.BalanceStateKey(fact.Sender(), fact.Currency()),
			err,
		), nil
	case senderBal.Big().Compare(fee) < 0:
		return nil, mitumbase.NewBaseOperationProcessReasonError(
			"not enough balance of sender, %q",
			fact.Sender(),
		), nil
	}

	v, ok := senderBalSt.Value().(statecurrency.BalanceStateValue)
	if !ok {
		return nil, mitumbase.NewBaseOperationProcessReasonError("expected BalanceStateValue, not %T", senderBalSt.Value()), nil
	}

	if err := state.CheckExistsState(statecurrency.AccountStateKey(currencyPolicy.Feeer().Receiver()), getStateFunc); err != nil {
		return nil, nil, err
	} else if feeRcvrSt, found, err := getStateFunc(statecurrency.BalanceStateKey(currencyPolicy.Feeer().Receiver(), fact.currency)); err != nil {
		return nil, nil, err
	} else if !found {
		return nil, nil, errors.Errorf("feeer receiver %s not found", currencyPolicy.Feeer().Receiver())
	} else if feeRcvrSt.Key() != senderBalSt.Key() {
		r, ok := feeRcvrSt.Value().(statecurrency.BalanceStateValue)
		if !ok {
			return nil, nil, errors.Errorf("expected %T, not %T", statecurrency.BalanceStateValue{}, feeRcvrSt.Value())
		}
		sts = append(sts, common.NewBaseStateMergeValue(
			feeRcvrSt.Key(),
			statecurrency.NewAddBalanceStateValue(r.Amount.WithBig(fee)),
			func(height mitumbase.Height, st mitumbase.State) mitumbase.StateValueMerger {
				return statecurrency.NewBalanceStateValueMerger(height, feeRcvrSt.Key(), fact.currency, st)
			},
		))

		sts = append(sts, common.NewBaseStateMergeValue(
			senderBalSt.Key(),
			statecurrency.NewDeductBalanceStateValue(v.Amount.WithBig(fee)),
			func(height mitumbase.Height, st mitumbase.State) mitumbase.StateValueMerger {
				return statecurrency.NewBalanceStateValueMerger(height, senderBalSt.Key(), fact.currency, st)
			},
		))
	}

	return sts, nil, nil
}

func (opp *CreateDataProcessor) Close() error {
	createDataProcessorPool.Put(opp)

	return nil
}
