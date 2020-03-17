#!/bin/bash

#docker version: 17.0+
#docker-compose version: 1.7.1+
#Board version: 0.8.0+

set -e

usage=$'This shell script will uninstall Board images and data volume. Only run it under the installation directory. \nUsage:    uninstalil [OPTINOS]  \nOptions:\n  -s      Silent uninstall.\n  --help  Show this help info.'
item=0
defaultDataVolume="/data/board"
silentFlag=flase

while [ $# -gt 0 ]; do
        case $1 in
            --help)
            echo "$usage"
            exit 0;;
            -s)
            echo "Uninstall without any user interaction."
            silentFlag=true
            ;;
            *)
            echo "$usage"
            exit 1;;
        esac
        shift || true
done

# The hostname in board.cfg has not been modified
if  [ ! -f docker-compose*.yml ] 
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

function delete_images {
	docker-compose down --rmi all
}

function remove_data {
	rm -rf $defaultDataVolume
}

echo "[Step $item]: checking uninstallation environment ..."; let item+=1
check_docker
check_dockercompose

echo "[Step $item]: checking existing instance of Board ..."; let item+=1
if [ -n "$(docker-compose ps -q)"  ]
then
	echo "stopping existing Board instance ..."
	docker-compose down
fi
echo ""

echo "[Step $item]: remove Board images..."; let item+=1
	delete_images
echo ""

echo "[Step $item]: prepare removing Board data..."

if [ $silentFlag == "true" ]
then 
        echo "start deleting..."
        remove_data
        echo "Done."
else
        if read -t 10 -p "Really want to delete Board data? Please input [yes] to confirm: " flag
        then
                if [ $flag == "yes" ]
                then
                        echo "You input [$flag] for deletion, start data deletion after 5 seconds..."
                        sleep 5s
                        echo "start deleting..."
                        remove_data
                        echo "Done."
                else
                        echo "You input [$flag], skip data deletion."
                fi
        else
                echo ""
                echo "Sorry ,timeout!"
        fi
fi
	
echo ""

echo $"----Board uninstaller running complete.----
For more information, please visit http://10.110.18.40:10080/inspursoft/board"

