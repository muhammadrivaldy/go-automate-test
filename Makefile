build:
	docker image build -t service .

create-migration:
	@read -p "input the migration name: " name; \
	migrate create -ext sql -dir migrations $$name