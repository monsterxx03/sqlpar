build: parser
	go build -o bin/sqlpar

parser:
	go generate
	rm y.output
