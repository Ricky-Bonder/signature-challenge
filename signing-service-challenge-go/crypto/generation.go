package crypto

type KeyPairGenerator interface {
	Generate() (interface{}, error)
	GetAlgorithm() string
}
