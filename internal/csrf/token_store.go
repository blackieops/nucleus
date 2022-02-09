package csrf

type TokenStore interface {
	Consume(string) error
	Generate() (string, error)
}
