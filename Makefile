lint:
	@ echo "-> Running linters..."
	@ go tool golangci-lint run  ./...
	@ echo "-> Done."

vuln:
	@ echo "-> Running vulnerability checks..."
	@ go tool govulncheck ./...
	@ echo "-> Done."

test: |
	@ echo "-> Running unit tests..."
	@ go test -v ./...
	@ echo "-> Done."
