# Instalar el contrato.
docker exec -it cli peer chaincode install -n audit -p github.com/chaincode/audit -v 1.0
#Coopiar el bloque genesis al contenedor
docker cp channelall.block peer0.audit.ugr.es:/channelall.block
#Unirse al canal 
docker exec -e "CORE_PEER_MSPCONFIGPATH=/var/hyperledger/users/Admin@audit.ugr.es/msp" peer0.audit.ugr.es peer channel join -b channelall.block