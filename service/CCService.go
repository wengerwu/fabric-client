package service

import (
	"fabric-client/sdkInit"
	"time"

	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

var ClientMap map[string]*sdkInit.Client

type Setup struct {
	ChaincodeID string
	Client      *channel.Client
	LClient     *ledger.Client
}

func (setup *Setup) Execute(fcn string, args [][]byte) (channel.Response, error) {
	request := channel.Request{
		ChaincodeID: setup.ChaincodeID,
		Fcn:         fcn,
		Args:        args,
	}
	response, err := setup.Client.Execute(request)
	return response, err
}

func (setup *Setup) QueryBlockByTxID(txID fab.TransactionID)  (*common.Block, error) {
	block, err :=setup.LClient.QueryBlockByTxID(txID)
	return block, err
}

func (setup *Setup) SetEvent(eventFilter string, eventCallbackUrl string, handler func(eventFilter string, callbackUrl string, event *fab.CCEvent)) error {
	reg, notifier, err := setup.Client.RegisterChaincodeEvent(setup.ChaincodeID, eventFilter)
	if err != nil {
		return err
	}
	defer setup.Client.UnregisterChaincodeEvent(reg)

	select {
	case ccEvent := <-notifier:
		handler(eventFilter, eventCallbackUrl, ccEvent)
		break
	case <-time.After(time.Second * 20):
		handler(eventFilter, eventCallbackUrl, nil)
		break
	}

	return nil
}

func (setup *Setup) Query(fcn string, args [][]byte) (channel.Response, error) {
	request := channel.Request{
		ChaincodeID: setup.ChaincodeID,
		Fcn:         fcn,
		Args:        args,
	}

	response, err := setup.Client.Query(request)
	return response, err
}
