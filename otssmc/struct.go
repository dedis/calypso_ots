package otssmc

import (
	"gopkg.in/dedis/cothority.v1/skipchain"
	"gopkg.in/dedis/crypto.v0/abstract"
	"gopkg.in/dedis/onet.v1"
	"gopkg.in/dedis/onet.v1/crypto"
)

type DecryptRequestData struct {
	WriteSBF       *skipchain.SkipBlockFix
	ReadSBF        *skipchain.SkipBlockFix
	InclusionProof *skipchain.BlockLink
	ACPublicKeys   []abstract.Point
}

type DecryptRequest struct {
	RootIndex int
	Roster    *onet.Roster
	Data      *DecryptRequestData
	Signature *crypto.SchnorrSig
}

type DecryptReply struct {
	DecShares []*DecryptedShare
}

type DecryptedShare struct {
	K  abstract.Point
	Cs []abstract.Point
}

// Protocol messages

type AnnounceDecrypt struct {
	DecReqData *DecryptRequestData
	Signature  *crypto.SchnorrSig
	RootIndex  int
}

type StructAnnounceDecrypt struct {
	*onet.TreeNode
	AnnounceDecrypt
}

type ShareReply struct {
	DecShare *DecryptedShare
}

type StructShareReply struct {
	*onet.TreeNode
	ShareReply
}
