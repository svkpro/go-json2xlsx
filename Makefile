.PHONY: up
up:
	@docker-compose up -d

.PHONY: down
down:
	@docker-compose down

.PHONY: restart
restart:
ifeq ($(strip $(service)),)
	@echo "\e[31mInter the service name (make service={YOUR_SERVICE} restart)\e[39m"
else
	@docker-compose stop $(service)
	@docker-compose up -d $(service)
endif

.PHONY: list
list:
	@docker-compose ps

.PHONY: help
help:
	@echo "\nYou can use the following make commands:\n"
	@echo "make up - In order to up services"
	@echo "make down - In order to down services"
	@echo "make service={YOU_SERVICE} restart - In order to restart services"
	@echo "make list - In order to look through the service list\n"