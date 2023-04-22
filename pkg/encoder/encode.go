package encoder

type Encoder struct {
}

const (
	base         = 62
	characterSet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

// New is a constructor function
func New() *Encoder {
	return &Encoder{}
}

// EncodeToBase62 accepts a number and returns it as base62 string
func (e Encoder) EncodeToBase62(number uint64) string {
	encoded := ""
	for number > 0 {
		r := number % base
		number /= base
		encoded = string(characterSet[r]) + encoded
	}

	return encoded
}
