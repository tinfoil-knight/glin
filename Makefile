default:
	@echo "Specify a command to run with make"

codegen:
	@echo "> Generating code"
	stringer -type=TokenType
	go run tool/generate_ast.go .
	make format

format:
	@echo "> Formatting the source"
	go fmt .

clean:
	@echo "> Cleaning build files and cache"
	go clean
	rm -rf tmp bin glin

run:
	@echo "> Starting the CLI"
	nodemon -e go --signal SIGINT --exec 'go' run .

.PHONY: codegen format clean run