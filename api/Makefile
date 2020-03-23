# Go parameters
    GOCMD=go
    GOBUILD=$(GOCMD) build
    GOCLEAN=$(GOCMD) clean
    GOTEST=$(GOCMD) test
    GOGET=$(GOCMD) get
    ROOT_PATH=$(abspath $(dir $(lastword $(MAKEFILE_LIST)))/..)
    MYSQL_FILES_PATH=$(ROOT_PATH)/dev_env/mysqlfiles
    BINARY_PATH=$(ROOT_PATH)/cmd/bin
    BINARY_NAME=evidentor
    BINARY_UNIX=$(BINARY_NAME)_unix
    
	
    
    all: test build
    build: 
			$(GOBUILD)  -o $(BINARY_PATH)/$(BINARY_NAME) -v
    test: 
			$(GOTEST) -v 
    clean: 
			$(GOCLEAN)
			rm -f $(BINARY_PATH)/$(BINARY_NAME)
    run:
			$(GOBUILD) -o $(BINARY_PATH)$(BINARY_NAME) -v 
			./$(BINARY_PATH)$(BINARY_NAME)
    
	# Cross compilation
    build-linux:
			CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

    build-windows:
				CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

    docker-db:
			
			-docker rm -f mysqldb-evidentor

			docker run --name mysqldb-evidentor -v $(MYSQL_FILES_PATH):/var/lib/mysql -e MYSQL_ROOT_PASSWORD=root -d mysql:5.7

    docker-api:
			
			$(GOCLEAN)
			rm -f $(BINARY_PATH)/$(BINARY_NAME)

			$(GOBUILD)  -o $(BINARY_PATH)/$(BINARY_NAME) -v

			$(ROOT_PATH)/cmd/./image-build-docker.sh latest

			-docker rm -f evidentor-api

			docker run -it --name evidentor-api --network="host"  --env-file $(ROOT_PATH)/.env evidentor:latest
    
    docker-build:
			
			$(GOCLEAN)
			rm -f $(BINARY_PATH)/$(BINARY_NAME)

			$(GOBUILD)  -o $(BINARY_PATH)/$(BINARY_NAME) -v

			$(ROOT_PATH)/cmd/./image-build-docker.sh latest

			-docker rm -f mysqldb-evidentor
			-docker rm -f evidentor-api

			docker run --name mysqldb-evidentor -v $(MYSQL_FILES_PATH):/var/lib/mysql -e MYSQL_ROOT_PASSWORD=root -d mysql:5.7

			docker run -it --name evidentor-api --network="host"  --env-file $(ROOT_PATH)/.env evidentor:latest
			