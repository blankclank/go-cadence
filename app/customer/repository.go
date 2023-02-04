package customer

import "errors"

var (
	ErrCustomerNotFound = errors.New("the customer was not found")
	ErrUpdateCustomer   = errors.New("failed to update the customer")
)

type CustomerRepository interface {
	Get(string) (Customer, error)
	Update(Customer) error
}
