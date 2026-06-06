package ports

type Consumer interface {
	Consume(handler func([]byte) error) error
}	