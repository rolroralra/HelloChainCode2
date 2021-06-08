package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a Stock
type SmartContract struct {
	contractapi.Contract
}

// Stock describes basic details of what makes up a simple stock
type Stock struct {
	ID    string `json:"ID"`
	Make  string `json:"make"`
	Model string `json:"model"`
	Count int    `json:"count"`
	Owner string `json:"owner"`
}

// InitLedger adds a base set of stocks to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	stocks := []Stock{
		Stock{ID: "CAR0", Make: "Toyota", Model: "Prius", Count: 10, Owner: "Tomoko"},
		Stock{ID: "CAR1", Make: "Ford", Model: "Mustang", Count: 200, Owner: "Brad"},
		Stock{ID: "CAR2", Make: "Hyundai", Model: "Tucson", Count: 320, Owner: "Jin Soo"},
		Stock{ID: "CAR3", Make: "Volkswagen", Model: "Passat", Count: 25, Owner: "Max"},
		Stock{ID: "CAR4", Make: "Tesla", Model: "S", Count: 15, Owner: "Adriana"},
	}

	for _, stock := range stocks {
		stockJSON, err := json.Marshal(stock)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(stock.ID, stockJSON)
		if err != nil {
			return fmt.Errorf("Failed to put to world state. %v", err)
		}
	}

	return nil

}

// QueryStock returns the stock stored in the world state with given id.
func (s *SmartContract) QueryStock(ctx contractapi.TransactionContextInterface, id string) (*Stock, error) {
	stockJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state: %v", err)
	}
	if stockJSON == nil {
		return nil, fmt.Errorf("The stock %s does not exist", id)
	}

	var stock Stock
	err = json.Unmarshal(stockJSON, &stock)
	if err != nil {
		return nil, err
	}

	return &stock, nil
}

// In CouchDB,QueryStock returns the stock stored in the world state with given id.
func (s *SmartContract) QueryStockCouchDB(ctx contractapi.TransactionContextInterface, query string) ([]*Stock, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state: %v", err)
	}
	if resultsIterator == nil {
		return nil, fmt.Errorf("The stock does not exist")
	}

	defer resultsIterator.Close()

	var stocks []*Stock
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var stock Stock
		err = json.Unmarshal(queryResponse.Value, &stock)
		if err != nil {
			return nil, err
		}
		stocks = append(stocks, &stock)
	}
	return stocks, nil
}

// QueryAllStocks returns all stocks found in world state
func (s *SmartContract) QueryAllStocks(ctx contractapi.TransactionContextInterface) ([]*Stock, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all stocks in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var stocks []*Stock
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var stock Stock
		err = json.Unmarshal(queryResponse.Value, &stock)
		if err != nil {
			return nil, err
		}
		stocks = append(stocks, &stock)
	}

	return stocks, nil
}

// AddStock issues a new stock to the world state with given details.
func (s *SmartContract) AddStock(ctx contractapi.TransactionContextInterface, id string, make string, model string, count int, owner string) error {
	exists, err := s.StockExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("The stock %s already exists", id)
	}

	stock := Stock{
		ID:    id,
		Make:  make,
		Model: model,
		Count: count,
		Owner: owner,
	}
	stockJSON, err := json.Marshal(stock)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(id, stockJSON)
	if err != nil {
		return fmt.Errorf("Failed to put to world state. %v", err)
	}
	return nil
}

// ChangeOwner updates the owner field of stock with given id in world state.
func (s *SmartContract) ChangeOwner(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
	stock, err := s.QueryStock(ctx, id)
	if err != nil {
		return err
	}

	stock.Owner = newOwner
	stockJSON, err := json.Marshal(stock)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, stockJSON)
}

// StockExists returns true when stock with given ID exists in world state
func (s *SmartContract) StockExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	stockJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("Failed to read from world state: %v", err)
	}

	return stockJSON != nil, nil
}

// queryHistoryStocks
func (s *SmartContract) QueryHistoryStocks(ctx contractapi.TransactionContextInterface, id string) ([]*Stock, error) {

	historyIer, error := ctx.GetStub().GetHistoryForKey(id)

	if error != nil {
		return nil, error
	}

	var stocks []*Stock
	for historyIer.HasNext() {
		queryResponse, err := historyIer.Next()
		var stock Stock
		if err != nil {
			return nil, err
		}
		if queryResponse.IsDelete {
			continue
		} else {
			err = json.Unmarshal(queryResponse.Value, &stock)
			if err != nil {
				return nil, err
			}
		}

		stocks = append(stocks, &stock)
	}

	return stocks, nil
}

// deleteStock
func (s *SmartContract) DeleteStock(ctx contractapi.TransactionContextInterface, id string) error {

	_, err := s.QueryStock(ctx, id)
	if err != nil {
		return err
	}

	return ctx.GetStub().DelState(id)
}

//pushStock
func (s *SmartContract) PushStock(ctx contractapi.TransactionContextInterface, id string, count int) error {
	stock, err := s.QueryStock(ctx, id)
	if err != nil {
		return err
	}

	if count <= 0 {
		return fmt.Errorf("Non positive count. %v < %v", stock.Count, count)
	}

	stock.Count += count
	stockJSON, err := json.Marshal(stock)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, stockJSON)
}

//popStock
func (s *SmartContract) PopStock(ctx contractapi.TransactionContextInterface, id string, count int) error {
	stock, err := s.QueryStock(ctx, id)
	if err != nil {
		return err
	}

	if count <= 0 {
		return fmt.Errorf("Non positive count. %v < %v", stock.Count, count)
	}

	if stock.Count < count {
		return fmt.Errorf("Not enough count. %v < %v", stock.Count, count)
	}

	stock.Count -= count
	stockJSON, err := json.Marshal(stock)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, stockJSON)
}
