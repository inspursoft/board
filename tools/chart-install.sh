#!/bin/bash

#docker version: 17.0
#docker-compose version: 1.7.1 
#Board version: 0.8.0

set -e

namespace=$1
namespace=${namespace:-board}
version_tag=$(cat VERSION | head -n 1)

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

function confirm {
        read -r -p "$1" input
        case $input in
                y|Y)
                        echo "$2"
                        ;;
                n|N)
                        echo "$3"
                        exit 0
                        ;;
                *)
                        echo "Invalid input..."
                        exit 1
                        ;;
        esac
}

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

function load_images {
	# Check if the tar package exists
	if [ -f board*.tgz ]
	then
		echo "[Step $item]: loading Board images ..."; let item+=1
		docker load -i ./board*.tgz
	else
		echo "Can not find board*.tgz"
		exit 1
	fi
	# Parse image_registry_url
	if [[ $(cat ./board.cfg) =~ registry_ip[[:blank:]]*=[[:blank:]]*([0-9a-zA-Z._/:-]*) ]]
	then
		registry_ip=${BASH_REMATCH[1]}
		echo "Parse registry_ip = $registry_ip"
	else
		echo "Failed to parse registry_ip in board.cfg"
		exit 1
	fi
	if [[ $(cat ./board.cfg) =~ registry_port[[:blank:]]*=[[:blank:]]*([0-9a-zA-Z._/:-]*) ]]
	then
		registry_port=${BASH_REMATCH[1]}
		echo "Parse registry_port = $registry_port"
	else
		echo "Failed to parse registry_port in board.cfg"
		exit 1
	fi
	image_registry_url=$registry_ip:$registry_port
	echo "image_registry_url = $image_registry_url"
	# docker tag and push images to registry
	for image in $(docker images --format "{{.Repository}}:{{.Tag}}" | grep $version_tag | grep -v $image_registry_url);
	do
		docker tag $image $image_registry_url/openboard/$image &> /dev/null
		docker push $image_registry_url/openboard/$image &> /dev/null
		echo "Push $image_registry_url/openboard/$image ok"
	done
	docker tag docker.elastic.co/elasticsearch/elasticsearch:7.9.3 $image_registry_url/elasticsearch/elasticsearch:7.9.3
	docker tag docker.elastic.co/kibana/kibana:7.9.3 $image_registry_url/kibana/kibana:7.9.3
	docker tag quay.io/fluentd_elasticsearch/fluentd:v3.0.4 $image_registry_url/fluentd_elasticsearch/fluentd:v3.0.4
	docker push $image_registry_url/elasticsearch/elasticsearch:7.9.3 &> /dev/null
	docker push $image_registry_url/kibana/kibana:7.9.3 &> /dev/null
	docker push $image_registry_url/fluentd_elasticsearch/fluentd:v3.0.4 &> /dev/null
}

function create_pv {
        if [[ $(cat ./board.cfg) =~ nfs_path[[:blank:]]*=[[:blank:]]*([0-9a-zA-Z._/:-]*) ]]
        then
                nfs_path=${BASH_REMATCH[1]}
                echo "Parse nfs_path = $nfs_path"
        else
                echo "Failed to parse nfs_path in board.cfg"
                exit 1
        fi
        if [[ $(cat ./board.cfg) =~ nfs_server[[:blank:]]*=[[:blank:]]*([0-9a-zA-Z._/:-]*) ]]
        then
                nfs_server=${BASH_REMATCH[1]}
                echo "Parse nfs_server = $nfs_server"
        else
                echo "Failed to parse nfs_server in board.cfg"
                exit 1
        fi

	sed "s#__nfs_path__#$nfs_path#g" pv.tpl > pv.yaml
	sed "s#__nfs_server__#$nfs_server#g" -i pv.yaml
	kubectl create -f pv.yaml
}

echo "[Step $item]: checking installation environment ..."; let item+=1
check_docker
check_helm

load_images

# Install gitlab
option=gitlab

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
./prepare_chart --conf board.cfg --tag $version_tag
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
if [ -n "$(helm ls board --all -q)"  ]
then
	echo "Delete existing Board instance ..."
	helm del --purge board
fi
echo ""

echo "[Step $item]: checking nfs server active ..."; let item+=1
confirm "Please confirm NFS is configured. Are you sure to continue? [Y/n]" "Continue." "Please start and configure NFS server."
if [[ $(systemctl status nfs | grep Active) =~ Active[[:blank:]]*:[[:blank:]]*(active?) ]]
then
        read -r -p "Do you need to create pv? [Y/n]" input
        case $input in
                y|Y)
                        echo "Creating..."
                        create_pv
                        ;;
                n|N)
                        ;;
                *)
                        echo "Invalid input..."
                        exit 1
                        ;;
        esac
else
	echo "NFS serve is not active."
	exit 1
fi

# Confirm whether deleted claimRef in pv
confirm "Please confirm that the status of each pv is Available. Are you sure to continue? [Y/n]" "Continue."

echo "[Step $item]: starting Board ..."
echo "namespace = $namespace"
helm install --name board --namespace $namespace charts/board

echo ""

echo $"----Board has been installed and started successfully.----

Now you should be able to visit the admin portal at ${protocol}://${hostname}. 
For more details, please visit https://github.com/inspursoft/board .
"
