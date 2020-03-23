#!/bin/bash

#DOCKER variables
DOCKER_REGISTRY="xxxdev-dockerregistry.horisen.org:5043"
DOCKER_USERNAME="xxxhorisendev-team"
DOCKER_PASSWORD="xxx39TSK0Jv38085119"

RED='\033[0;31m'
GREEN='\033[32;1m'
NC='\033[0m' # No Color

BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

ACTIONS=("start" "install" "update" "get" "build" "deploy-vagrant" "deploy" "log" "bump-version"  "desc" "docker")

#DEPLOY config
DEPL_SOURCE_DIR="cmd/link-li-api/"

declare -A DEPL_TARGET_DIR=(
    ["staging"]="xxxadmin@194.0.137.77:/opt/link-li-api/"
    ["production"]="xxxadmin@194.0.137.78:/opt/link-li-api/"
)
declare -A DEPL_TARGET_CACHE=(
)

DEPL_DEF_PARAM="--rsh=\"ssh -p8142\" --exclude-from=./deploy_exclude.txt --progress"
declare -A DEPL_TARGET_PARAM=(
    ["staging"]=$DEPL_DEF_PARAM
    ["production"]=$DEPL_DEF_PARAM
)

declare -A DEPL_SSH=(
    ["staging"]="-p 8142 xxxadmin@194.0.137.77"
    ["production"]="-p 8142 xxxadmin@194.0.137.78"
)
declare -A DEPL_CMD=(
    ["staging"]="pkill -f link-li-api"
    ["production"]="pkill -f link-li-api"
)

declare -A DEPL_NAME_CMD=(
    ["staging"]="cd /opt/link-li-api/;mv link-li-api link-li-api-old"
    ["production"]=""
)

declare -A DEPL_LOG=(
    ["staging"]="/opt/link-li-api/logs/app.log"
    ["production"]="/opt/link-li-api/logs/app.log"
)


function start_swagger {
    printf "\nPlease wait while we boot up Swagger container \n"
    docker run -d -p 1200:8080 -e URL=http://localhost:1201/link-li-api.yaml swaggerapi/swagger-ui

    cd $BASEDIR/docs/api
    php -S localhost:1201 router.php &
    SWAGGER_PID=$!

    printf "\nSwagger UI available @ ${GREEN}http://localhost:1200${NC}\n"  
}

function stop_swagger {
    printf "\n\nKilling swagger: $SWAGGER_PID...\n"
    kill $SWAGGER_PID
    docker stop $(docker ps -a -q --filter ancestor=swaggerapi/swagger-ui --format="{{.ID}}")        
}

function action_start {
    start_swagger

    printf "\nStarting Link.li redirector and renderer...\n"    
    cd $BASEDIR/cmd/link-li-redirector
    ./link-li-redirector&
    REDIRECTOR_PID=$!

    printf "\nStarting Link.li page designer @ ${GREEN}http://localhost:8002/{pageId}${NC} ...\n"    
    cd $BASEDIR/docs/grape.js
    php -S localhost:8002 -t ./ &
    GRAPE_PID=$!
 

    trap cleanup SIGINT
    cleanup()
    {
        printf "\n\nKilling pids: $REDIRECTOR_PID, $GRAPE_PID...\n"
        kill $REDIRECTOR_PID $GRAPE_PID

        stop_couchbase
        stop_swagger

        exit 0
    }

    action_desc

    cd $BASEDIR/cmd/link-li-api
    ./link-li-api    
}

function action_desc {
    printf "\n========${RED}Link.li Development Environment${NC}========\n"    
    printf "\nPage Editor: ${GREEN}http://localhost:8002/e509c2ec-c62c-462e-89b0-26ed2d309d85${NC}\n"
    echo "------------------"
    printf "\nSwagger UI: ${GREEN}http://localhost:1200${NC}\n"
    echo "------------------"
    printf "\nCB UI: ${GREEN}http://localhost:8091/ui/index.html${NC}\n"
    printf "\n================================================\n"
}

function action_install {
    # golang installed or not?
    if [ -z "$GOPATH" ]; then
        # no golang, just get binaries 
        action_get
    else
        cd $GOPATH/src/git.horisen.biz
        if [ ! -d "corego" ]; then
            git clone ssh://git-user@git.horisen.biz:8142/opt/git/horisen-pro/corego.git
        fi
        if [ ! -d "bg" ]; then
            git clone ssh://git-user@git.horisen.biz:8142/opt/git/bg.git
        fi            
        go get github.com/sirupsen/logrus
        go get github.com/go-redis/redis
        go get github.com/jinzhu/gorm
        go get github.com/streadway/amqp
        go get github.com/nyaruka/phonenumbers
        go get github.com/asaskevich/govalidator
	    go get github.com/gofrs/uuid
        go get github.com/go-sql-driver/mysql
        go get github.com/gorilla/mux
        

        echo "Building binaries .."
        cd $BASEDIR
        action_build
    fi

    cd $BASEDIR
    echo "Preparing couchbase db ..."
    sudo mkdir -p /opt/couchbase/var
    sudo unzip docs/db/lib.zip -d /opt/couchbase/var

    printf "\n127.0.0.1 link.test\n" >> /etc/hosts

    echo "Starting dev. env. ..."
    action_start    
}

