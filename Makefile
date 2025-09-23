gen_auth:
	protoc -I proto proto/auth.proto --go_out=./gen/go/auth/ --go_opt=paths=source_relative --go-grpc_out=./gen/go/auth/ --go-grpc_opt=paths=source_relative
gen_relations:
	protoc -I proto proto/relations.proto --go_out=./gen/go/relations/ --go_opt=paths=source_relative --go-grpc_out=./gen/go/relations/ --go-grpc_opt=paths=source_relative
auth:
	go run cmd/auth/main.go

migrate_init:
	migrate -path=migrations/ -database "postgresql://yaro21:r2n7kzcp@localhost:54320/discord?sslmode=disable" -verbose up
	
migrate_test:
	go build ./cmd/migrator/main.go
	go run ./cmd/migrator/main.go --storage_path=./storage/auth.db --migrations_path=./tests/migrations --migrations_table=migrations_test