.PHONY:  default  deploy  refresh  runserver  test  test-coverage  test-docker  test-release

default: test

deploy:
	GOOS=linux GOARCH=amd64 go build -o build/lambda-handler cmd/lambda-handler/main.go
	cd build && zip -9 lambda-handler.zip lambda-handler
	cd ../terraform/ && terraform apply -auto-approve

refresh:
	cookiecutter gh:sjansen/cookiecutter-golang --output-dir .. --config-file .cookiecutter.yaml --no-input --overwrite-if-exists
	git checkout go.mod go.sum

runserver:
	go run cmd/localdev/main.go -b ../webui/build/

test:
	@scripts/run-all-tests
	@echo ========================================
	@git grep TODO  -- '**.go' || true
	@git grep FIXME -- '**.go' || true

test-coverage: test-docker
	go tool cover -html=dist/coverage.txt

test-docker:
	@scripts/docker-up-test

test-release:
	git stash -u -k
	goreleaser release --rm-dist --skip-publish
	-git stash pop