function action_get {
    URL=http://horisen:q1k4v1pN72624939@vgrnis.horisen.com:8134/link-li-api

    curl  $URL/link-li-api > cmd/link-li-api/link-li-api
    chmod +x cmd/link-li-api/link-li-api 
}

function action_update {
    git pull

    # golang installed or not?
    if [ -z "$GOPATH" ]; then
        # no golang, just get binaries 
        action_get
    else
        cd $GOPATH/src/git.horisen.biz/bg
        git pull        

        echo "Building binaries ..."
        cd $BASEDIR
        action_build
    fi

    printf "\nYou can now start dev env.\n\n"    
}

function action_build {
    BUILD_TIME=$(date --rfc-3339=seconds)
    BUILD_VERSION="0.0.0"
    BUILD_COMMIT=$(git rev-list -1 HEAD)
    if [ ! -f VERSION ]; then
        echo -e "${WARNING_FLAG} Could not find a VERSION file, please create it.${RESET}"        
    else      
        BUILD_VERSION=`cat VERSION`
    fi    
    LD_FLAGS="-X 'git.horisen.biz/corego/pkg/core.BuildTime=$BUILD_TIME' -X 'git.horisen.biz/corego/pkg/core.BuildVersion=$BUILD_VERSION' -X 'git.horisen.biz/corego/pkg/core.BuildCommit=$BUILD_COMMIT'"    

    cd $BASEDIR/cmd/link-li-api
    go build -ldflags="$LD_FLAGS"
}

function action_docker {
    action_build

    #get version
    cd $BASEDIR
    if [ ! -f VERSION ]; then
        echo -e "${WARNING_FLAG} Could not find a VERSION file, please create it.${RESET}"
        exit 1 
    else      
        TAG=`cat VERSION`
        TAG="v$TAG"
    fi

    #build and deploy
    cd $BASEDIR/cmd/link-li-api

    IMAGE_NAME="horisen/link-li-api"

    docker login -u ${DOCKER_USERNAME} -p ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}
    docker build -t ${DOCKER_REGISTRY}/${IMAGE_NAME}:${TAG} .
    docker tag ${DOCKER_REGISTRY}/${IMAGE_NAME}:${TAG} ${DOCKER_REGISTRY}/${IMAGE_NAME}:latest
    docker push ${DOCKER_REGISTRY}/${IMAGE_NAME}
    docker tag ${DOCKER_REGISTRY}/${IMAGE_NAME}:latest ${DOCKER_REGISTRY}/${IMAGE_NAME}:${TAG}


    docker push ${DOCKER_REGISTRY}/${IMAGE_NAME}
    docker logout
}


function action_deploy-vagrant {    
    scp $BASEDIR/cmd/link-li-api/link-li-api vagrant@192.168.1.134:/var/www/link-li-api/
    scp $BASEDIR/cmd/link-li-redirector/link-li-redirector vagrant@192.168.1.134:/var/www/link-li-api/
}

function action_deploy {
    TARGET=$1
    ACTION=$2

    TARGET_DIR=${DEPL_TARGET_DIR[$TARGET]}
    TARGET_PARAM=${DEPL_TARGET_PARAM[$TARGET]}
    CLEAN_CACHE_URL=${DEPL_TARGET_CACHE[$TARGET]}
    TARGET_SSH=${DEPL_SSH[$TARGET]}
    TARGET_CMD=${DEPL_CMD[$TARGET]}    
    TARGET_NAME_CMD=${DEPL_NAME_CMD[$TARGET]}

    if [ "$TARGET_DIR" = "" ]
    then
        MODE=""
        printf "\n\033[31;1m Target [$TARGET] is not defined \033[0m\n"
        printf "cmd format is: \033[32;1mdeploy {target} {mode}\033[0m\n\n"
        exit 1
    fi
    

    if [ "$ACTION" = "upload" ]
    then
        MODE=""
        if [ "$TARGET_NAME_CMD" != "" ]
        then
            CMD="ssh $TARGET_SSH \"${TARGET_NAME_CMD}\""
            printf "\n\nMaking backup file:\n$CMD\n"
            eval $CMD
        fi
        printf "\n\033[31;1m Running LIVE mode\033[0m\n"
    else
        MODE="-n"
        printf "\n\033[32;1m Running SIMULATION mode \033[0m\n"
    fi

    CMD="rsync $TARGET_PARAM -r $MODE -v -c   ${DEPL_SOURCE_DIR} $TARGET_DIR"

    printf "\n\ncmd to be executed:\n$CMD\n"

    eval $CMD    

    #cleanup cache
    if [ "$ACTION" = "upload" ] && [ $CLEAN_CACHE_URL ]
    then
        printf "\n\ngoing to clean cache at URL:\n$CLEAN_CACHE_URL\n"
        printf "\n\033[32;1mPress enter to continue\033[0m\n"
        read -p ""
        curl $CLEAN_CACHE_URL
    fi

    #restart server
    if [ "$ACTION" = "upload" ] && [ "$TARGET_SSH" ] && [ "$TARGET_CMD" ]
    then
        CMD="ssh $TARGET_SSH \"${TARGET_CMD}\""
        printf "\n\nremote cmd to be executed:\n$CMD\n"     
        eval $CMD
    fi        
}

