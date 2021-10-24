format:
	@echo "> Formatting the source"
	go fmt .


clean:
	@echo "> Cleaning build files and cache"
	go clean
	rm -rf tmp bin glin

.PHONY: format clean