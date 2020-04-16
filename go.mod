module github.com/calypso-demo/ots

go 1.13

replace github.com/dedis/protobuf => ./imported/protobuf

replace gopkg.in/dedis/onet.v1 => ./imported/onet

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/bford/golang-x-crypto v0.0.0-20160518072526-27db609c9d03 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.3.0
	go.dedis.ch/kyber/v3 v3.0.0-pre4
	go.etcd.io/bbolt v1.3.4 // indirect
	golang.org/x/crypto v0.0.0-20200406173513-056763e48d71 // indirect
	gopkg.in/dedis/cothority.v1 v1.0.0-20180112132810-9daa49171eb7
	gopkg.in/dedis/crypto.v0 v0.0.0-20170824083343-8f53a63e87fd
	gopkg.in/dedis/onet.v1 v1.0.0-00010101000000-000000000000
	gopkg.in/satori/go.uuid.v1 v1.2.0 // indirect
)
