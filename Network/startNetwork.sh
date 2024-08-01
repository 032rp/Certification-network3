#!/bin/bash

echo "Register ca admin for each organisation"

docker-compose -f docker/docker-compose-ca.yaml up -d
sleep 3
sudo chmod -R 777 organizations/
echo "------------Register and enroll the users for each organization—-----------"
chmod +x registerEnroll.sh
./registerEnroll.sh
sleep 3

echo "—-------------Build the infrastructure—-----------------"

docker-compose -f docker/docker-compose-certnetwork.yaml up -d
sleep 3
echo "-------------Generate the genesis block—-------------------------------"

export FABRIC_CFG_PATH=./config
export CHANNEL_NAME=mychannel
configtxgen -profile CertChannelUsingRaft -outputBlock ./channel-artifacts/${CHANNEL_NAME}.block -channelID $CHANNEL_NAME
sleep 2
echo "------ Create the application channel------"

echo "Joining Orderer 1"
export ORDERER_CA=./organizations/ordererOrganizations/certnetwork.com/orderers/orderer1.certnetwork.com/msp/tlscacerts/tlsca.certnetwork.com-cert.pem
export ORDERER_ADMIN_TLS_SIGN_CERT=./organizations/ordererOrganizations/certnetwork.com/orderers/orderer1.certnetwork.com/tls/server.crt
export ORDERER_ADMIN_TLS_PRIVATE_KEY=./organizations/ordererOrganizations/certnetwork.com/orderers/orderer1.certnetwork.com/tls/server.key
osnadmin channel join --channelID $CHANNEL_NAME --config-block ./channel-artifacts/$CHANNEL_NAME.block -o localhost:7053 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
sleep 2
osnadmin channel list -o localhost:7053 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
sleep 2

echo "Joining Orderer 2"
export ORDERER_CA=./organizations/ordererOrganizations/certnetwork.com/orderers/orderer2.certnetwork.com/msp/tlscacerts/tlsca.certnetwork.com-cert.pem
export ORDERER_ADMIN_TLS_SIGN_CERT=./organizations/ordererOrganizations/certnetwork.com/orderers/orderer2.certnetwork.com/tls/server.crt
export ORDERER_ADMIN_TLS_PRIVATE_KEY=./organizations/ordererOrganizations/certnetwork.com/orderers/orderer2.certnetwork.com/tls/server.key
osnadmin channel join --channelID $CHANNEL_NAME --config-block ./channel-artifacts/$CHANNEL_NAME.block -o localhost:7055 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
sleep 2
osnadmin channel list -o localhost:7055 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
sleep 2

echo "Joining Orderer 3"
export ORDERER_CA=./organizations/ordererOrganizations/certnetwork.com/orderers/orderer3.certnetwork.com/msp/tlscacerts/tlsca.certnetwork.com-cert.pem
export ORDERER_ADMIN_TLS_SIGN_CERT=./organizations/ordererOrganizations/certnetwork.com/orderers/orderer3.certnetwork.com/tls/server.crt
export ORDERER_ADMIN_TLS_PRIVATE_KEY=./organizations/ordererOrganizations/certnetwork.com/orderers/orderer3.certnetwork.com/tls/server.key
osnadmin channel join --channelID $CHANNEL_NAME --config-block ./channel-artifacts/$CHANNEL_NAME.block -o localhost:7057 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
sleep 2
osnadmin channel list -o localhost:7057 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
sleep 2

export FABRIC_CFG_PATH=./peercfg
export CHANNEL_NAME=mychannel
export CORE_PEER_LOCALMSPID=iitMSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/iit.certnetwork.com/peers/peer0.iit.certnetwork.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/iit.certnetwork.com/users/Admin@iit.certnetwork.com/msp
export CORE_PEER_ADDRESS=localhost:8051
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/certnetwork.com/orderers/orderer1.certnetwork.com/msp/tlscacerts/tlsca.certnetwork.com-cert.pem
export IIT_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/iit.certnetwork.com/peers/peer0.iit.certnetwork.com/tls/ca.crt
export MHRD_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/mhrd.certnetwork.com/peers/peer0.mhrd.certnetwork.com/tls/ca.crt
export NPCL_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/npcl.certnetwork.com/peers/peer0.npcl.certnetwork.com/tls/ca.crt
sleep 2

