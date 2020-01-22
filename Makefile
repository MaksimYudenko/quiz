.PHONY: build
build:
	@echo '>>' build quiz
	@GOOS=linux GOARCH=amd64 go build -o ./build/quiz -i ./cmd/main.go
	@echo '>>>' task '"build"' completed successfully

.PHONY: image
image:
	@echo '>>' image quiz
	@docker build -t quiz ./build/ > /dev/null
	@echo '>>>' task '"image"' completed successfully

.PHONY: clean
clean:
	docker images -q -f "dangling=true" | xargs -I {} docker rmi {}
	@docker rmi -f "quiz"

.PHONY: start
start:
	@docker run -it quiz

.PHONY: all
all: build image start