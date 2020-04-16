package main

import (
	"bytes"
	"errors"

	"github.com/BurntSushi/toml"
	"github.com/calypso-demo/ots"
	ocs "github.com/calypso-demo/ots/onchain-secrets"
	"github.com/calypso-demo/ots/otssmc"

	"github.com/calypso-demo/ots/otsclient"
	"gopkg.in/dedis/crypto.v0/abstract"
	"gopkg.in/dedis/crypto.v0/ed25519"
	"gopkg.in/dedis/crypto.v0/random"
	"gopkg.in/dedis/crypto.v0/share/pvss"
	"gopkg.in/dedis/onet.v1"
	"gopkg.in/dedis/onet.v1/log"
	"gopkg.in/dedis/onet.v1/network"
)

func init() {
	onet.SimulationRegister("OTSDemo", NewOTSSimulation)
}

type OTSSimulation struct {
	onet.SimulationBFTree
}

func NewOTSSimulation(config string) (onet.Simulation, error) {
	otss := &OTSSimulation{}
	_, err := toml.Decode(config, otss)
	if err != nil {
		return nil, err
	}
	return otss, nil
}

func (otss *OTSSimulation) Setup(dir string, hosts []string) (*onet.SimulationConfig, error) {
	sc := &onet.SimulationConfig{}
	otss.CreateRoster(sc, hosts, 2000)
	err := otss.CreateTree(sc)
	if err != nil {
		return nil, err
	}
	return sc, nil
}

func (otss *OTSSimulation) Node(config *onet.SimulationConfig) error {
	return otss.SimulationBFTree.Node(config)
}

func (otss *OTSSimulation) Run(config *onet.SimulationConfig) error {
	scPubKeys := config.Roster.Publics()
	numTrustee := config.Tree.Size()
	mesgSize := 1024 * 1024
	mesg := make([]byte, mesgSize)
	for i := 0; i < mesgSize; i++ {
		mesg[i] = 'w'
	}
	scurl, err := otsclient.CreateSkipchain(config.Roster)
	if err != nil {
		return err
	}
	for round := 0; round < otss.Rounds; round++ {
		dataPVSS := otsclient.DataPVSS{
			Suite:        ed25519.NewAES128SHA256Ed25519(false),
			SCPublicKeys: scPubKeys,
			NumTrustee:   numTrustee,
		}

		wrPrivKey := dataPVSS.Suite.Scalar().Pick(random.Stream)
		wrPubKey := dataPVSS.Suite.Point().Mul(nil, wrPrivKey)
		// Reader's pk/sk pair
		privKey := dataPVSS.Suite.Scalar().Pick(random.Stream)
		pubKey := dataPVSS.Suite.Point().Mul(nil, privKey)

		err = otsclient.SetupPVSS(&dataPVSS, pubKey)
		if err != nil {
			return err
		}
		encMesg, hashEnc, err := otsclient.EncryptMessage(&dataPVSS, mesg)
		if err != nil {
			return err
		}
		writeSB, err := otsclient.WriteRequest(scurl, &dataPVSS, hashEnc, pubKey, wrPrivKey)
		if err != nil {
			return err
		}

		// Bob gets it from Alice
		writeID := writeSB.Hash
		// Get write transaction from skipchain
		writeSB, writeData, writeSig, err := otsclient.GetWriteSB(scurl, writeID)
		if err != nil {
			return err
		}
		sigVerErr := otsclient.VerifyWriteSignature(writeData, writeSig, wrPubKey)
		if sigVerErr != nil {
			return sigVerErr
		}
		validHash := otsclient.VerifyEncMesg(writeData, encMesg)
		if validHash != 0 {
			return errors.New("Cannot verify encrypted message")
		}
		readSB, err := otsclient.ReadRequest(scurl, writeID, privKey)
		if err != nil {
			return err
		}
		updWriteSB, err := otsclient.GetUpdatedWriteSB(scurl, writeID)
		if err != nil {
			return err
		}

		acPubKeys := readSB.Roster.Publics()
		readSBF := readSB.SkipBlockFix
		p, err := config.Overlay.CreateProtocol("otssmc", config.Tree, onet.NilServiceID)
		if err != nil {
			return err
		}

		// GetDecryptedShares call preparation
		idx := readSB.Index - updWriteSB.Index - 1
		if idx < 0 {
			return errors.New("Forward-link index is negative")
		}
		inclusionProof := updWriteSB.GetForward(idx)
		if inclusionProof == nil {
			return errors.New("Forward-link does not exist")
		}
		data := &otssmc.DecryptRequestData{
			WriteSBF:       updWriteSB.SkipBlockFix,
			ReadSBF:        readSBF,
			InclusionProof: inclusionProof,
			ACPublicKeys:   acPubKeys,
		}
		proto := p.(*otssmc.OTSDecrypt)
		proto.DecReqData = data
		proto.RootIndex = 0
		msg, err := network.Marshal(data)
		if err != nil {
			return err
		}
		sig, err := ots.SignMessage(msg, privKey)
		if err != nil {
			return err
		}
		proto.Signature = &sig
		go p.Start()
		reencShares := <-proto.DecShares
		tmpDecShares, err := otsclient.ElGamalDecrypt(reencShares, privKey)
		if err != nil {
			return err
		}
		size := len(tmpDecShares)
		decShares := make([]*pvss.PubVerShare, size)
		for i := 0; i < size; i++ {
			decShares[tmpDecShares[i].S.I] = tmpDecShares[i]
		}
		var validKeys []abstract.Point
		var validEncShares []*pvss.PubVerShare
		var validDecShares []*pvss.PubVerShare
		for i := 0; i < size; i++ {
			validKeys = append(validKeys, writeData.SCPublicKeys[i])
			validEncShares = append(validEncShares, writeData.EncShares[i])
			validDecShares = append(validDecShares, decShares[i])
		}
		recSecret, err := pvss.RecoverSecret(dataPVSS.Suite, writeData.G, validKeys, validEncShares, validDecShares, dataPVSS.Threshold, dataPVSS.NumTrustee)
		if err != nil {
			return err
		}
		recvMesg, err := otsclient.DecryptMessage(recSecret, encMesg)
		if err != nil {
			return err
		}
		log.Info("Recovered message?:", bytes.Compare(mesg, recvMesg) == 0)
	}
	return nil
}

func prepareDummyDP(scurl *ocs.SkipChainURL, scRoster *onet.Roster, pairCount int) error {
	scPubKeys := scRoster.Publics()
	numTrustee := len(scPubKeys)
	dp := otsclient.DataPVSS{
		Suite:        ed25519.NewAES128SHA256Ed25519(false),
		SCPublicKeys: scPubKeys,
		NumTrustee:   numTrustee,
	}
	return otsclient.AddDummyWRPairs(scurl, &dp, pairCount)
}
