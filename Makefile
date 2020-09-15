build: yacc
	go build -o bin/sqlpar

yacc:
	go generate
	rm y.output

test:
	go test --race -v $$(go list ./...| grep -v -e /vendor/)
