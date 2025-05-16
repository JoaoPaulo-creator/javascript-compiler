build:
	go run main.go input2.js
	clang output.ll -o output
	./output
