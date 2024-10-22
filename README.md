# Simulation_Modbus

## Build and run
Build
```
docker build -t texolaurora/simulation_modbus .
```
Run
```
docker run --rm -it --name simulation_modbus -p 502:502 texolaurora/simulation_modbus Sensor01.yml
```
Run SH
```
docker run --rm -it --entrypoint  sh texolaurora/simulation_modbus
```
Tag
```
docker tag simulation_modbus:0.2.eval texolaurora/simulation_modbus:eval
```
list all running Docker containers
```
docker ps -a
```
list all Docker images
```
docker images
```
Delete Docker images
```
docker rmi texolaurora/simulation_modbus:0.1.eval
```
Stop and Delete Container
```
docker stop container_id
docker rm container_id
```
