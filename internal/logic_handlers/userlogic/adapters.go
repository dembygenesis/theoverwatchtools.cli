package userlogic

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . persistor
type persistor interface {
	//GetUsers(ctx context.Context, tx persistence.TransactionHandler, filters *model.UserFilters) (*model.PaginatedUsers, error)
}
