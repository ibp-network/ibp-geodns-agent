module github.com/ibp-network/ibp-geodns-agent

go 1.24.2

require (
	github.com/ibp-network/ibp-geodns-libs v0.2.0
	github.com/nats-io/nats.go v1.48.0
)

require (
	github.com/klauspost/compress v1.18.3 // indirect
	github.com/nats-io/nkeys v0.4.12 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.47.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
)

replace github.com/ibp-network/ibp-geodns-libs => ../ibp-geodns-libs
