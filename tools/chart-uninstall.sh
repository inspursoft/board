#!/bin/bash

#docker version: 17.0+
#docker-compose version: 1.7.1+
#Board version: 0.8.0+

set -e

usage=$'This shell script will uninstall Board. Only run it under the installation directory. \nUsage:    uninstall [OPTINOS]  \nOptions:\n  -s      Silent uninstall.\n  --help  Show this help info.'
item=0

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
if  [ ! -d charts ] 
then
	echo "Charts directory does not exist."
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

function check_helm {
    if ! helm version &> /dev/null;then
        echo $"Helm is required but not found!
Please make sure that Helm client or Helm server is ready!
For more details, please visit: https://v2.helm.sh/docs/"
        exit $?
    else
        if [[ $(helm version) =~ ((Client:.*).(Server:.*)) ]]; then
            client_version=${BASH_REMATCH[2]}
            server_version=${BASH_REMATCH[3]}

            if [[ $client_version =~ (([0-9]+).([0-9]+).([0-9]+)) ]]; then
                client_version_main=${BASH_REMATCH[2]}
                echo "client version:" $client_version_main

                if [ "$client_version_main" -ne 2 ]
                then
                    echo "Only support helm 2!"
                    exit 1
                fi
            else
                echo "Unknown version"
                exit 1
            fi
            
            if [[ $server_version =~ (([0-9]+).([0-9]+).([0-9]+)) ]]; then
                server_version_main=${BASH_REMATCH[2]}
                echo "server version:" $server_version_main

                if [ "$server_version_main" -ne 2 ]
                then
                    echo "Only support helm 2!"
                    exit 1
                fi
            else
                echo "Unable to get helm server version"
                exit 1
            fi
        else
            echo "Failed to get helm version."
            exit 1
        fi
        echo "helm is ok"

    fi
}

echo "[Step $item]: checking uninstallation environment ..."; let item+=1
check_docker
check_helm

sed -i "s/^hostname.*$/hostname = reg.mydomain.com/" board.cfg

if [[ $(cat ./config/apiserver/env) =~ DEVOPS_OPT=(legacy?) ]]
then
cd archive
fi

echo "[Step $item]: checking existing instance of Board ..."; let item+=1
if [ -n "$(helm ls board --all -q)"  ]
then
	echo "Delete existing Board instance ..."
	helm del --purge board
fi
echo ""

echo $"----Board uninstaller running complete.----

For more details, please visit https://github.com/inspursoft .
"
