go_get_pokerblow_deps:
	echo "machine machine gitlab.com" > ~/.netrc &&  echo "  login $(username)@gitlab.com" >>  ~/.netrc && echo "  password $(token)" >> ~/.netrc
	go env -w GOPRIVATE=gitlab.com/pokerblow
	go get gitlab.com/pokerblow/go-auth
	go get gitlab.com/pokerblow/go-config


run_pubsub_emulator:
	gcloud beta emulators pubsub start --project=pokerblow --host-port=localhost:8282

run:
	PB_PORT=4040 PB_JWT_VERIFICATION_DISABLE=true $$(gcloud beta emulators pubsub env-init) && go run main.go

test:
	go test ./...

token ?= eyJhbGciOiJSUzI1NiIsImtpZCI6IjUxMDM2YWYyZDgzOWE4NDJhZjQzY2VjZmJiZDU4YWYxYTc1OGVlYTIiLCJ0eXAiOiJKV1QifQ.eyJ1c2VybmFtZSI6ImRlbmlzIiwiaXNzIjoiaHR0cHM6Ly9zZWN1cmV0b2tlbi5nb29nbGUuY29tL3Bva2VyYmxvdyIsImF1ZCI6InBva2VyYmxvdyIsImF1dGhfdGltZSI6MTU5ODY0Njg3OCwidXNlcl9pZCI6IjJmUjVNWXlqcU1TTWd6THRka2RyS1prYzE4dDEiLCJzdWIiOiIyZlI1TVl5anFNU01nekx0ZGtkcktaa2MxOHQxIiwiaWF0IjoxNTk4NjQ2ODc4LCJleHAiOjE1OTg2NTA0NzgsImVtYWlsIjoiZGVuaXNnbG90b3Y5OEBtYWlsLnJ1IiwiZW1haWxfdmVyaWZpZWQiOmZhbHNlLCJmaXJlYmFzZSI6eyJpZGVudGl0aWVzIjp7ImVtYWlsIjpbImRlbmlzZ2xvdG92OThAbWFpbC5ydSJdfSwic2lnbl9pbl9wcm92aWRlciI6InBhc3N3b3JkIn19.Kcn9QroIR-62xGlzGTHdvx2uNRfiqAtUJyYBfzg74Mt_v4XozZTW-6O_teFmFoRasJsOr49uW4i9ntkkgoc6FgDoo1jTi_1yMMx3_gNS9qSAMMcmMscqelOHQdgxsi9mJMwltHqHOf-AsoYl7qbc_HCf5ShYBtljZlkUXY_pMGvy0ePupNMFiWxoTYmNiIelaz0d-O9oVzns8XOm6O6A5qsFLx6hnNNsS7cBMMbc9zqhsySAZhMzYHdkd-LvL8QtUUjAQQfsbit9hPFa4irFEf7gfYOX61kUeKlMBP6-f6O8Q8GmgieBJ30Ly8wQwJcgHcYkTAFlZhD8FRO08QesgQ
ws_user:
	@wscat -c "ws://localhost:8080/api/ws/news?token=$(token)"


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

lint:
	golangci-lint run

install_lint:
	GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.26.0
