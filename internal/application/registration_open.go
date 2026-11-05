package application

import (
	"context"
	"fmt"
)

func (a *Application) RegistrationOpen(ctx context.Context) (bool, error) {
	taken, err := a.users.HasUsers(ctx)
	if err != nil {
		return false, fmt.Errorf("checking for users: %w", err)
	}
	return !taken, nil
}
