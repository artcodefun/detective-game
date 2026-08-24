package ports

import "context"

// TransactionManager runs fn in a database transaction.
//
// The transaction is committed when fn returns nil and rolled back otherwise.
// Repositories must use txCtx for every operation that belongs to the same
// atomic change.
type TransactionManager interface {
	WithTx(ctx context.Context, fn func(txCtx context.Context) error) error
}
