
compile_ws:
	protoc --proto_path=mqws --go_out=mqws --go_opt=paths=source_relative mqws/mqws.proto

compile_table:
	protoc --proto_path=mqtable --go_out=mqtable --go_opt=paths=source_relative mqtable/mqtable.proto

compile_profile:
	protoc --proto_path=mqprofile --go_out=mqprofile --go_opt=paths=source_relative mqprofile/mqprofile.proto

compile_bank:
	protoc -I . mqbank/mqbank.proto --go_out=plugins=grpc:.

compile_mail:
	protoc -I . mqmail/api.proto --go_out=plugins=grpc:.

compile_messenger:
	protoc -I . mqmessenger/api.proto --go_out=plugins=grpc:.

compile_market:
	protoc -I . mqmarket/api.proto --go_out=plugins=grpc:.

compile_survival:
	protoc -I . mqsurvival/api.proto --go_out=plugins=grpc:.
