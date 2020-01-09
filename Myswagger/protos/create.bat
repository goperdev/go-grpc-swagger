protoc --go_out=plugins=grpc:. test.proto
protoc --grpc-gateway_out=logtostderr=true:. test.proto
protoc --swagger_out=logtostderr=true:. test.proto
pause
exit
