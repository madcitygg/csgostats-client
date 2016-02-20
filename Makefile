.PHONY: release

version = $(shell sed -n 's/.*version = "\([^"]*\)\s*"/\1/p' main.go)

help:
	@echo "release - build binaries for current version (${version})"

release:
	mkdir -p release/
	-rm release/*

	GOOS=windows GOARCH=386 go build -o release/csgostats-client.exe . && \
		cd release && \
		zip -X -o csgostats-client_${version}_windows_386.zip csgostats-client.exe && \
		rm csgostats-client.exe

	GOOS=windows GOARCH=amd64 go build -o release/csgostats-client.exe . && \
		cd release && \
		zip -X -o csgostats-client_${version}_windows_amd64.zip csgostats-client.exe && \
		rm csgostats-client.exe

	GOOS=darwin GOARCH=amd64 go build -o release/csgostats-client . && \
		cd release && \
		chmod +x csgostats-client && \
		zip -X -o csgostats-client_${version}_darwin_amd64.zip csgostats-client && \
		rm csgostats-client

	GOOS=linux GOARCH=amd64 go build -o release/csgostats-client . && \
		cd release && \
		chmod +x csgostats-client && \
		zip -X -o csgostats-client_${version}_linux_amd64.zip csgostats-client && \
		rm csgostats-client
