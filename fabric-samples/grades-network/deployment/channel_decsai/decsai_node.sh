#Instalación de contrato.
docker exec -it cli peer chaincode install -n grades -p github.com/chaincode/grades -v 1.0
echo "creando canal"
docker exec -e "CORE_PEER_MSPCONFIGPATH=/var/hyperledger/users/Admin@decsai.ugr.es/msp" peer0.decsai.ugr.es peer channel create -o  orderer.etsiit.ugr.es:7050 -c channeldecsai -f /var/hyperledger/configs/channeldecsai.tx
echo "Uniendose al canal."
docker exec -e "CORE_PEER_MSPCONFIGPATH=/var/hyperledger/users/Admin@decsai.ugr.es/msp" peer0.decsai.ugr.es peer channel join -b channeldecsai.block
echo "instanciando el canal."
docker exec -it cli peer chaincode instantiate -o orderer.etsiit.ugr.es:7050 -C channeldecsai -n grades github.com/chaincode/grades -v 1.0 -c '{"Args": []}' -P "OR('decsaiMSP.member')"
