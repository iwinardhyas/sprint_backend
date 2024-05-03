package utils

import (
	"fmt"

	"github.com/iwinardhyas/sprint_backend/pkg/repository"
)

// GetCredentialsByRole func for getting credentials from a role name.
func GetCredentialsByRole(role string) ([]string, error) {
	// Define credentials variable.
	var credentials []string

	// Switch given role.
	switch role {
	case repository.AdminRoleName:
		// Admin credentials (all access).
		credentials = []string{
			repository.UserRoleName,
			repository.UserRoleName,
			repository.UserRoleName,
		}
	case repository.ModeratorRoleName:
		credentials = []string{
			repository.UserRoleName,
			repository.UserRoleName,
		}
	case repository.UserRoleName:
		credentials = []string{
			repository.UserRoleName,
		}
	default:
		// Return error message.
		return nil, fmt.Errorf("role '%v' does not exist", role)
	}

	return credentials, nil
}
