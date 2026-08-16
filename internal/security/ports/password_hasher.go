package ports

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(storedHash string, password string) error
}
