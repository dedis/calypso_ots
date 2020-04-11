module github.com/calypso-demo/ots

go 1.13

replace github.com/dedis/protobuf => ./imported/protobuf

replace gopkg.in/dedis/onet.v1 => ./imported/onet

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/daviddengcn/go-colortext v1.0.0 // indirect
	github.com/dedis/protobuf v0.0.0-00010101000000-000000000000 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/montanaflynn/stats v0.6.3 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/shirou/gopsutil v2.20.3+incompatible // indirect
	go.etcd.io/bbolt v1.3.4 // indirect
	golang.org/x/crypto v0.0.0-20200406173513-056763e48d71 // indirect
	gopkg.in/dedis/cothority.v1 v1.0.0-20180112132810-9daa49171eb7
	gopkg.in/dedis/crypto.v0 v0.0.0-20170824083343-8f53a63e87fd
	gopkg.in/dedis/onet.v1 v1.0.0-00010101000000-000000000000
	gopkg.in/satori/go.uuid.v1 v1.2.0 // indirect
	gopkg.in/tylerb/graceful.v1 v1.2.15 // indirect
	gopkg.in/urfave/cli.v1 v1.20.0 // indirect
)
