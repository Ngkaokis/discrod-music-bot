.PHONY:dev 
dev:
	air -c air.toml

.PHONY:build 
build:
	go build -o app	

.PHONY:run 
run:
	./app	
