package networking

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/NethermindEth/posmoni/internal/utils"
	log "github.com/sirupsen/logrus"
)

// ExecutionClient : Struct ExecutionAPI interface implementation
type ExecutionClient struct {
	// Execution node endpoint to connect to
	Endpoint string
	// Time between retries when a request fails
	RetryDuration time.Duration
}

type Eth1Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err Eth1Error) Error() string {
	return fmt.Sprintf("Error %d (%s)", err.Code, err.Message)
}

// Inspired by https://github.com/onrik/ethrpc/blob/master/ethrpc.go
/*
Call :
Logic for calling a ETH json-rpc method.

params :-
a. method string

b. params []any

returns :-
a. json.RawMessage
Result field of the json response
b. error
Error if any
*/
func (ec *ExecutionClient) Call(method string, params ...any) (json.RawMessage, error) {
	request := eth1Request{
		ID:      1,
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	response, err := utils.PostRequest(ec.Endpoint, "application/json", bytes.NewBuffer(body), true, ec.RetryDuration)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf(ReadBodyError, err)
	}

	var resp eth1Response
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, *resp.Error
	}

	return resp.Result, nil
}

/*
SyncStatus :
Check sync status of the execution client using the json-rpc API method 'eth_syncing'.

params :-
none

returns :-
a. ExecutionSyncingStatus
Sync status of the execution client
*/
func (ec *ExecutionClient) SyncStatus() ExecutionSyncingStatus {
	result, err := ec.Call("eth_syncing")
	if err != nil {
		return ExecutionSyncingStatus{Endpoint: ec.Endpoint, Error: err}
	}
	log.Debugf("Result: %s", string(result))

	var ess ExecutionSyncingStatus
	ess, err = unmarshalData(result, ess)
	if err != nil && err.Error() != "json: cannot unmarshal bool into Go value of type networking.ExecutionSyncingStatus" {
		return ExecutionSyncingStatus{Endpoint: ec.Endpoint, Error: err}
	}

	// If is not syncing (is synced), result is 'false'
	ess.IsSyncing = ess.CurrentBlock != ""

	return ess
}
