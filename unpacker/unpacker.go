package unpacker

// Unpacker interface is implemented by all supported
// archive files unpackers
type Unpacker interface {
	Unpack() error
}
