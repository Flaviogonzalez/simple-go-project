broker_binary=brokerApp
listener_binary=listenerApp
auth_binary=authApp

build_listener:
	chdir ..\listener-service && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${listener_binary} ./

build_broker:
	chdir ..\broker-service && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${broker_binary} ./

build_auth:
	chdir ..\auth-service && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${auth_binary} ./

build: build_broker build_listener build_auth
	docker-compose up --build -d

buildonly: build_broker build_listener build_auth
	docker-compose up --build -d broker-service listener-service auth-service
