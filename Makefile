include /.$(PWD)/app/.env

run: run-test-db
rm: remove-last-container

LAST_CT_IT:=$(shell docker ps -l -q)

.PHONY: run-test-db
run-test-db: 
	@echo $(DB_PWD)
	@docker run -e 'ACCEPT_EULA=Y' -e 'MSSQL_SA_PASSWORD=$(DB_PWD)' -p 1499:1433 -d mcr.microsoft.com/mssql/server:2017-latest
	@sleep 21	
	@sqlcmd -S localhost,1499 -U sa -P $(DB_PWD) -d master -i db/create-db.sql			

.PHONY: remove-last-container
remove-last-container:
	@docker stop $(LAST_CT_IT)	
	@sleep 10	
	@docker rm $(LAST_CT_IT)	