
include .$(DB_PWD)/app/.env

start: build-test-db run-test-db init-test-db
build: build-test-db
run: run-test-db
initdb: init-test-db
rm: remove-last-container

LAST_CT_IT:=$(shell docker ps -l -q)

.PHONY: build-test-db
build-test-db: 
	@docker build --rm -f ./db/.Dockerfile --build-arg "DB_PWD=$(DB_PWD)" -t auth-db ./db/

.PHONY: run-test-db
run-test-db: 
	@docker run -d -p 5432:5432 --name auth-db auth-db 

.PHONY: init-test-db
init-test-db:
	@sleep 25
	@cat ./db/init_db.sql | docker exec -i auth-db psql -U bilrafal -d auth

.PHONY: remove-last-container
remove-last-container:
	@docker stop auth-db	
	@sleep 5	
	@docker rm auth-db	
	@sleep 5
	@docker volume prune

.PHONY: sleep20
sleep20: 
	@sleep 20