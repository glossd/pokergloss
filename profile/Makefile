

go_get_pokerblow_deps:
	echo "machine machine gitlab.com" > ~/.netrc &&  echo "  login $(username)@gitlab.com" >>  ~/.netrc && echo "  password $(token)" >> ~/.netrc
	go env -w GOPRIVATE=gitlab.com/pokerblow
	go get -u gitlab.com/pokerblow/go-auth
	go get -u gitlab.com/pokerblow/go-config

run:
	@make db_run
	PB_PROFILE=local PB_STDOUT_LEVEL=debug PB_PROFILE=local go run main.go

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