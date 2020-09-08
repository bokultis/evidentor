#!/bin/bash

#colors
RED='\033[0;31m'
GREEN='\033[32;1m'
NC='\033[0m' # No Color

GOBUILD="go build"
GOCLEAN="go clean"
GOTEST="go test"
GOGET="go get"
BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
MYSQL_FILES_PATH="${BASEDIR}/dev_env/mysqlfiles"
BINARY_PATH="${BASEDIR}/cmd/bin"
BINARY_NAME="evidentor"
BINARY_UNIX="${BINARY_NAME}_unix"

DOCKER_API="evidentor-api"
DOCKER_DATABASE="mysqldb-evidentor"

DB_USER="root"
DB_PASS="root"
DB_NAME="evidentor"

ACTIONS=("start" "start_swagger" "docker_db" "docker_api" "docker_kill_all")

function start_swagger {

    stop_swagger 

    #start swagger
    printf "\nPlease wait while we boot up Swagger container \n"
    docker run -d -p 1200:8080 -e URL=http://localhost:1201/evidentor.yaml swaggerapi/swagger-ui

    cd $BASEDIR/docs
    php -S localhost:1201 router.php &
    SWAGGER_PID=$!

    printf "\nSwagger PID: ${SWAGGER_PID} UI available @ ${GREEN}http://localhost:1200${NC}\n"
}

function stop_swagger {
    #kill php proccess if any that holds port used for display yaml
    SWAGGER_PID=$(ps aux | grep '[l]ocalhost:1201' | awk '{print $2}')
    if [ ! -z $SWAGGER_PID ]
    then
        printf "\n\nKilling swagger: $SWAGGER_PID ...\n"
        kill $SWAGGER_PID
    fi
    # stop swagger container
    docker stop $(docker ps -a -q --filter ancestor=swaggerapi/swagger-ui --format="{{.ID}}") > /dev/null 2>&1
}


function action_docker_db {
    docker rm -f ${DOCKER_DATABASE}
    docker run --name ${DOCKER_DATABASE} -d -v ${MYSQL_FILES_PATH}:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=${DB_PASS} -d mysql:5.7
}

function action_docker_api {
    go clean
    rm -f ${BINARY_PATH}/${BINARY_NAME}
    
    cd $BASEDIR/api
    ${GOBUILD}  -o ${BINARY_PATH}/${BINARY_NAME} -v

    cd ${BASEDIR}/cmd/ 
    ./image-build-docker.sh latest  > /dev/null 2>&1

    cd ${BASEDIR}
    docker rm -f ${DOCKER_API}
    docker run -d --name ${DOCKER_API} --network="host" -v ${BASEDIR}/log:/var/log  --env-file=.env evidentor:latest 
}

function action_docker_kill_all {
    docker container stop $(docker container ls -aq)
    docker container rm $(docker container ls -aq)
}

function backup_database {
    # Backup
    docker exec ${DOCKER_DATABASE} /usr/bin/mysqldump -u ${DB_USER} --password=${DB_PASS} ${DB_NAME} > backup_${DB_NAME}.sql > /dev/null 2>&1
}

function restore_database {
    # Restore
    cat backup_${DB_NAME}.sql | docker exec -i ${DOCKER_DATABASE} /usr/bin/mysql -u ${DB_USER} --password=${DB_PASS} ${DB_NAME} > /dev/null 2>&1
}

function schedule_db_backup {
    trap cleanup SIGINT
    cleanup(){
        backup_database
        stop_swagger
        printf "\n================================================\n"
        printf "\nYou are done and exited ${GREEN}cleanly${NC}!\n"
        printf "\n================================================\n"
        exit 0
    }
    while true  
    do  
        backup_database
        sleep 5 
    done
}

function action_start {
    printf "\n======== ${RED}Evidentor Development Environment${NC} ========\n" 
    echo "------------------"
    printf "\nStarting ${GREEN}Mysql Docker Container${NC} ... \n"    
    action_docker_db

    printf "\nStarting ${GREEN}Evidentor API Docker Container${NC} ... \n" 
    action_docker_api

    start_swagger > /dev/null
    echo "------------------"
    printf "\nSwagger UI: ${GREEN}http://localhost:1200${NC}\n"
    echo "------------------"
    
    gnome-terminal --command="bash -c 'tail -f -n 25 ${BASEDIR}/log/logrus.log; $SHELL'"

    printf "\nRestore Database ${GREEN}Docker ${DOCKER_DATABASE} : ${DB_NAME}${NC} ... \n" 
    restore_database

    printf "\nTo clean exit press Ctrl+c \n"
    schedule_db_backup
    
}

function run_action {
    local ACTION_NAME="action_$1"
    if typeset -f "${ACTION_NAME}" > /dev/null; then
        $ACTION_NAME $2 $3
    else
        printf "\nAction [$1] is not available.\nAvailable actions: "
        echo "${ACTIONS[@]}"
        printf "\n\n"
        exit 1
    fi    
}

run_action $1 $2 $3

