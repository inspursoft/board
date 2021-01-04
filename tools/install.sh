#!/bin/bash

#docker version: 17.0 
#docker-compose version: 1.7.1 
#Board version: 0.8.0

set -e

usage=$'Please set hostname and other necessary attributes in board.cfg first. DO NOT use localhost or 127.0.0.1 for hostname, because Board needs to be accessed by external clients.'
item=0

while [ $# -gt 0 ]; do
        case $1 in
            --help)
            echo "$usage"
            exit 0;;
            *)
            echo "$usage"
            exit 1;;
        esac
        shift || true
done

workdir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $workdir

# The hostname in board.cfg has not been modified
if grep 'hostname = reg.mydomain.com' &> /dev/null board.cfg
then
	echo $usage
	exit 1
fi

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
	echo "[Step $item]: loading Board images ..."; let item+=1
	docker load -i ./board*.tgz
fi
echo ""

option=legacy

if [[ $(cat ./board.cfg) =~ devops_opt[[:blank:]]*=[[:blank:]]*(gitlab?) ]]
then
option=${BASH_REMATCH[1]}
docker run --rm -v $(pwd)/board.cfg:/app/instance/board.cfg gitlab-helper:1.0
fi

echo "[Step $item]: preparing environment ...";  let item+=1
#if [ -n "$host" ]
#then
#	sed "s/^hostname = .*/hostname = $host/g" -i ./board.cfg
#fi
./prepare
echo ""

protocol=http
hostname=reg.mydomain.com

if [[ $(cat ./board.cfg) =~ ui_url_protocol[[:blank:]]*=[[:blank:]]*(https?) ]]
then
protocol=${BASH_REMATCH[1]}
fi

if [[ $(grep 'hostname[[:blank:]]*=' ./board.cfg) =~ hostname[[:blank:]]*=[[:blank:]]*(.*) ]]
then
hostname=${BASH_REMATCH[1]}
fi

if [ $option == "legacy" ]
then
	cd archive
fi

echo "[Step $item]: checking existing instance of Board ..."; let item+=1
if [ -n "$(docker-compose ps -q)"  ]
then
	echo "stopping existing Board instance ..."
	docker-compose down
fi
echo ""

echo "[Step $item]: creating Board network ..."; let item+=1
docker network create board &> /dev/null || true

echo "[Step $item]: starting Board ..."
docker-compose up -d

echo ""

echo $"----Board has been installed and started successfully.----

Now you should be able to visit the admin portal at ${protocol}://${hostname}. 
For more details, please visit https://github.com/inspursoft/board .
"
