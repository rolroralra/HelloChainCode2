package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a Car
type SmartContract struct {
	contractapi.Contract
}

// Car describes basic details of what makes up a simple car
type Car struct {
	ID     string `json:"ID"`
	Make   string `json:"make"`
	Model  string `json:"model"`
	Count  int    `json:"count"`
	Owner  string `json:"owner"`
}

// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	cars := []Car{
		Car{ID: "CAR0", Make: "Toyota", Model: "Prius", Count: 10, Owner: "Tomoko"},
		Car{ID: "CAR1", Make: "Ford", Model: "Mustang", Count: 200, Owner: "Brad"},
		Car{ID: "CAR2", Make: "Hyundai", Model: "Tucson", Count: 320, Owner: "Jin Soo"},
		Car{ID: "CAR3", Make: "Volkswagen", Model: "Passat", Count: 25, Owner: "Max"},
		Car{ID: "CAR4", Make: "Tesla", Model: "S", Count: 15, Owner: "Adriana"},
	}

	for _, car := range cars {
		carJSON, err := json.Marshal(car)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(car.ID, carJSON)
		if err != nil {
			return fmt.Errorf("Failed to put to world state. %v", err)
		}
	}

	return nil

}

// QueryCar returns the car stored in the world state with given id.
func (s *SmartContract) QueryCar(ctx contractapi.TransactionContextInterface, id string) (*Car, error) {
	carJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state: %v", err)
	}
	if carJSON == nil {
		return nil, fmt.Errorf("The car %s does not exist", id)
	}

	var car Car
	err = json.Unmarshal(carJSON, &car)
	if err != nil {
		return nil, err
	}

	return &car, nil
}

// In CouchDB,QueryCar returns the car stored in the world state with given id.
func (s *SmartContract) QueryCarCouchDB(ctx contractapi.TransactionContextInterface, query string) ([]*Car, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state: %v", err)
	}
	if resultsIterator == nil {
		return nil, fmt.Errorf("The car does not exist")
	}

	defer resultsIterator.Close()

	var cars []*Car
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var car Car
		err = json.Unmarshal(queryResponse.Value, &car)
		if err != nil {
			return nil, err
		}
		cars = append(cars, &car)
	}
	return cars, nil
}

// QueryAllCars returns all cars found in world state
func (s *SmartContract) QueryAllCars(ctx contractapi.TransactionContextInterface) ([]*Car, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all cars in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var cars []*Car
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var car Car
		err = json.Unmarshal(queryResponse.Value, &car)
		if err != nil {
			return nil, err
		}
		cars = append(cars, &car)
	}

	return cars, nil
}

// AddCar issues a new car to the world state with given details.
func (s *SmartContract) AddCar(ctx contractapi.TransactionContextInterface, id string, make string, model string, count int, owner string) error {
	exists, err := s.CarExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("The car %s already exists", id)
	}

	car := Car{
		ID:     id,
		Make:   make,
		Model:  model,
		Count: count,
		Owner:  owner,
	}
	carJSON, err := json.Marshal(car)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(id, carJSON)
	if err != nil {
		return fmt.Errorf("Failed to put to world state. %v", err)
	}
	return nil
}

// ChangeOwner updates the owner field of car with given id in world state.
func (s *SmartContract) ChangeOwner(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
	car, err := s.QueryCar(ctx, id)
	if err != nil {
		return err
	}

	car.Owner = newOwner
	carJSON, err := json.Marshal(car)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, carJSON)
}

// CarExists returns true when car with given ID exists in world state
func (s *SmartContract) CarExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	carJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("Failed to read from world state: %v", err)
	}

	return carJSON != nil, nil
}

// queryHistoryCars
func (s *SmartContract) QueryHistoryCars(ctx contractapi.TransactionContextInterface, id string) ([]*Car, error) {

	historyIer, error := ctx.GetStub().GetHistoryForKey(id)

	if error != nil {
		return nil, error
	}

	var cars []*Car
	for historyIer.HasNext() {
		queryResponse, err := historyIer.Next()
		var car Car
		if err != nil {
			return nil, err
		}
		if queryResponse.IsDelete {
			continue
		} else {
			err = json.Unmarshal(queryResponse.Value, &car)
			if err != nil {
				return nil, err
			}
		}

		cars = append(cars, &car)
	}

	return cars, nil
}

// deleteCar
func (s *SmartContract) DeleteCar(ctx contractapi.TransactionContextInterface, id string) error {

	_, err := s.QueryCar(ctx, id)
	if err != nil {
		return err
	}

	return ctx.GetStub().DelState(id)
}

//pushCar 
func (s *SmartContract) PushCar(ctx contractapi.TransactionContextInterface, id string, count int) error {
	car, err := s.QueryCar(ctx, id)
	if err != nil {
		return err
	}

	if count <= 0 {
		return fmt.Errorf("Non positive count. %v < %v", car.Count, count)
  }

	car.Count += count
	carJSON, err := json.Marshal(car)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, carJSON)
}

//popCar
func (s *SmartContract) PopCar(ctx contractapi.TransactionContextInterface, id string, count int) error {
	car, err := s.QueryCar(ctx, id)
	if err != nil {
		return err
	}

	if count <= 0 {
		return fmt.Errorf("Non positive count. %v < %v", car.Count, count)
  }

	if car.Count < count {
		return fmt.Errorf("Not enough count. %v < %v", car.Count, count)
  }

	car.Count -= count
	carJSON, err := json.Marshal(car)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, carJSON)
}
