#!/bin/bash
echo
echo Arrancando un engine que atiende peticiones por localhost:9092, limita el parque a 5 visitantes y manda peticiones a un servidor de tiempos de espera situado en localhost:9094.
echo
./fwq_engine localhost 9092 5 localhost 9094

