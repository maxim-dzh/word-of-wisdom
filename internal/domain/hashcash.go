package domain

// HashcashHeader is an object representation of the hashcash header
type HashcashHeader struct {
	Version   int
	Bits      uint
	Timestamp int64
	Resource  string
	Random    string
	Counter   uint64
}
