package utils

import "github.com/ethereum/go-ethereum/rlp"

func EncodeToRLPs[T comparable](s []T) []rlp.RawValue {
	var raws []rlp.RawValue
	for _, v := range s {
		if enc, err := rlp.EncodeToBytes(v); err == nil {
			if len(enc) > 0 {
				raws = append(raws, enc)
			}
		}
	}
	return raws
}
