# Instalar el contrato.
docker exec -it cli peer chaincode install -n audit -p github.com/chaincode/audit -v 1.0

#Crea el canal
docker exec -e "CORE_PEER_MSPCONFIGPATH=/var/hyperledger/users/Admin@atc.ugr.es/msp" peer0.atc.ugr.es peer channel create -o orderer.etsiit.ugr.es:7050 -c channelall -f /var/hyperledger/configs/channelall.tx

#Unir el nodo atc al canal
docker exec -e "CORE_PEER_MSPCONFIGPATH=/var/hyperledger/users/Admin@atc.ugr.es/msp" peer0.atc.ugr.es peer channel join -b channelall.block

#Copia el blocke del canal a la carpeta
docker cp peer0.atc.ugr.es:channelall.block .

#Comando para copiar el bloque genesis a los otros nodos.
scp channelall.block nodo@192.168.56.222:/home/nodo/Fabric-Grades/fabric-samples/grades-network/deployment/channel_all/channelall.block
scp channelall.block nodo@192.168.56.223:/home/nodo/Fabric-Grades/fabric-samples/grades-network/deployment/channel_all/channelall.block