echo "—---------------Join peer0 iit to the channel—-------------"
peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block
sleep 3
echo "—---------------Join peer1 iit to the channel—-------------"
export CORE_PEER_ADDRESS=localhost:9051
peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block
sleep 3

echo "peerchannel list"
peer channel list

echo "iit Anchor peer update"

peer channel fetch config channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer1.certnetwork.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
sleep 2
cd channel-artifacts
configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq '.data.data[0].payload.data.config' config_block.json > config.json
cp config.json config_copy.json
jq '.channel_group.groups.Application.groups.iitMSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.iit.certnetwork.com","port": 8051}]},"version": "0"}}' config_copy.json > modified_config.json
configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id ${CHANNEL_NAME} --original config.pb --updated modified_config.pb --output config_update.pb
configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb
cd ..
peer channel update -f channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer1.certnetwork.com --tls --cafile $ORDERER_CA
sleep 3

echo "—---------------package chaincode—-------------"

peer lifecycle chaincode package certnet.tar.gz --path ${PWD}/../chaincode/ --lang golang --label certnet_1.0
sleep 1

echo "—---------------install chaincode in IIT peer—-------------"

peer lifecycle chaincode install certnet.tar.gz
sleep 5
export CORE_PEER_ADDRESS=localhost:8051
peer lifecycle chaincode install certnet.tar.gz
sleep 5
peer lifecycle chaincode queryinstalled
sleep 1

export CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid certnet.tar.gz)

echo "—---------------Approve chaincode in IIT peer—-------------"

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer1.certnetwork.com --channelID $CHANNEL_NAME --name CertNetwork --version 1.0 --collections-config ../chaincode/collection.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
sleep 2



export FABRIC_CFG_PATH=./peercfg
export CHANNEL_NAME=mychannel
export CORE_PEER_LOCALMSPID=mhrdMSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/mhrd.certnetwork.com/peers/peer0.mhrd.certnetwork.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/mhrd.certnetwork.com/users/Admin@mhrd.certnetwork.com/msp
export CORE_PEER_ADDRESS=localhost:10051
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/certnetwork.com/orderers/orderer1.certnetwork.com/msp/tlscacerts/tlsca.certnetwork.com-cert.pem
export IIT_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/iit.certnetwork.com/peers/peer0.iit.certnetwork.com/tls/ca.crt
export MHRD_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/mhrd.certnetwork.com/peers/peer0.mhrd.certnetwork.com/tls/ca.crt
export NPCL_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/npcl.certnetwork.com/peers/peer0.npcl.certnetwork.com/tls/ca.crt
sleep 2

echo "—---------------Join peer0 mhrd to the channel—-------------"
peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block
sleep 3
echo "—---------------Join peer1 mhrd to the channel—-------------"
export CORE_PEER_ADDRESS=localhost:11051
peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block
sleep 3
echo "peerchannel list"
peer channel list

echo "mhrd Anchor peer update"
peer channel fetch config channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer1.certnetwork.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
sleep 2
cd channel-artifacts
configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq '.data.data[0].payload.data.config' config_block.json > config.json
cp config.json config_copy.json
jq '.channel_group.groups.Application.groups.mhrdMSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.mhrd.certnetwork.com","port": 10051}]},"version": "0"}}' config_copy.json > modified_config.json
configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id $CHANNEL_NAME --original config.pb --updated modified_config.pb --output config_update.pb
configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb
cd ..
peer channel update -f channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer1.certnetwork.com --tls --cafile $ORDERER_CA
sleep 3

