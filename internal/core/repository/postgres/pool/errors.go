package core_postgres_pool

import "errors"

var (
	ErrNoRows             = errors.New("no rows")
	ErrViolatesForeignKey = errors.New("violetes foreign key")
	ErrUnknown            = errors.New("unknown")
)
