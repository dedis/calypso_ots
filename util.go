package ots

import (
	"crypto/sha256"

	"gopkg.in/dedis/crypto.v0/abstract"
	"gopkg.in/dedis/onet.v1/crypto"
	"gopkg.in/dedis/onet.v1/network"
)

func CreatePointH(suite abstract.Suite, pubKey abstract.Point) (abstract.Point, error) {
	binPubKey, err := pubKey.MarshalBinary()
	if err != nil {
		return nil, err
	}
	tmpHash := sha256.Sum256(binPubKey)
	labelHash := tmpHash[:]
	h, _ := suite.Point().Pick(nil, suite.Cipher(labelHash))
	return h, nil
}

func SignMessage(msg []byte, privKey abstract.Scalar) (crypto.SchnorrSig, error) {
	tmpHash := sha256.Sum256(msg)
	msgHash := tmpHash[:]
	return crypto.SignSchnorr(network.Suite, privKey, msgHash)
}