echo "—---------------install chaincode in MHRD peer—-------------"

peer lifecycle chaincode install certnet.tar.gz
sleep 5
export CORE_PEER_ADDRESS=localhost:10051
peer lifecycle chaincode install certnet.tar.gz
sleep 5
peer lifecycle chaincode queryinstalled
sleep 1

export CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid certnet.tar.gz)

echo "—---------------Approve chaincode in MHRD peer—-------------"

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer1.certnetwork.com --channelID $CHANNEL_NAME --name CertNetwork --version 1.0 --collections-config ../chaincode/collection.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
sleep 2

export FABRIC_CFG_PATH=./peercfg
export CHANNEL_NAME=mychannel
export CORE_PEER_LOCALMSPID=npclMSP
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/npcl.certnetwork.com/peers/peer0.npcl.certnetwork.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/npcl.certnetwork.com/users/Admin@npcl.certnetwork.com/msp
export CORE_PEER_ADDRESS=localhost:12051
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/certnetwork.com/orderers/orderer1.certnetwork.com/msp/tlscacerts/tlsca.certnetwork.com-cert.pem
export IIT_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/iit.certnetwork.com/peers/peer0.iit.certnetwork.com/tls/ca.crt
export MHRD_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/mhrd.certnetwork.com/peers/peer0.mhrd.certnetwork.com/tls/ca.crt
export NPCL_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/npcl.certnetwork.com/peers/peer0.npcl.certnetwork.com/tls/ca.crt
sleep 2

echo "—---------------Join peer0 npcl to the channel—-------------"
peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block
sleep 3
echo "—---------------Join peer1 npcl to the channel—-------------"
export CORE_PEER_ADDRESS=localhost:13051
peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block
sleep 3
echo "peer channel list"
peer channel list


echo "npcl Anchor peer update"

peer channel fetch config channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer1.certnetwork.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
sleep 2
cd channel-artifacts
configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq '.data.data[0].payload.data.config' config_block.json > config.json
cp config.json config_copy.json
jq '.channel_group.groups.Application.groups.npclMSP.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.npcl.certnetwork.com","port": 12051}]},"version": "0"}}' config_copy.json > modified_config.json
configtxlator proto_encode --input config.json --type common.Config --output config.pb
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id ${CHANNEL_NAME} --original config.pb --updated modified_config.pb --output config_update.pb
configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb
cd ..
peer channel update -f channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer1.certnetwork.com --tls --cafile $ORDERER_CA
sleep 3
peer channel getinfo -c $CHANNEL_NAME

echo "—---------------install chaincode in NPCL peer—-------------"

peer lifecycle chaincode install certnet.tar.gz
sleep 3
export CORE_PEER_ADDRESS=localhost:12051
peer lifecycle chaincode install certnet.tar.gz
sleep 5
peer lifecycle chaincode queryinstalled
sleep 1

export CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid certnet.tar.gz)

echo "—---------------Approve chaincode in NPCL peer—-------------"

peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer1.certnetwork.com --channelID $CHANNEL_NAME --name CertNetwork --version 1.0 --collections-config ../chaincode/collection.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
sleep 2

peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_NAME --name CertNetwork --version 1.0 --sequence 1 --collections-config ../chaincode/collection.json --tls --cafile $ORDERER_CA --output json

peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer1.certnetwork.com --channelID $CHANNEL_NAME --name CertNetwork  --version 1.0 --sequence 1 --collections-config ../chaincode/collection.json --tls --cafile $ORDERER_CA --peerAddresses localhost:8051 --tlsRootCertFiles $IIT_PEER_TLSROOTCERT --peerAddresses localhost:10051 --tlsRootCertFiles $MHRD_PEER_TLSROOTCERT --peerAddresses localhost:12051 --tlsRootCertFiles $NPCL_PEER_TLSROOTCERT
sleep 2
peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name CertNetwork  --cafile $ORDERER_CA

