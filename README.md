#                               Certification Verification using Blockchain
###                                 Certificate Network Setup Guide
##### This guide provides step-by-step instructions for setting up a Hyperledger Fabric network for managing an certificate Network. The network architecture includes three organizations: IIT, MHRD, and NPCL.

### Prerequisites
##### Docker installed on your system
##### Docker Compose installed on your system
##### Hyperledger Fabric binaries and Docker images

### Setup Instructions

##### ***Open a command terminal with in Network folder, let's call this terminal as host terminal

```
cd Network/

```

##  **************** Host terminal ********************

### ------------Register the CA admin for each Organization—----------------

##### ***Build the docker-compose-ca.yaml in the docker folder

```
docker-compose -f docker/docker-compose-ca.yaml up -d

```

```
sudo chmod -R 777 organizations/
```
### ------------Register and enroll the users for each organization—-----------

##### ***Build the registerEnroll.sh script file


```
chmod +x registerEnroll.sh

./registerEnroll.sh
```


### —-------------Build the infrastructure—-----------------

##### ***Build the docker-compose-certnetwork.yaml in the docker folder

```
docker-compose -f docker/docker-compose-certnetwork.yaml up -d
```


### -------------Generate the genesis block—-------------------------------

##### ***Build the configtx.yaml file in the config folder

```
export FABRIC_CFG_PATH=./config

export CHANNEL_NAME=mychannel
```

```
configtxgen -profile CertChannelUsingRaft -outputBlock ./channel-artifacts/${CHANNEL_NAME}.block -channelID $CHANNEL_NAME
```


### ------ Create the Application channel------
###                                    -------Join ORDERER1-----------
```
export ORDERER_CA=./organizations/ordererOrganizations/certnetwork.com/orderers/orderer1.certnetwork.com/msp/tlscacerts/tlsca.certnetwork.com-cert.pem
export ORDERER_ADMIN_TLS_SIGN_CERT=./organizations/ordererOrganizations/certnetwork.com/orderers/orderer1.certnetwork.com/tls/server.crt
export ORDERER_ADMIN_TLS_PRIVATE_KEY=./organizations/ordererOrganizations/certnetwork.com/orderers/orderer1.certnetwork.com/tls/server.key

```

```
osnadmin channel join --channelID $CHANNEL_NAME --config-block ./channel-artifacts/$CHANNEL_NAME.block -o localhost:7053 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
```


```
osnadmin channel list -o localhost:7053 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
```
###                                 -------Join ORDERER2-----------
```
export ORDERER_CA=./organizations/ordererOrganizations/certnetwork.com/orderers/orderer2.certnetwork.com/msp/tlscacerts/tlsca.certnetwork.com-cert.pem
export ORDERER_ADMIN_TLS_SIGN_CERT=./organizations/ordererOrganizations/certnetwork.com/orderers/orderer2.certnetwork.com/tls/server.crt
export ORDERER_ADMIN_TLS_PRIVATE_KEY=./organizations/ordererOrganizations/certnetwork.com/orderers/orderer2.certnetwork.com/tls/server.key
osnadmin channel join --channelID $CHANNEL_NAME --config-block ./channel-artifacts/$CHANNEL_NAME.block -o localhost:7055 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
osnadmin channel list -o localhost:7055 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
```
###                                -------Join ORDERER3-----------
```
export ORDERER_CA=./organizations/ordererOrganizations/certnetwork.com/orderers/orderer3.certnetwork.com/msp/tlscacerts/tlsca.certnetwork.com-cert.pem
export ORDERER_ADMIN_TLS_SIGN_CERT=./organizations/ordererOrganizations/certnetwork.com/orderers/orderer3.certnetwork.com/tls/server.crt
export ORDERER_ADMIN_TLS_PRIVATE_KEY=./organizations/ordererOrganizations/certnetwork.com/orderers/orderer3.certnetwork.com/tls/server.key
osnadmin channel join --channelID $CHANNEL_NAME --config-block ./channel-artifacts/$CHANNEL_NAME.block -o localhost:7057 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
osnadmin channel list -o localhost:7057 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
```
##### ***Open another terminal with in Network folder, let's call this terminal as IIT terminal.

##  **************** IIT terminal ********************

