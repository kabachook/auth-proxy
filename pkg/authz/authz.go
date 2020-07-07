package authz

type Authz interface {
	Authorize(username string) (bool, error)
}
