#!/bin/bash
PROJECT_DIR=$PWD
MODE="$1"
function ayuda() {
  echo "Uso: "
  echo "  ops.sh "
  echo "    <mode> -'generar', 'crearcanales', 'remplazar', 'limpiar'"
  echo "      - 'generar' - Genera certificados"
  echo "      - 'crearcanales' - Crea los canales channelall channelatc channeldecsai"
  echo "      - 'remplazar' - Copia las claves privadas para las CA de los docker compose de cada nodo."
  echo "      - 'limpiar' - Se eliminan los contenedores para volver a instalar la red."
}


function generarCertificados(){

	echo "##########################################################"
	echo "Generando certificados utilizando la herramienta cryptogen"
	echo "##########################################################"
  
	if [ -d /crypto-config ]; then
		rm -rf /crypto-config
	fi

  if [ "$?" -ne 0 ]; then
    echo "Ha habido un problema con la herramienta rypto-config.yaml"
    exit 1
  fi

 ../bin/cryptogen generate --config=./crypto-config.yaml
  echo
}


function generarCanales(){
  
    if [ ! -d ./channel-artifacts ]; then
	    	mkdir channel-artifacts
    fi

  echo "###################"
	echo "Creando bloque genesis."
	echo "###################"
  
  export FABRIC_CFG_PATH=$PWD
	../bin/configtxgen -profile FourOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block
  
  echo "###################"
	echo "Generarndo canal channelAll"
	echo "###################"

  ../bin/configtxgen -profile ChannelAll -outputCreateChannelTx ./channel-artifacts/channelall.tx -channelID channelall   

  echo "##################################"
	echo "Generarndo canal departamento atc"
	echo "##################################"
  ../bin/configtxgen -profile Channelatc -outputCreateChannelTx ./channel-artifacts/channelatc.tx -channelID channelatc
 
  echo "####################################"
	echo "Generarndo canal departamento decsai"
	echo "#####################################"
  ../bin/configtxgen -profile Channeldecsai -outputCreateChannelTx ./channel-artifacts/channeldecsai.tx -channelID channeldecsai
}


function remplazarClavesPrivadas(){
  
	echo "######################################################"
	echo "Copiando las claves a los archivos docker-compose.yml"
	echo "######################################################"

  CURRENT_DIR=$PWD
  cd crypto-config/peerOrganizations/atc.ugr.es/ca/
  PRIV_KEY=$(ls *_sk)
  cd $CURRENT_DIR"/deployment"
  sed -i "s/CA_ATC_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose-node1.yml

  cd ../crypto-config/peerOrganizations/decsai.ugr.es/ca/
  PRIV_KEY=$(ls *_sk)
  cd $CURRENT_DIR"/deployment"
  sed -i "s/CA_DECSAI_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose-node2.yml

  cd ../crypto-config/peerOrganizations/audit.ugr.es/ca/
  PRIV_KEY=$(ls *_sk)
  cd $CURRENT_DIR"/deployment"
  sed -i "s/CA_AUDIT_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose-node3.yml
 
}

function limpiarRed(){
  echo "########################"
	echo "Eliminando contenedores."
	echo "########################"
  docker stop $(docker ps -aq)
  docker rm $(docker ps -aq)
  docker volume prune
  docker rmi $(docker images |grep 'net-peer')

  echo "Contenedores removidos correctamente."
}

#Comandos para creacion de certificados y canales.
if [ "${MODE}" == "generar" ]; then
  generarCertificados
elif [ "${MODE}" == "crearcanales" ]; then 
  generarCanales
elif [ "${MODE}" == "remplazar" ]; then 
  remplazarClavesPrivadas
elif [ "${MODE}" == "limpiar" ]; then 
limpiarRed
else
  echo "****${MODE}"
  ayuda
  exit 1
fi