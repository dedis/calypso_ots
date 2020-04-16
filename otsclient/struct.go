package otsclient

import (
	"gopkg.in/dedis/crypto.v0/abstract"
	"gopkg.in/dedis/crypto.v0/share/pvss"
)

type DataPVSS struct {
	NumTrustee   int
	Threshold    int
	Suite        abstract.Suite
	G            abstract.Point
	H            abstract.Point
	Secret       abstract.Scalar
	SCPublicKeys []abstract.Point
	EncShares    []*pvss.PubVerShare
	EncProofs    []abstract.Point
}
