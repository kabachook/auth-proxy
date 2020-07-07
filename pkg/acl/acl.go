package acl

type ACL interface {
	// Check: checks if user with provided username and route is allowed
	Check(username string, route string) bool
}