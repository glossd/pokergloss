go_get_pokerblow_deps:
	echo "machine machine gitlab.com" > ~/.netrc &&  echo "  login $(username)@gitlab.com" >>  ~/.netrc && echo "  password $(token)" >> ~/.netrc
	go env -w GOPRIVATE=gitlab.com/pokerblow
	go get gitlab.com/pokerblow/go-auth
	go get gitlab.com/pokerblow/go-config


run:
	@make db_run
	PB_JWT_VERIFICATION_DISABLE=true PB_STDOUT_LEVEL=debug go run main.go

test:
	go test ./...

unit_test:
	go test $$(go list ./... | grep -v /e2e)

e2e_test:
	PB_STDOUT_LEVEL=debug go test ./e2e/...

lint:
	golangci-lint run


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


token ?= eyJhbGciOiJSUzI1NiIsImtpZCI6IjQ5YWQ5YmM1ZThlNDQ3OTNhMjEwOWI1NmUzNjFhMjNiNDE4ODA4NzUiLCJ0eXAiOiJKV1QifQ.eyJwaWN0dXJlIjoiaHR0cHM6Ly9zdG9yYWdlLmdvb2dsZWFwaXMuY29tL3Bva2VyYmxvdy1hdmF0YXJzLzJmUjVNWXlqcU1TTWd6THRka2RyS1prYzE4dDEiLCJ1c2VybmFtZSI6ImRlbmlzIiwiaXNzIjoiaHR0cHM6Ly9zZWN1cmV0b2tlbi5nb29nbGUuY29tL3Bva2VyYmxvdyIsImF1ZCI6InBva2VyYmxvdyIsImF1dGhfdGltZSI6MTU5OTE0NDc1OCwidXNlcl9pZCI6IjJmUjVNWXlqcU1TTWd6THRka2RyS1prYzE4dDEiLCJzdWIiOiIyZlI1TVl5anFNU01nekx0ZGtkcktaa2MxOHQxIiwiaWF0IjoxNTk5ODI0OTY2LCJleHAiOjE1OTk4Mjg1NjYsImVtYWlsIjoiZGVuaXNnbG90b3Y5OEBtYWlsLnJ1IiwiZW1haWxfdmVyaWZpZWQiOmZhbHNlLCJmaXJlYmFzZSI6eyJpZGVudGl0aWVzIjp7ImVtYWlsIjpbImRlbmlzZ2xvdG92OThAbWFpbC5ydSJdfSwic2lnbl9pbl9wcm92aWRlciI6InBhc3N3b3JkIn19.uUFnfxWFGXAlPFjS9Lf5QjB_YLx824dqDjnOQi-xYquaCRjE9yMHZ-Tk8MKLhbjjWWWK0eleghfqIbZQmTowXhOkq9h5fCb5RgWUrRe2otRm5APhWrX-xXkkDpl5uEcYj9xly-1MHrCUBQdUCFaQdkwD960LJf-8saFmv_FWND-QveRyJxja10Dt8L6I6YI_NjoZTgYe24JamNNgxkFrfqFNN6fPN_X_SQYJXHetkX0zm_SL6UODwgId6ioTSOM0QrLt1ZcdGiq1eNuX8s9aQqHs8XYPRm4Bt0zARX3J4RReX9OZa64O9d26WDnLVz_QdPzzIz8itxpiZnYIjnGVRA


generate_docs:
	docker run --rm -v $$PWD:/app -w /app pokerblow/swag:1.6.7 /root/swag init

generate_client:
	docker run --rm -v $$PWD:/local openapitools/openapi-generator-cli:v4.3.1 generate -i /local/docs/swagger.yaml -g typescript-axios -o /local/client-typescript

clean_client_gen:
	rm -rf client-typescript
	rm -rf docs

install_lint:
	GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.26.0
