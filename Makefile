build: yacc
	go build -o bin/sqlpar

yacc:
	go generate
	rm y.output
