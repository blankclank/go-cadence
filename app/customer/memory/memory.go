package memory

import (
	"fmt"
	"go-cadence/tavern/customer"
	"sync"
)

type CustomerMemoryRepository struct {
	customers map[string]customer.Customer
	sync.Mutex
}

func New() *CustomerMemoryRepository {
	return &CustomerMemoryRepository{
		customers: make(map[string]customer.Customer),
	}
}

func (mr *CustomerMemoryRepository) Get(id string) (customer.Customer, error) {
	if customer, ok := mr.customers[id]; ok {
		return customer, nil
	}

	return customer.Customer{}, customer.ErrCustomerNotFound
}

func (mr *CustomerMemoryRepository) Update(c customer.Customer) error {
	if _, ok := mr.customers[c.Name]; !ok {
		return fmt.Errorf("customer does not exits: %w", customer.ErrUpdateCustomer)
	}

	mr.Lock()
	mr.customers[c.Name] = c
	mr.Unlock()
	return nil
}
