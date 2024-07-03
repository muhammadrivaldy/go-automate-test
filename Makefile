build:
	docker image build -t go-automate-test .

create-migration:
	@read -p "input the migration name: " name; \
	migrate create -ext sql -dir migrations $$name

run-integration:
	docker compose --profile integration-tests up --abort-on-container-exit --renew-anon-volumes --force-recreate --exit-code-from integration-tests