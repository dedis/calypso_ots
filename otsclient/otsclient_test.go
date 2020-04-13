package otsclient

import (
	"testing"

	"github.com/stretchr/testify/require"
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

	//scPubKeys := roster.Publics()
	//dataPVSS := util.DataPVSS{
	//Suite:        ed25519.NewAES128SHA256Ed25519(false),
	//SCPublicKeys: scPubKeys,
	//NumTrustee:   numServers,
	//}

	// Writer's pk/sk pair
	//wrPrivKey := dataPVSS.Suite.Scalar().Pick(random.Stream)
	//wrPubKey := dataPVSS.Suite.Point().Mul(nil, wrPrivKey)
	// Reader's pk/sk pair
	//privKey := dataPVSS.Suite.Scalar().Pick(random.Stream)
	//pubKey := dataPVSS.Suite.Point().Mul(nil, privKey)
	//err = SetupPVSS(&dataPVSS, pubKey)
	//require.NoError(t, err)

	//mesgSize := 1024 * 1024
	//mesg := make([]byte, mesgSize)
	//for i := 0; i < mesgSize; i++ {
	//mesg[i] = 'w'
	//}
	////encMesg, hashEnc, err := EncryptMessage(&dataPVSS, mesg)
	//_, hashEnc, err := EncryptMessage(&dataPVSS, mesg)
	//require.NoError(t, err)

	//// Creating write transaction
	//writeSB, err := CreateWriteTxn(scurl, &dataPVSS, hashEnc, pubKey, wrPrivKey)
	//require.NoError(t, err)

	//// Bob gets it from Alice
	//writeID := writeSB.Hash
	//// Get write transaction from skipchain
	////writeSB, writeTxnData, sig, err := GetWriteTxnSB(scurl, writeID)
	//writeSB, _, _, err = GetWriteTxnSB(scurl, writeID)
	//require.NoError(t, err)
}

func TestMain(m *testing.M) {
	log.MainTest(m)
}
