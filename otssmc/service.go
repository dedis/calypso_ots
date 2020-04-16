package otssmc

import (
	"gopkg.in/dedis/onet.v1"
	"gopkg.in/dedis/onet.v1/log"
	"gopkg.in/dedis/onet.v1/network"
)

const ServiceName = "OTSSMCService"

const (
	// ErrorParse indicates an error while parsing the protobuf-file.
	ErrorParse = iota + 4000
)

type OTSSMCService struct {
	*onet.ServiceProcessor
}

func init() {
	onet.RegisterNewService(ServiceName, newOTSSMCService)
	network.RegisterMessage(&DecryptRequest{})
	network.RegisterMessage(&DecryptReply{})
	// network.RegisterMessage(&util.OTSDecryptReqData{})
	// network.RegisterMessage(&util.DecryptedShare{})
}

func (s *OTSSMCService) DecryptRequest(req *DecryptRequest) (*DecryptReply, onet.ClientError) {
	log.Lvl3("DecryptRequest received in service")
	// Tree with depth = 1
	childCount := len(req.Roster.List) - 1
	log.Lvl3("Number of childs:", childCount)
	tree := req.Roster.GenerateNaryTreeWithRoot(childCount, s.ServerIdentity())
	if tree == nil {
		return nil, onet.NewClientErrorCode(ErrorParse, "couldn't create tree")
	}

	pi, err := s.CreateProtocol(ProtocolName, tree)
	if err != nil {
		return nil, onet.NewClientError(err)
	}

	otsDec := pi.(*OTSDecrypt)
	otsDec.DecReqData = req.Data
	otsDec.Signature = req.Signature
	otsDec.RootIndex = req.RootIndex
	err = pi.Start()
	if err != nil {
		return nil, onet.NewClientError(err)
	}

	reply := &DecryptReply{
		DecShares: <-pi.(*OTSDecrypt).DecShares,
	}
	return reply, nil
}

func (s *OTSSMCService) NewProtocol(tn *onet.TreeNodeInstance, conf *onet.GenericConfig) (onet.ProtocolInstance, error) {
	log.Lvl3("OTSDecrypt Service received New Protocol event")
	pi, err := NewProtocol(tn)
	return pi, err
}

func newOTSSMCService(c *onet.Context) onet.Service {
	s := &OTSSMCService{
		ServiceProcessor: onet.NewServiceProcessor(c),
	}
	err := s.RegisterHandler(s.DecryptRequest)
	log.Lvl3("OTSSMC Service registered")
	if err != nil {
		log.ErrFatal(err, "Couldn't register message:")
	}
	return s
}
