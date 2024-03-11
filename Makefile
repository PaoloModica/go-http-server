.PHONY: start-pg stop-pg test-integration

start-pg:
	docker-compose up --build --force-recreate --detach

stop-pg:
	docker-compose down --remove-orphans

test:
	$(MAKE) start-pg
	@sleep 5
	go test
	$(MAKE) stop-pg
