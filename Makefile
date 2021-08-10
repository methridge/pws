.DEFAULT_GOAL := build

BINARY_NAME = pws
LD_FLAGS = "-X main.api=${WU_API_URL} \
-X main.sid=${WU_SID} \
-X main.units=${WU_UNITS} \
-X main.key=${WU_API_KEY}"

build:
	@go mod tidy
	@go build -o ${BINARY_NAME} -ldflags ${LD_FLAGS} pws.go

run:
	@go mod tidy
	@go build -o ${BINARY_NAME} -ldflags ${LD_FLAGS} pws.go
	./${BINARY_NAME}

install:
	mv ${BINARY_NAME} ${HOME}/bin

clean:
	go clean