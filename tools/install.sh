#!/bin/bash

#docker version: 17.0 
#docker-compose version: 1.7.1 
#Board version: 7.0.0


item=0

if [ -n "$1" && -f "$1" ];then
tar zxvf $1 -C /data
else
	if [ ! -e "/data/pre-env"]
	then
		echo "Please add the file pre-env.tar.gz file for add node to the directory Deploy!"
		exit 1
	fi
	echo "Welcome to Adminserver!"
fi

workdir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $workdir

function check_docker {
	if ! docker --version &> /dev/null
	then
		echo "Need to install docker(17.0+) first and run this script again."
		exit 1
	fi
	
	# docker has been installed and check its version
	if [[ $(docker --version) =~ (([0-9]+).([0-9]+).([0-9]+)) ]]
	then
		docker_version=${BASH_REMATCH[1]}
		docker_version_part1=${BASH_REMATCH[2]}
		docker_version_part2=${BASH_REMATCH[3]}
		
		# the version of docker does not meet the requirement
		if [ "$docker_version_part1" -lt 17 ] || ([ "$docker_version_part1" -eq 17 ] && [ "$docker_version_part2" -lt 0 ])
		then
			echo "Need to upgrade docker package to 17.0+."
			exit 1
		else
			echo "docker version: $docker_version"
		fi
	else
		echo "Failed to parse docker version."
		exit 1
	fi
}

function check_dockercompose {
	if ! docker-compose --version &> /dev/null
	then
		echo "Need to install docker-compose(1.7.1+) by yourself first and run this script again."
		exit $?
	fi
	
	# docker-compose has been installed, check its version
	if [[ $(docker-compose --version) =~ (([0-9]+).([0-9]+).([0-9]+)) ]]
	then
		docker_compose_version=${BASH_REMATCH[1]}
		docker_compose_version_part1=${BASH_REMATCH[2]}
		docker_compose_version_part2=${BASH_REMATCH[3]}
		
		# the version of docker-compose does not meet the requirement
		if [ "$docker_compose_version_part1" -lt 1 ] || ([ "$docker_compose_version_part1" -eq 1 ] && [ "$docker_compose_version_part2" -lt 6 ])
		then
			echo "Need to upgrade docker-compose package to 1.7.1+."
		else
			echo "docker-compose version: $docker_compose_version"
		fi
	else
		echo "Failed to parse docker-compose version."
		exit 1
	fi
}

echo "[Step $item]: checking installation environment ..."; let item+=1
check_docker
check_dockercompose

if [ -f board*.tgz ]
then
	echo "[Step $item]: loading Board & Adminserver images ..."; let item+=1
	docker load -i ./board*.tgz
fi
echo ""

echo "[Step $item]: checking existing instance of Adminserver ..."; let item+=1
if [ -n "$(docker ps -aqf name=deploy_log_1)"  ]
then
	echo "stopping existing Board instance ..."
	docker-compose -f ./docker-compose-rest.yml down
fi
echo ""
if [ -n "$(docker ps -aqf name=deploy_db_1)"  ]
then
	echo "stopping existing Database instance ..."
	docker-compose -f ./docker-compose-db.yml down
fi
echo ""
if [ -n "$(docker ps -aqf name=deploy_adminserver_1)"  ]
then
	echo "stopping existing Adminserver instance ..."
	docker-compose -f ./docker-compose-adminserver.yml down
fi
echo ""

echo "[Step $item]: starting Board-adminserver ..."
docker-compose -f ./docker-compose-adminserver.yml up -d

echo $"----Board-adminserver has been installed and started successfully.----

You can visit it on http://your-IP:8082 .

For more details, please visit http://open.inspur.com/TechnologyCenter/board .
"

