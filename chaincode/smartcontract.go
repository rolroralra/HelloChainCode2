package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a Product
type SmartContract struct {
	contractapi.Contract
}

type Product struct {
	ID          string    `json:"ID"`
	ModelID     string    `json:"modelID"`
	ModelName   string    `json:"modelName"`
	Make        string    `json:"make"`
	Status      int       `json:"status"`
	UpdatedAt   string `json:"updatedAt"`
	Description string       `json:"description"`
}

// InitLedger adds a base set of products to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	products := []Product{
		{ID: "PRODUCT-00001", ModelID: "MODEL-00001", ModelName: "GalaxyS7", Make: "SAMSUNG", Status: 1, UpdatedAt: "2020-06-08", Description: "등록"},
		{ID: "PRODUCT-00002", ModelID: "MODEL-00002", ModelName: "GalaxyS9", Make: "SAMSUNG", Status: 1, UpdatedAt: "2020-06-08", Description: "등록"},
		{ID: "PRODUCT-00003", ModelID: "MODEL-00003", ModelName: "GalaxyS10", Make: "SAMSUNG", Status: 1, UpdatedAt: "2020-06-08", Description: "등록"},
		{ID: "PRODUCT-00004", ModelID: "MODEL-00004", ModelName: "GalaxyS11", Make: "SAMSUNG", Status: 1, UpdatedAt: "2020-06-08", Description: "등록"},
		{ID: "PRODUCT-00005", ModelID: "MODEL-00005", ModelName: "GalaxyS20", Make: "SAMSUNG", Status: 1, UpdatedAt: "2020-06-08", Description: "등록"},
		{ID: "PRODUCT-00006", ModelID: "MODEL-00006", ModelName: "GalaxyS20", Make: "SAMSUNG", Status: 1, UpdatedAt: "2020-06-08", Description: "등록"},
		{ID: "PRODUCT-00007", ModelID: "MODEL-00007", ModelName: "GalaxyS20", Make: "SAMSUNG", Status: 1, UpdatedAt: "2020-06-08", Description: "등록"},
	}

	for _, product := range products {
		productJSON, err := json.Marshal(product)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(product.ID, productJSON)
		if err != nil {
			return fmt.Errorf("Failed to put to world state. %v", err)
		}
	}

	return nil

}

// QueryProduct returns the product stored in the world state with given id.
func (s *SmartContract) QueryProduct(ctx contractapi.TransactionContextInterface, id string) (*Product, error) {
	productJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state: %v", err)
	}
	if productJSON == nil {
		return nil, fmt.Errorf("The product %s does not exist", id)
	}

	var product Product
	err = json.Unmarshal(productJSON, &product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// In CouchDB,QueryProduct returns the product stored in the world state with given id.
func (s *SmartContract) QueryProductCouchDB(ctx contractapi.TransactionContextInterface, query string) ([]*Product, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state: %v", err)
	}
	if resultsIterator == nil {
		return nil, fmt.Errorf("The product does not exist")
	}

	defer resultsIterator.Close()

	var products []*Product
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var product Product
		err = json.Unmarshal(queryResponse.Value, &product)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	return products, nil
}

// QueryAllProducts returns all products found in world state
func (s *SmartContract) QueryAllProducts(ctx contractapi.TransactionContextInterface) ([]*Product, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all products in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var products []*Product
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var product Product
		err = json.Unmarshal(queryResponse.Value, &product)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}

// AddProduct issues a new product to the world state with given details.
func (s *SmartContract) AddProduct(ctx contractapi.TransactionContextInterface, id string, modelID string, modelName string, make string, status int, updatedAt string, description string) error {
	exists, err := s.ProductExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("The product %s already exists", id)
	}

	product := Product{
		ID:    id,
		ModelID:  modelID,
		ModelName: modelName,
		Make: make,
		Status: status,
		UpdatedAt: updatedAt,
		Description: description,
	}

	productJSON, err := json.Marshal(product)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(id, productJSON)
	if err != nil {
		return fmt.Errorf("Failed to put to world state. %v", err)
	}
	return nil
}

// ChangeOwner updates the owner field of product with given id in world state.
//func (s *SmartContract) ChangeOwner(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
//	product, err := s.QueryProduct(ctx, id)
//	if err != nil {
//		return err
//	}
//
//	//product.Owner = newOwner
//	productJSON, err := json.Marshal(product)
//	if err != nil {
//		return err
//	}
//
//	return ctx.GetStub().PutState(id, productJSON)
//}

// ProductExists returns true when product with given ID exists in world state
func (s *SmartContract) ProductExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	productJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("Failed to read from world state: %v", err)
	}

	return productJSON != nil, nil
}

// queryHistoryProducts
func (s *SmartContract) QueryHistoryProducts(ctx contractapi.TransactionContextInterface, id string) ([]*Product, error) {

	historyIer, error := ctx.GetStub().GetHistoryForKey(id)

	if error != nil {
		return nil, error
	}

	var products []*Product
	for historyIer.HasNext() {
		queryResponse, err := historyIer.Next()
		var product Product
		if err != nil {
			return nil, err
		}
		if queryResponse.IsDelete {
			continue
		} else {
			err = json.Unmarshal(queryResponse.Value, &product)
			if err != nil {
				return nil, err
			}
		}

		products = append(products, &product)
	}

	return products, nil
}

// deleteProduct
func (s *SmartContract) DeleteProduct(ctx contractapi.TransactionContextInterface, id string) error {

	_, err := s.QueryProduct(ctx, id)
	if err != nil {
		return err
	}

	return ctx.GetStub().DelState(id)
}
