package chaincode_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-samples/asset-transfer-fabcar/chaincode-go/chaincode"
	"github.com/hyperledger/fabric-samples/asset-transfer-fabcar/chaincode-go/chaincode/mocks"
	"github.com/stretchr/testify/require"
)

type transactionContext interface {
	contractapi.TransactionContextInterface
}

type chaincodeStub interface {
	shim.ChaincodeStubInterface
}

func TestAddStock(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	assetTransfer := chaincode.SmartContract{}
	err := assetTransfer.AddModel(transactionContext, "MODEL-00001", "KIA", 10)
	require.NoError(t, err)

	chaincodeStub.PutStateReturns(fmt.Errorf("MODEL-00001"))
	err = assetTransfer.AddModel(transactionContext, "MODEL-00001", "KIA", 3)
	require.EqualError(t, err, "Failed to put to world state. MODEL-00001")
}

func TestQueryModel(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	expectedAsset := &chaincode.Model{ID: "Stock1"}
	// expectedAsset1 := &chaincode.Stock{ID: "Stock2"}
	bytes, err := json.Marshal(expectedAsset)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	assetTransfer := chaincode.SmartContract{}
	asset, err := assetTransfer.QueryModel(transactionContext, "")
	require.NoError(t, err)
	fmt.Print(expectedAsset)
	fmt.Print(asset)
	require.Equal(t, expectedAsset, asset)
}
