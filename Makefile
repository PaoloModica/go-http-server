.PHONY: start-pg stop-pg test-integration

start-pg:
	docker-compose up --build --force-recreate --detach

stop-pg:
	docker-compose down --remove-orphans

test-integration:
	$(MAKE) start-pg
	@sleep 5
	go test -run server_integration_test.go
	$(MAKE) stop-pg
