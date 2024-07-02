build:
	docker image build -t go-automate-test .

create-migration:
	@read -p "input the migration name: " name; \
	migrate create -ext sql -dir migrations $$name