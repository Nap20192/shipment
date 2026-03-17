package kernel

type DomainEvent interface {
	Name() string
	Payload() []byte
}
