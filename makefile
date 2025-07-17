PROTO_DIR=proto
OUT_DIR=.
PROTO_FILES=$(wildcard $(PROTO_DIR)/*.proto)

build:
	protoc -I. --go_out=$(OUT_DIR) --go-grpc_out=$(OUT_DIR) $(PROTO_FILES)

clean:
	rm -f $(OUT_DIR)/*.pb.go