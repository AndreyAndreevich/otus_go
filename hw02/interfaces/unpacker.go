package interfaces

// Unpacker is interface of unpacking string
type Unpacker interface {
	Unpack(string) (string, error)
}