function action_log  {
    TARGET=$1

    TARGET_SSH=${DEPL_SSH[$TARGET]}
    TARGET_LOG=${DEPL_LOG[$TARGET]}

    if [ "$TARGET_SSH" = "" ] || [ "$TARGET_LOG" = "" ]
    then
        printf "\n\033[31;1m Target [$TARGET] is not defined \033[0m\n"
        printf "cmd format is: \033[32;1mlog {target}\033[0m\n\n"
        exit 1
    fi
    CMD="ssh -t $TARGET_SSH \"tail -f $TARGET_LOG\""
    eval $CMD
}

function action_bump-version  {
    NOW="$(date +'%B %d, %Y')"
    RED="\033[1;31m"
    GREEN="\033[0;32m"
    YELLOW="\033[1;33m"
    BLUE="\033[1;34m"
    PURPLE="\033[1;35m"
    CYAN="\033[1;36m"
    WHITE="\033[1;37m"
    RESET="\033[0m"

    LATEST_HASH=`git log --pretty=format:'%h' -n 1`

    QUESTION_FLAG="${GREEN}?"
    WARNING_FLAG="${YELLOW}!"
    NOTICE_FLAG="${CYAN}❯"

    ADJUSTMENTS_MSG="${QUESTION_FLAG} ${CYAN}Now you can make adjustments to ${WHITE}CHANGELOG.md${CYAN}. Then press enter to continue."
    PUSHING_MSG="${NOTICE_FLAG} Pushing new version to the ${WHITE}origin${CYAN}..."

    if [ ! -f VERSION ]; then
        echo -e "${WARNING_FLAG} Could not find a VERSION file, please create it.${RESET}"
        exit 1
    fi    

    BASE_STRING=`cat VERSION`
    BASE_LIST=(`echo $BASE_STRING | tr '.' ' '`)
    V_MAJOR=${BASE_LIST[0]}
    V_MINOR=${BASE_LIST[1]}
    V_PATCH=${BASE_LIST[2]}
    echo -e "${NOTICE_FLAG} Current version: ${WHITE}$BASE_STRING"
    echo -e "${NOTICE_FLAG} Latest commit hash: ${WHITE}$LATEST_HASH"
    V_MINOR=$((V_MINOR + 1))
    V_PATCH=0
    SUGGESTED_VERSION="$V_MAJOR.$V_MINOR.$V_PATCH"
    echo -ne "${QUESTION_FLAG} ${CYAN}Enter a version number [${WHITE}$SUGGESTED_VERSION${CYAN}]: "
    read INPUT_STRING
    if [ "$INPUT_STRING" = "" ]; then
        INPUT_STRING=$SUGGESTED_VERSION
    fi
    echo -e "${NOTICE_FLAG} Will set new version to be ${WHITE}$INPUT_STRING"
    echo $INPUT_STRING > VERSION
    echo "## $INPUT_STRING ($NOW)" > tmpfile
    git log --pretty=format:"  - %s" "api-$BASE_STRING"...HEAD >> tmpfile
    echo "" >> tmpfile
    echo "" >> tmpfile
    cat CHANGELOG.md >> tmpfile
    mv tmpfile CHANGELOG.md
    echo -e "$ADJUSTMENTS_MSG"
    read
    echo -e "$PUSHING_MSG"
    git add CHANGELOG.md VERSION
    git commit -m "Bump version to ${INPUT_STRING}."
    git tag -a -m "Tag version ${INPUT_STRING}." "api-$INPUT_STRING"
    git push origin --tags
    git push

    echo -e "${NOTICE_FLAG} Finished.${RESET}"
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