protoc -I. \
--go_out=../server/pb \
--go_out=../client/pb  \
--go_opt=paths=source_relative \
--go-grpc_out=../server/pb \
--go-grpc_out=../client/pb \
--go-grpc_opt=paths=source_relative \
employee.proto person.proto