##### ***Build the core.yaml in peercfg folder
```
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
```
### ---------------Join peer0-iit to the channel—-------------
```
peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block
```
### —---------------Join peer1-iit to the channel—-------------
```
export CORE_PEER_ADDRESS=localhost:9051
peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block
peer channel list
```


##### ***Open another terminal with in Network folder, let's call this terminal as MHRD terminal.

##  **************** MHRD terminal *****************

```
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
```
### —---------------Join peer0-mhrd to the channel—-------------
```
peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block
```
### —---------------Join peer1-mhrd to the channel—-------------
export CORE_PEER_ADDRESS=localhost:11051
peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block
```
peer channel list
```

##### ***Open another terminal with in Network folder, let's call this terminal as NPCL terminal.

##  **************** NPCL terminal ******************

```
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
```
### —---------------Join peer0-npcl to the channel—-------------
```
peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block
```
### —---------------Join peer1-npcl to the channel—-------------
```
export CORE_PEER_ADDRESS=localhost:13051
peer channel join -b ./channel-artifacts/$CHANNEL_NAME.block
peer channel list
```

### —-------------------- Anchor peer update —---------------------

##  **************** IIT terminal ******************

```
peer channel fetch config channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer1.certnetwork.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
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
```
##  **************** MHRD terminal ******************

```
peer channel fetch config channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer1.certnetwork.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
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

```
##  **************** NPCL terminal ******************

```
peer channel fetch config channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer1.certnetwork.com -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
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
```
```
peer channel getinfo -c $CHANNEL_NAME
```


### —-----------------Chaincode lifecycle—-------------------

##### ***Build the chaincode***

##### ***Make sure that Chaincode is available in chaincode folder which is at the same level of Network. 

##  **************** IIT Terminal ******************


—---------------package chaincode—-------------"
```
peer lifecycle chaincode package certnet.tar.gz --path ${PWD}/../chaincode/ --lang golang --label certnet_1.0
```

### —---------------Install chaincode in IIT peer1—-------------
```
peer lifecycle chaincode install certnet.tar.gz
```
### —---------------Install chaincode in IIT peer0—-------------
```
export CORE_PEER_ADDRESS=localhost:8051
peer lifecycle chaincode install certnet.tar.gz
```
```
peer lifecycle chaincode queryinstalled
export CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid certnet.tar.gz)
```

###  **************** MHRD Terminal ******************

### —---------------Install chaincode in MHRD peer1—-------------
```
peer lifecycle chaincode install certnet.tar.gz
```
### —---------------Install chaincode in MHRD peer0—-------------
```
export CORE_PEER_ADDRESS=localhost:10051
peer lifecycle chaincode install certnet.tar.gz
```
```
peer lifecycle chaincode queryinstalled
export CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid certnet.tar.gz)
```

##  **************** NPCL Terminal ******************
### —---------------Install chaincode in NPCL peer1—-------------
```
peer lifecycle chaincode install certnet.tar.gz
```
### —---------------Install chaincode in NPCL peer0—-------------
```
export CORE_PEER_ADDRESS=localhost:12051
peer lifecycle chaincode install certnet.tar.gz
```
```
peer lifecycle chaincode queryinstalled
export CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid certnet.tar.gz)
```
### Approve Chaincode from each Organizations
##  **************** IIT Terminal ******************

```
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer1.certnetwork.com --channelID $CHANNEL_NAME --name CertNetwork --version 1.0 --collections-config ../chaincode/collection.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
```

##  **************** MHRD Terminal ******************

```
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer1.certnetwork.com --channelID $CHANNEL_NAME --name CertNetwork --version 1.0 --collections-config ../chaincode/collection.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
```

##  **************** NPCL Terminal ******************
```
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer1.certnetwork.com --channelID $CHANNEL_NAME --name CertNetwork --version 1.0 --collections-config ../chaincode/collection.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
```

###     Commit Chaincode from any Terminal
##  **************** IIT Terminal ******************

