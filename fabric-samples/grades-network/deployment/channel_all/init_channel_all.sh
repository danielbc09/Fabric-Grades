# Instalar el contrato.

docker exec -it cli peer chaincode instantiate -o orderer.etsiit.ugr.es:7050 -C channelall -n audit github.com/chaincode/audit -v 1.0 -c '{"Args": []}' -P "OR('atcMSP.member', 'decsaiMSP.member','auditMSP.member')"

docker exec -e "CORE_PEER_GOSSIP_USELEADERELECTION=true" -it cli peer chaincode invoke -o orderer.etsiit.ugr.es:7050 -C channelall -n audit -c '{"function":"initLedger","Args":[]}'

docker exec -e "CORE_PEER_GOSSIP_USELEADERELECTION=true" -it cli peer chaincode invoke -o orderer.etsiit.ugr.es:7050 -C channelall -n audit -c '{"function":"queryAllAuditGrades","Args":[]}'
