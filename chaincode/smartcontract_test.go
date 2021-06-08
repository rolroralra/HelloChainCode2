package chaincode_test

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type transactionContext interface {
	contractapi.TransactionContextInterface
}

type chaincodeStub interface {
	shim.ChaincodeStubInterface
}
//
//func TestQuaryAllProducts(t *testing.T) {
//	chaincodeStub := &mocks.ChaincodeStub{}
//	transactionContext := &mocks.TransactionContext{}
//	transactionContext.GetStubReturns(chaincodeStub)
//
//	//expectedAsset := &chaincode.Model{ID: "Stock1"}
//	// expectedAsset1 := &chaincode.Stock{ID: "Stock2"}
//	//bytes, err := json.Marshal(expectedAsset)
//	//require.NoError(t, err)
//
//	chaincodeStub.GetStateReturns(bytes, nil)
//	assetTransfer := chaincode.SmartContract{}
//	asset, err := assetTransfer.QueryAllProducts(transactionContext)
//	require.NoError(t, err)
//	//fmt.Print(expectedAsset)
//	fmt.Print(asset)
//	//require.Equal(t, expectedAsset, asset)
//}

//func TestAddStock(t *testing.T) {
//	chaincodeStub := &mocks.ChaincodeStub{}
//	transactionContext := &mocks.TransactionContext{}
//	transactionContext.GetStubReturns(chaincodeStub)
//
//	assetTransfer := chaincode.SmartContract{}
//	err := assetTransfer.AddModel(transactionContext, "MODEL-00001", "KIA", 10)
//	require.NoError(t, err)
//
//	chaincodeStub.PutStateReturns(fmt.Errorf("MODEL-00001"))
//	err = assetTransfer.AddModel(transactionContext, "MODEL-00001", "KIA", 3)
//	require.EqualError(t, err, "Failed to put to world state. MODEL-00001")
//}
//
//func TestQueryModel(t *testing.T) {
//	chaincodeStub := &mocks.ChaincodeStub{}
//	transactionContext := &mocks.TransactionContext{}
//	transactionContext.GetStubReturns(chaincodeStub)
//
//	expectedAsset := &chaincode.Model{ID: "Stock1"}
//	// expectedAsset1 := &chaincode.Stock{ID: "Stock2"}
//	bytes, err := json.Marshal(expectedAsset)
//	require.NoError(t, err)
//
//	chaincodeStub.GetStateReturns(bytes, nil)
//	assetTransfer := chaincode.SmartContract{}
//	asset, err := assetTransfer.QueryModel(transactionContext, "")
//	require.NoError(t, err)
//	fmt.Print(expectedAsset)
//	fmt.Print(asset)
//	require.Equal(t, expectedAsset, asset)
//}
