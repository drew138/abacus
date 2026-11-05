package application

import (
	"context"
	"fmt"
)

func (a *Application) CheckRegistration(ctx context.Context) (bool, error) {
	users, err := a.repo.CountUsers(ctx)
	if err != nil {
		return false, fmt.Errorf("checking for users: %w", err)
	}
	return users == 0, nil
}
