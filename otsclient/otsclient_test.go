package otsclient

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/dedis/crypto.v0/abstract"
	"gopkg.in/dedis/crypto.v0/ed25519"
	"gopkg.in/dedis/crypto.v0/random"
	"gopkg.in/dedis/crypto.v0/share/pvss"
	"gopkg.in/dedis/onet.v1"
	"gopkg.in/dedis/onet.v1/log"
)

func TestSimple(t *testing.T) {
	numServers := 7
	local := onet.NewTCPTest()
	_, roster, _ := local.GenTree(numServers, true)
	defer local.CloseAll()
	require.NotNil(t, roster)

	scurl, err := CreateSkipchain(roster)
	require.NoError(t, err)
	require.NotNil(t, scurl)

	scPubKeys := roster.Publics()
	dataPVSS := DataPVSS{
		Suite:        ed25519.NewAES128SHA256Ed25519(false),
		SCPublicKeys: scPubKeys,
		NumTrustee:   numServers,
	}

	// Writer's pk/sk pair
	wrPrivKey := dataPVSS.Suite.Scalar().Pick(random.Stream)
	wrPubKey := dataPVSS.Suite.Point().Mul(nil, wrPrivKey)
	//Reader's pk/sk pair
	privKey := dataPVSS.Suite.Scalar().Pick(random.Stream)
	pubKey := dataPVSS.Suite.Point().Mul(nil, privKey)
	err = SetupPVSS(&dataPVSS, pubKey)
	require.NoError(t, err)

	mesgSize := 1024 * 1024
	mesg := make([]byte, mesgSize)
	for i := 0; i < mesgSize; i++ {
		mesg[i] = 'w'
	}
	encMesg, hashEnc, err := EncryptMessage(&dataPVSS, mesg)
	require.NoError(t, err)

	//// Creating write transaction
	writeSB, err := WriteRequest(scurl, &dataPVSS, hashEnc, pubKey, wrPrivKey)
	require.NoError(t, err)

	// Bob gets it from Alice
	writeID := writeSB.Hash
	// Get write transaction from skipchain
	writeSB, writeData, sig, err := GetWriteSB(scurl, writeID)
	require.NoError(t, err)

	sigVerErr := VerifyWriteSignature(writeData, sig, wrPubKey)
	require.NoError(t, sigVerErr)

	validHash := VerifyEncMesg(writeData, encMesg)
	require.Equal(t, validHash, 0)

	// Creating read transaction
	readSB, err := ReadRequest(scurl, writeID, privKey)
	require.NoError(t, err)

	updWriteSB, err := GetUpdatedWriteSB(scurl, writeID)
	require.NoError(t, err)

	acPubKeys := readSB.Roster.Publics()
	// Bob obtains the SC public keys from T_W
	scPubKeys = writeData.SCPublicKeys
	decShares, err := GetDecryptedShares(roster, updWriteSB, readSB.SkipBlockFix, acPubKeys, scPubKeys, privKey, readSB.Index)
	require.NoError(t, err)

	var validKeys []abstract.Point
	var validEncShares []*pvss.PubVerShare
	var validDecShares []*pvss.PubVerShare
	sz := len(decShares)
	for i := 0; i < sz; i++ {
		validKeys = append(validKeys, writeData.SCPublicKeys[i])
		validEncShares = append(validEncShares, writeData.EncShares[i])
		validDecShares = append(validDecShares, decShares[i])
	}

	// Normally Bob doesn't have dataPVSS but we are
	// using it only for PVSS parameters for simplicity
	recSecret, err := pvss.RecoverSecret(dataPVSS.Suite, writeData.G, validKeys, validEncShares, validDecShares, dataPVSS.Threshold, dataPVSS.NumTrustee)
	require.NoError(t, err)
	require.NotNil(t, recSecret)

	recMesg, err := DecryptMessage(recSecret, encMesg)
	require.NoError(t, err)
	require.Equal(t, bytes.Compare(recMesg, mesg), 0)
}

func TestMain(m *testing.M) {
	log.MainTest(m)
}
