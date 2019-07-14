DEFAULT=all

all: update-vendor update-gofmt

update-vendor:
	rm -rf ../vendor
	go mod tidy
	go mod vendor
	bazel run //:gazelle -- update-repos -from_file=go.mod

update-gofmt:
	cd .. && ls | grep -v vendor | xargs -I {} find {} -name "*.go" | xargs -I {} gofmt -s -w {}
