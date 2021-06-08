package chaincode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a Model
type SmartContract struct {
	contractapi.Contract
}

// Model (재고관리-- 모델 기준)
type Model struct {
	ID    string `json:"ID"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// Product (Seq = 0 : MetaData, Seq >= 1 : Product History)
type Product struct {
	ID          string    `json:"ID"`
	Seq         int       `json:"seq"`
	ModelID     string    `json:"modelID"`
	Status      int       `json:"status"`
	UpdateAt    time.Time `json:"updateAt"`
	Description int       `json:"description"`
}

// InitLedger adds a base set of models to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	models := []Model{
		{ID: "MODEL-00001", Name: "GalaxyS10", Count: 10},
		{ID: "MODEL-00002", Name: "GalaxyS11", Count: 20},
		{ID: "MODEL-00003", Name: "GalaxyS20", Count: 1000},
		{ID: "MODEL-00004", Name: "GalaxyS21", Count: 2000},
		{ID: "MODEL-00005", Name: "GalaxyS22", Count: 0},
	}

	for _, model := range models {
		modelJSON, err := json.Marshal(model)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(model.ID, modelJSON)
		if err != nil {
			return fmt.Errorf("Failed to put to world state. %v", err)
		}
	}

	return nil

}

// QueryModel returns the model stored in the world state with given id.
func (s *SmartContract) QueryModel(ctx contractapi.TransactionContextInterface, id string) (*Model, error) {
	modelJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state: %v", err)
	}
	if modelJSON == nil {
		return nil, fmt.Errorf("The model %s does not exist", id)
	}

	var model Model
	err = json.Unmarshal(modelJSON, &model)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

// In CouchDB,QueryModel returns the model stored in the world state with given id.
func (s *SmartContract) QueryModelCouchDB(ctx contractapi.TransactionContextInterface, query string) ([]*Model, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state: %v", err)
	}
	if resultsIterator == nil {
		return nil, fmt.Errorf("The model does not exist")
	}

	defer resultsIterator.Close()

	var models []*Model
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var model Model
		err = json.Unmarshal(queryResponse.Value, &model)
		if err != nil {
			return nil, err
		}
		models = append(models, &model)
	}
	return models, nil
}

// QueryAllModels returns all models found in world state
func (s *SmartContract) QueryAllModels(ctx contractapi.TransactionContextInterface) ([]*Model, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all models in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var models []*Model
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var model Model
		err = json.Unmarshal(queryResponse.Value, &model)
		if err != nil {
			return nil, err
		}
		models = append(models, &model)
	}

	return models, nil
}

// AddModel issues a new model to the world state with given details.
func (s *SmartContract) AddModel(ctx contractapi.TransactionContextInterface, id string, name string, count int) error {
	exists, err := s.ModelExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("The model %s already exists", id)
	}

	model := Model{
		ID:    id,
		Name:  name,
		Count: count,
	}
	modelJSON, err := json.Marshal(model)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(id, modelJSON)
	if err != nil {
		return fmt.Errorf("Failed to put to world state. %v", err)
	}
	return nil
}

// ChangeOwner updates the owner field of model with given id in world state.
func (s *SmartContract) ChangeOwner(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
	model, err := s.QueryModel(ctx, id)
	if err != nil {
		return err
	}

	//model.Owner = newOwner
	modelJSON, err := json.Marshal(model)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, modelJSON)
}

// ModelExists returns true when model with given ID exists in world state
func (s *SmartContract) ModelExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	modelJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("Failed to read from world state: %v", err)
	}

	return modelJSON != nil, nil
}

// queryHistoryModels
func (s *SmartContract) QueryHistoryModels(ctx contractapi.TransactionContextInterface, id string) ([]*Model, error) {

	historyIer, error := ctx.GetStub().GetHistoryForKey(id)

	if error != nil {
		return nil, error
	}

	var models []*Model
	for historyIer.HasNext() {
		queryResponse, err := historyIer.Next()
		var model Model
		if err != nil {
			return nil, err
		}
		if queryResponse.IsDelete {
			continue
		} else {
			err = json.Unmarshal(queryResponse.Value, &model)
			if err != nil {
				return nil, err
			}
		}

		models = append(models, &model)
	}

	return models, nil
}

// deleteModel
func (s *SmartContract) DeleteModel(ctx contractapi.TransactionContextInterface, id string) error {

	_, err := s.QueryModel(ctx, id)
	if err != nil {
		return err
	}

	return ctx.GetStub().DelState(id)
}

//pushModel
func (s *SmartContract) PushModel(ctx contractapi.TransactionContextInterface, id string, count int) error {
	model, err := s.QueryModel(ctx, id)
	if err != nil {
		return err
	}

	if count <= 0 {
		return fmt.Errorf("Non positive count. %v < %v", model.Count, count)
	}

	model.Count += count
	modelJSON, err := json.Marshal(model)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, modelJSON)
}

//popModel
func (s *SmartContract) PopModel(ctx contractapi.TransactionContextInterface, id string, count int) error {
	model, err := s.QueryModel(ctx, id)
	if err != nil {
		return err
	}

	if count <= 0 {
		return fmt.Errorf("Non positive count. %v < %v", model.Count, count)
	}

	if model.Count < count {
		return fmt.Errorf("Not enough count. %v < %v", model.Count, count)
	}

	model.Count -= count
	modelJSON, err := json.Marshal(model)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, modelJSON)
}
