
precommit: ensure format addlicense test check
	@echo "ready to commit"

ensure:
	go mod verify
	go mod vendor

format:
	go get golang.org/x/tools/cmd/goimports
	find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	find . -type f -name '*.go' -not -path './vendor/*' -exec goimports -w "{}" +

test:
	go test -cover -race $(shell go list ./... | grep -v /vendor/)

addlicense:
	go get github.com/google/addlicense
	addlicense -c "Benjamin Borbe" -y 2019 -l bsd ./*.go ./cmd/create-mqtt-data/*.go

check: lint vet errcheck

lint:
	go get golang.org/x/lint/golint
	golint -min_confidence 1 $(shell go list ./... | grep -v /vendor/)

vet:
	go vet $(shell go list ./... | grep -v /vendor/)

errcheck:
	go get github.com/kisielk/errcheck
	errcheck -ignore '(Close|Write|Fprint)' $(shell go list ./... | grep -v /vendor/)

run:
	docker network create kafka || echo 'network already exists'
	docker-compose up -d
	docker-compose logs -f

ksqlcli:
	docker run -ti \
	--net=kafka \
	confluentinc/cp-ksql-cli:5.3.1 \
	http://ksql-server:8088
