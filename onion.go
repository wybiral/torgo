package torgo

type Onion struct {
	Ports map[int]string
	ServiceID string
	PrivateKey string
	PrivateKeyType string
}
