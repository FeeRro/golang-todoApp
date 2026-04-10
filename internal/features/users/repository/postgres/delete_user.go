package users_postgres_repository

import (
	"context"
	"fmt"

	core_error "github.com/FeeRro/golang-todoApp/internal/core/errors"
)

func (r UserRepository) DeleteUser(
	ctx context.Context,
	id int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM todolist.users
	WHERE id=$1;
	`

	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("User whit id='%d': %w", id, core_error.ErrNotFound)
	}

	return nil
}