```
peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_NAME --name CertNetwork --version 1.0 --sequence 1 --collections-config ../chaincode/collection.json --tls --cafile $ORDERER_CA --output json
```
```
peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer1.certnetwork.com --channelID $CHANNEL_NAME --name CertNetwork  --version 1.0 --sequence 1 --collections-config ../chaincode/collection.json --tls --cafile $ORDERER_CA --peerAddresses localhost:8051 --tlsRootCertFiles $IIT_PEER_TLSROOTCERT --peerAddresses localhost:10051 --tlsRootCertFiles $MHRD_PEER_TLSROOTCERT --peerAddresses localhost:12051 --tlsRootCertFiles $NPCL_PEER_TLSROOTCERT
```
```
peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name CertNetwork  --cafile $ORDERER_CA
```
### Invoke(PrivateData) the Chaincode from IIT Terminal
```
export NAME=$(echo -n "Ram" | base64 | tr -d \\n)
export EMAIL=$(echo -n "rp@gmail.com" | base64 | tr -d \\n)
export DOB=$(echo -n "06/07/1995" | base64 | tr -d \\n)
export COLLEGE=$(echo -n "ShyamalDas" | base64 | tr -d \\n)
```
```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer1.certnetwork.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n CertNetwork --peerAddresses localhost:8051 --tlsRootCertFiles $IIT_PEER_TLSROOTCERT --peerAddresses localhost:10051 --tlsRootCertFiles $MHRD_PEER_TLSROOTCERT --peerAddresses localhost:12051 --tlsRootCertFiles $NPCL_PEER_TLSROOTCERT -c '{"Args":["StudentContract:RegisterStudent","ST0001"]}' --transient "{\"name\":\"$NAME\",\"email\":\"$EMAIL\",\"dateofbirth\":\"$DOB\",\"college\":\"$COLLEGE\"}"
```

### To Query the chaincode 
peer chaincode query -C $CHANNEL_NAME -n CertNetwork -c '{"Args":["StudentrContract:GetStudent","ST0001"]}'


### --------- Stop the Certificate-network --------------

##  **************** Host terminal ******************

```
docker-compose -f docker/docker-compose-certnetwork.yaml down
```
```
docker-compose -f docker/docker-compose-ca.yaml down
```
```
docker rm -f $(docker ps -a | awk '($2 ~ /dev-peer.*/) {print $1}')
```
```
docker volume rm $(docker volume ls -q)
```
```
sudo rm -rf organizations/
```
```
sudo rm -rf channel-artifacts/
```
```
sudo rm certnet.tar.gz
```
```
docker ps -a
```

##### ***if there still exists the containers then execute the following commands.

```
docker rm $(docker container ls -q) --force
```
```
docker container prune
```
```
docker system prune
```
```
docker volume prune
```
```
docker network prune
```


### Run using startNetwork.sh script

##### ***Build startNetwork.sh script file

```
chmod +x startNetwork.sh
```

```
./startNetwork.sh
```


##### ***To submit transaction as iitMSP***
```
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

```

### Invoke(PrivateData) the Chaincode 
```
export NAME=$(echo -n "Ram" | base64 | tr -d \\n)
export EMAIL=$(echo -n "rp@gmail.com" | base64 | tr -d \\n)
export DOB=$(echo -n "06/07/1995" | base64 | tr -d \\n)
export COLLEGE=$(echo -n "ShyamalDas" | base64 | tr -d \\n)
```
```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer1.certnetwork.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n CertNetwork --peerAddresses localhost:8051 --tlsRootCertFiles $IIT_PEER_TLSROOTCERT --peerAddresses localhost:10051 --tlsRootCertFiles $MHRD_PEER_TLSROOTCERT --peerAddresses localhost:12051 --tlsRootCertFiles $NPCL_PEER_TLSROOTCERT -c '{"Args":["StudentContract:RegisterStudent","ST0001"]}' --transient "{\"name\":\"$NAME\",\"email\":\"$EMAIL\",\"dateofbirth\":\"$DOB\",\"college\":\"$COLLEGE\"}"
```

### To Query the chaincode 
peer chaincode query -C $CHANNEL_NAME -n CertNetwork -c '{"Args":["StudentrContract:GetStudent","ST0001"]}'

### To stop the network using script file

##### ***Build stopNetwork.sh script file
```
chmod +x stopNetwork.sh
```
```
./stopNetwork.sh
```

### To stop the network without clearing the data using docker

```
docker-compose -f ./docker/docker-compose-ca.yaml -f ./docker/docker-compose-certnetwork.yaml stop
```
### To start the network with existing docker containers
```
docker-compose -f ./docker/docker-compose-ca.yaml -f ./docker/docker-compose-certnetwork.yaml start
```
