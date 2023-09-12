package permission

// Permission constants
// You should not care about the values of these constants.
// When adding a new permission, it always has to be done at the end of the list.
// Because else the permissions of the existing users will be messed up.
// Also don't write explicit values for the constants

const (
	// ACCOUNTS_READ_ALL is the permission to read all accounts
	ACCOUNTS_READ_ALL uint64 = 1 << iota
	// ACCOUNTS_MANAGE_ALL is the permission to manage all accounts (includes read)
	ACCOUNTS_MANAGE_ALL
)
