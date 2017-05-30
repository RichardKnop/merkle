fmt:
	bash -c 'go list ./... | grep -v vendor | xargs -n1 go fmt'

test:
	bash -c 'go list ./... | grep -v vendor | xargs -n1 go test -timeout=10s'
