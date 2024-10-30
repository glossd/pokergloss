.DEFAULT_GOAL := run

go_get_pokerblow_deps:
	echo "machine machine gitlab.com" > ~/.netrc &&  echo "  login $(username)@gitlab.com" >>  ~/.netrc && echo "  password $(token)" >> ~/.netrc
	go env -w GOPRIVATE=gitlab.com/pokerblow
	go get -u gitlab.com/pokerblow/go-auth
	go get -u gitlab.com/pokerblow/go-config

go_get_pokerblow_deps_v2:
	git config --global url."git@gitlab.com:".insteadOf "https://gitlab.com/"

level ?= debug
prepare_run:
	@make db_run
#	@make bank_stub_run &
#	@make achieve_stub_run &

# PG_PROFILE=dev;PG_STDOUT_LEVEL=debug;PG_JWT_VERIFICATION_DISABLE=true;PUBSUB_EMULATOR_HOST=localhost:8282;PG_MULTI_FREEROLL_INTERVAL_MIN=1
run:
	@make prepare_run
	PG_PROFILE=dev PG_STDOUT_LEVEL=$(level) PG_JWT_VERIFICATION_DISABLE=true PUBSUB_EMULATOR_HOST=localhost:8282 PG_MULTI_FREEROLL_INTERVAL_MIN=1 go run main.go

run_scheduler:
	@make prepare_run
	PG_PROFILE=local PG_STDOUT_LEVEL=$(level) PG_TABLE_APP_TYPE=scheduler go run main.go

run_server:
	@make prepare_run
	PG_PROFILE=local PG_STDOUT_LEVEL=$(level) PG_TABLE_APP_TYPE=server go run main.go

test:
	go test ./...

unit_test:
	go test $$(go list ./... | grep -v /e2e)

e2e_test:
	PG_STDOUT_LEVEL=debug PG_PROFILE=e2e go test ./e2e/...

lint:
	golangci-lint run


bank_stub_run:
	@echo "Creating bank grpc stub"
	@go run cmd/bankstub/main.go
achieve_stub_run:
	@echo "Creating achieve grpc stub"
	@go run cmd/achievestub/main.go

CONTAINER_MONGO_NAME := test-mongo
db_run:
	@docker volume create test-mongo || true
	@docker top ${CONTAINER_MONGO_NAME} || docker run --rm -p 27017:27017 --name ${CONTAINER_MONGO_NAME} -v test-mongo:/data/db -d mongo:4.2
db_kill:
	@docker stop ${CONTAINER_MONGO_NAME}

remove_not_killed_db:
	@read -p "Are you sure? It will delete all containers and  volumes:" SURE; \
	test $$SURE == "yes" || exit 123;
	docker rm -f $$(docker ps -q -a)
	docker volume rm $$(docker volume ls -qf dangling=true)



token ?= eyJhbGciOiJSUzI1NiIsImtpZCI6IjUxMDM2YWYyZDgzOWE4NDJhZjQzY2VjZmJiZDU4YWYxYTc1OGVlYTIiLCJ0eXAiOiJKV1QifQ.eyJ1c2VybmFtZSI6ImRlbmlzIiwiaXNzIjoiaHR0cHM6Ly9zZWN1cmV0b2tlbi5nb29nbGUuY29tL3Bva2VyYmxvdyIsImF1ZCI6InBva2VyYmxvdyIsImF1dGhfdGltZSI6MTU5ODY0Njg3OCwidXNlcl9pZCI6IjJmUjVNWXlqcU1TTWd6THRka2RyS1prYzE4dDEiLCJzdWIiOiIyZlI1TVl5anFNU01nekx0ZGtkcktaa2MxOHQxIiwiaWF0IjoxNTk4NjQ2ODc4LCJleHAiOjE1OTg2NTA0NzgsImVtYWlsIjoiZGVuaXNnbG90b3Y5OEBtYWlsLnJ1IiwiZW1haWxfdmVyaWZpZWQiOmZhbHNlLCJmaXJlYmFzZSI6eyJpZGVudGl0aWVzIjp7ImVtYWlsIjpbImRlbmlzZ2xvdG92OThAbWFpbC5ydSJdfSwic2lnbl9pbl9wcm92aWRlciI6InBhc3N3b3JkIn19.Kcn9QroIR-62xGlzGTHdvx2uNRfiqAtUJyYBfzg74Mt_v4XozZTW-6O_teFmFoRasJsOr49uW4i9ntkkgoc6FgDoo1jTi_1yMMx3_gNS9qSAMMcmMscqelOHQdgxsi9mJMwltHqHOf-AsoYl7qbc_HCf5ShYBtljZlkUXY_pMGvy0ePupNMFiWxoTYmNiIelaz0d-O9oVzns8XOm6O6A5qsFLx6hnNNsS7cBMMbc9zqhsySAZhMzYHdkd-LvL8QtUUjAQQfsbit9hPFa4irFEf7gfYOX61kUeKlMBP6-f6O8Q8GmgieBJ30Ly8wQwJcgHcYkTAFlZhD8FRO08QesgQ

get_tables:
	@curl localhost:8080/api/table/tables -H "Authorization: Bearer $(token)"
post_table:
	@curl localhost:8080/api/table/tables -H "Authorization: Bearer $(token)"  -d '{"name":"my awesome table", "size":9, "bigBlind":2}'
delete_table:
	@curl -XDELETE localhost:8080/api/table/tables/$(id) -H "Authorization: Bearer $(token)"
reserve_seat:
	@curl localhost:8080/api/table/tables/$(id)/reserveSeat -H "Authorization: Bearer $(token)" -d '{"position":1}'
ws_table:
	@wscat -c "ws://localhost:8080/api/table/ws/tables/$(id)?token=$(token)"


generate_docs:
	docker run --rm -v $$PWD:/app -w /app pokerblow/swag:1.6.7 /root/swag init

generate_client:
	docker run --rm -v $$PWD:/local openapitools/openapi-generator-cli:v4.3.1 generate -i /local/docs/swagger.yaml -g typescript-axios -o /local/client-typescript

insert_package_json:
	@echo "{\"name\":\"@pokergloss/$(shell pwd | sed 's#.*/##')-client-typescript\",\"version\":\"0.0.1\",\"main\":\"dist/index.js\",\"types\":\"dist/index.d.ts\",\"files\":[\"dist\"],\"scripts\":{\"build\":\"tsc\",\"prepare\":\"npm run build\"},\"keywords\":[],\"author\":\"DennisGloss\",\"dependencies\":{\"axios\":\"^0.19.2\",\"url\":\"^0.11.0\"}}" > ./client-typescript/package.json

insert_ts_config:
	@echo '{"compilerOptions":{"declaration":true,"strictNullChecks":true,"target":"es5","outDir":"dist","module":"commonjs","moduleResolution":"node","sourceMap":true,"lib":["es2015","dom"],"rootDir":"."},"include":["."],"exclude":["node_modules","dist"]}' > ./client-typescript/tsconfig.json

npm_publish:
	make insert_package_json
	make insert_ts_config
	cd client-typescript && npm install && npm publish --access public

clean_client_gen:
	rm -rf client-typescript
	rm -rf docs

publish_client:
	make generate_docs
	make generate_client
	make npm_publish
	make clean_client_gen

install_lint:
	GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.26.0
