package main

// Config represents the configuration for a role.
type Config struct {
	CertPath     string `json:"certPath"`
	KeyDirectory string `json:"keyPath"`
	TLSCertPath  string `json:"tlsCertPath"`
	PeerEndpoint string `json:"peerEndpoint"`
	GatewayPeer  string `json:"gatewayPeer"`
	MSPID        string `json:"mspID"`
}

// Create a Profile map
var profile = map[string]Config{

	"iit": {
		CertPath:     "../Network/organizations/peerOrganizations/iit.certnetwork.com/users/User1@iit.certnetwork.com/msp/signcerts/cert.pem",
		KeyDirectory: "../Network/organizations/peerOrganizations/iit.certnetwork.com/users/User1@iit.certnetwork.com/msp/keystore/",
		TLSCertPath:  "../Network/organizations/peerOrganizations/iit.certnetwork.com/peers/peer0.iit.certnetwork.com/tls/ca.crt",
		PeerEndpoint: "localhost:8051",
		GatewayPeer:  "peer0.iit.certnetwork.com",
		MSPID:        "iitMSP",
	},

	"mhrd": {
		CertPath:     "../Network/organizations/peerOrganizations/mhrd.certnetwork.com/users/User1@mhrd.certnetwork.com/msp/signcerts/cert.pem",
		KeyDirectory: "../Network/organizations/peerOrganizations/mhrd.certnetwork.com/users/User1@mhrd.certnetwork.com/msp/keystore/",
		TLSCertPath:  "../Network/organizations/peerOrganizations/mhrd.certnetwork.com/peers/peer0.mhrd.certnetwork.com/tls/ca.crt",
		PeerEndpoint: "localhost:10051",
		GatewayPeer:  "peer0.mhrd.certnetwork.com",
		MSPID:        "mhrdMSP",
	},

	"npcl": {
		CertPath:     "../Network/organizations/peerOrganizations/npcl.certnetwork.com/users/User1@npcl.certnetwork.com/msp/signcerts/cert.pem",
		KeyDirectory: "../Network/organizations/peerOrganizations/npcl.certnetwork.com/users/User1@npcl.certnetwork.com/msp/keystore/",
		TLSCertPath:  "../Network/organizations/peerOrganizations/npcl.certnetwork.com/peers/peer0.npcl.certnetwork.com/tls/ca.crt",
		PeerEndpoint: "localhost:12051",
		GatewayPeer:  "peer0.npcl.certnetwork.com",
		MSPID:        "npclMSP",
	},

	"iit2": {
		CertPath:     "../Network/organizations/peerOrganizations/iit.certnetwork.com/users/User2@iit.certnetwork.com/msp/signcerts/cert.pem",
		KeyDirectory: "../Network/organizations/peerOrganizations/iit.certnetwork.com/users/User2@iit.certnetwork.com/msp/keystore/",
		TLSCertPath:  "../Network/organizations/peerOrganizations/iit.certnetwork.com/peers/peer0.iit.certnetwork.com/tls/ca.crt",
		PeerEndpoint: "localhost:8051",
		GatewayPeer:  "peer0.iit.certnetwork.com",
		MSPID:        "iitMSP",
	},
}
