package interfaces

type Unpacker interface {
	Unpack(string) (string, error)
}
