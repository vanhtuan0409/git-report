build:
	GOOS=darwin GOARCH=amd64 go build -o bin/greport github.com/vanhtuan0409/git-report/cmd/greport

release:
	make build
	tar -C bin -zcvf dist/greport-macos.tar.gz greport
