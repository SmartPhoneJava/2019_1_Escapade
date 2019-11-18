#!/bin/sh
echo ""
echo "  ----------------------------"
echo "  ----Docker swarm prepare----"
echo "  ----------------------------"
echo ""
trap 'echo " stop" ' INT TERM
#chmod +x prepare.sh && ./prepare.sh

echo "  0. Prepare other .sh"
chmod +x droplet_add.sh
chmod +x droplet_ufw.sh
chmod +x droplet_ufw.sh
chmod +x droplet_manager.sh
chmod +x droplet_worker.sh

# Укажите в переменной path путь до ssh ключа.
# Либо оставьте "id_rsa", но скопируйте id_rsa и id_rsa.pub
# в данный каталог 2019_1_Escapade/cmd/commands/
path="id_rsa"

# адреса серверов
addr1="142.93.79.194"
addr2="167.71.247.116"
addr3="138.197.68.111"
addr4="167.71.247.152"
addr5="167.71.243.102"

# названия машин
name1="api1"
name2="api2"
name3="api3"
name4="api4"
name5="api5"

# токен swarm менеджера
swarm_token="SWMTKN-1-34rkcp8qfted6l0lszq8fro121h1h3jkn90qmim45js4rolqae-8qfuwtpmja336e9pb9jf7t87o"

echo "  1. Create machines"

echo "  1.1. Create machine - "$addr1
./droplet_add.sh $addr1 $path $name1
./droplet_ufw.sh $addr1
./droplet_metrics.sh $addr1
./droplet_manager.sh $addr1

echo "  1.2. Create machine - "$addr2
./droplet_add.sh $addr2 $path $name2
./droplet_ufw.sh $addr2
./droplet_metrics.sh $addr2
./droplet_worker.sh $addr2 $swarm_token $addr1

echo "  1.3. Create machine - "$addr3
./droplet_add.sh $addr3 $path $name3
./droplet_ufw.sh $addr3
./droplet_metrics.sh $addr3
./droplet_worker.sh $addr3 $swarm_token $addr1

echo "  1.4. Create machine - 167.71.247.152"
./droplet_add.sh $addr4 $path $name4
./droplet_ufw.sh $addr4
./droplet_metrics.sh $addr4
./droplet_worker.sh $addr4 $swarm_token $addr1

echo "  1.5. Create machine - 167.71.243.102"
./droplet_add.sh $addr5 $path $name5
./droplet_ufw.sh $addr5
./droplet_metrics.sh $addr5
./droplet_worker.sh $addr5 $swarm_token $addr1

# echo "  1. Create machines"
# echo "  1.1. Create machine - 142.93.79.194"
# docker-machine create --driver=generic \
#     --generic-ip-address=142.93.79.194 \
#     --generic-ssh-user=root \
#     --generic-ssh-key=id_rsa \
#         api1

# echo "  1.2. Create machine - 167.71.247.116"
# docker-machine create --driver=generic \
#     --generic-ip-address=167.71.247.116 \
#     --generic-ssh-user=root \
#     --generic-ssh-key=id_rsa \
#         api2

# echo "  1.3. Create machine - 138.197.68.111"
# docker-machine create --driver=generic \
#     --generic-ip-address=138.197.68.111 \
#     --generic-ssh-user=root \
#     --generic-ssh-key=id_rsa \
#         api3

# echo "  1.4. Create machine - 167.71.247.152"
# docker-machine create --driver=generic \
#     --generic-ip-address=167.71.247.152 \
#     --generic-ssh-user=root \
#     --generic-ssh-key=id_rsa \
#         api4

# echo "  1.5. Create machine - 167.71.243.102"
# docker-machine create --driver=generic \
#     --generic-ip-address=167.71.243.102 \
#     --generic-ssh-user=root \
#     --generic-ssh-key=id_rsa \
#         api5