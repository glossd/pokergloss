
compile_achieve:
	protoc -I . grpcachieve/api.proto --go_out=plugins=grpc:.

compile_bank:
	protoc -I . grpcbank/api.proto --go_out=plugins=grpc:.

compile_market:
	protoc -I . grpcmarket/api.proto --go_out=plugins=grpc:.

compile_profile:
	protoc -I . grpcprofile/api.proto --go_out=plugins=grpc:.

compile_table_events:
	protoc -I . grpctableevents/api.proto --go_out=plugins=grpc:.

compile_table:
	protoc -I . grpctable/api.proto --go_out=plugins=grpc:.


docker_build_achieve:
	docker build -t registry.gitlab.com/pokerblow/go-grpc:stub-achieve-0.1 -f achieve.Dockerfile .
	docker push registry.gitlab.com/pokerblow/go-grpc:stub-achieve-0.1

docker_build_bank:
	docker build -t registry.gitlab.com/pokerblow/go-grpc:stub-bank-0.1 -f bank.Dockerfile .
	docker push registry.gitlab.com/pokerblow/go-grpc:stub-bank-0.1

docker_build_market:
	docker build -t registry.gitlab.com/pokerblow/go-grpc:stub-market-0.1 -f market.Dockerfile .
	docker push registry.gitlab.com/pokerblow/go-grpc:stub-market-0.1
