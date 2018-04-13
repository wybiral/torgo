package torgo

type Onion struct {
	Ports map[int]int
	ServiceID string
	PrivateKey string
	PrivateKeyType string
}
