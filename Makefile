.PHONY: build clean deploy

build:
	GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/get_statements ./cmd/aws_lambda/get_statements

clean:
	rm -rf ./bin

deploy: clean build
	# sls deploy --verbose
	sls deploy function -f get_statements_list
