#Inciando el contrato .
docker exec -e "CORE_PEER_GOSSIP_USELEADERELECTION=true" -it cli peer chaincode invoke -o orderer.etsiit.ugr.es:7050 -C channelatc -n grades -c '{"function":"initLedger","Args":[]}'
#Query de prueba
docker exec -e "CORE_PEER_GOSSIP_USELEADERELECTION=true" -it cli peer chaincode invoke -o orderer.etsiit.ugr.es:7050 -C channelatc -n grades -c '{"Args":["queryAllGrades"]}'


