#!/bin/bash

# declare an array of hosts to add to 127.0.0.1
hosts=("mongo" "mongo1" "mongo2" "postgres-coinbasepro" "postgres-polyon" "redis1")

for i in "${hosts[@]}"
do
	if [ -n "$(grep $HOSTNAME /etc/hosts)" ]
	then
		echo "$i already exists: $(grep $i /etc/hosts)"
	else
		sudo -- sh -c -e "echo '127.0.0.1 $i' >> /etc/hosts"
	fi
done

# Print the contents of the /etc/hosts file
cat /etc/hosts
