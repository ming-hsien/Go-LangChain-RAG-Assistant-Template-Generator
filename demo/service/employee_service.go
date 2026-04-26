package service

import (
	"fmt"

	"github.com/ming-hsien/lang-chain-template/internal/employee_status"
)

type EmployeeService struct {
	statusClient *employee_status.Client
}

func NewEmployeeService() *EmployeeService {
	return &EmployeeService{
		statusClient: employee_status.NewClient(),
	}
}

// CheckAvailability orchestrates the logic to check an employee's status
func (s *EmployeeService) CheckAvailability(name string) string {
	status, err := s.statusClient.FetchStatus(name)
	if err != nil {
		return fmt.Sprintf("Error checking availability for %s: %v", name, err)
	}

	return fmt.Sprintf("Employee Status [%s]: %s", name, status)
}
