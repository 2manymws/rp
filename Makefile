export GO111MODULE=on

default: test

ci: depsdev test

test: cert
	cp go.mod testdata/go_test.mod
	go mod tidy -modfile=testdata/go_test.mod
	go test ./... -modfile=testdata/go_test.mod -coverprofile=coverage.out -covermode=count

lint:
	golangci-lint run ./...

cert:
	mkdir -p testdata
	rm -f testdata/*.pem testdata/*.srl
	# cacert of *.example.com
	openssl req -x509 -newkey rsa:4096 -days 365 -nodes -sha256 -keyout testdata/cakey.pem -out testdata/cacert.pem -subj "/C=UK/ST=Test State/L=Test Location/O=Test Org/OU=Test Unit/CN=*.example.com/emailAddress=k1lowxb@gmail.com"
	# a.example.com
	openssl req -newkey rsa:4096 -nodes -keyout testdata/a.example.com.key.pem -out testdata/a.example.com.csr.pem -subj "/C=JP/ST=Test State/L=Test Location/O=Test Org/OU=Test Unit/CN=a.example.com/emailAddress=k1lowxb@gmail.com"
	openssl x509 -req -sha256 -in testdata/a.example.com.csr.pem -days 60 -CA testdata/cacert.pem -CAkey testdata/cakey.pem -CAcreateserial -out testdata/a.example.com.cert.pem -extfile testdata/a.example.com.openssl.cnf
	openssl verify -CAfile testdata/cacert.pem testdata/a.example.com.cert.pem
	# b.example.com
	openssl req -newkey rsa:4096 -nodes -keyout testdata/b.example.com.key.pem -out testdata/b.example.com.csr.pem -subj "/C=JP/ST=Test State/L=Test Location/O=Test Org/OU=Test Unit/CN=b.example.com/emailAddress=k1lowxb@gmail.com"
	openssl x509 -req -sha256 -in testdata/b.example.com.csr.pem -days 60 -CA testdata/cacert.pem -CAkey testdata/cakey.pem -CAcreateserial -out testdata/b.example.com.cert.pem -extfile testdata/b.example.com.openssl.cnf
	openssl verify -CAfile testdata/cacert.pem testdata/b.example.com.cert.pem

depsdev:
	go install github.com/Songmu/ghch/cmd/ghch@latest
	go install github.com/Songmu/gocredits/cmd/gocredits@latest
	go install github.com/k1LoW/octocov-go-test-bench/cmd/octocov-go-test-bench@latest

prerelease:
	git pull origin main --tag
	go mod download
	ghch -w -N ${VER}
	gocredits -w .
	git add CHANGELOG.md CREDITS go.mod go.sum
	git commit -m'Bump up version number'
	git tag ${VER}

prerelease_for_tagpr: depsdev
	gocredits -w .
	git add CHANGELOG.md CREDITS go.mod go.sum

release:
	git push origin main --tag

benchmark: depsdev
	go mod tidy -modfile=testdata/go_test.mod
	go test -modfile=testdata/go_test.mod -bench . -benchmem -benchtime 10000x -run Benchmark | octocov-go-test-bench --tee > custom_metrics_bencmark.json

.PHONY: default test benchmark
