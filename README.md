# realtor-tracker
A web app for tracking real estate price change and market trend

## Architecture
The application contains 3 components: indexer, analyzer, and webmvc. All components will be developed using Golang.

### Indexer
An independent service that runs periodically to fetch real estate information from the web. It will provide a RESTful endpoint to serve data needed by the other components.

### Analyzer
An offline data analysis component that process some basic statistical analysis on the real estate data.

### Webmvc
A Golang MVC framework used to serve the information.

## Dev Guide
Run `make test` to test the code before commit.

To run a long-live dev environment
```
docker run -v <absolute path to>/realtor-tracker:/go/src/github.com/<org>/realtor-tracker -p 9999:80 -it --rm realtor-tracker bash
```

To develop in the long-live dev container, make sure to add any new dependency in go.mod of the corresponding directory manually or by running the go command, and then run bazel.
```
cd <the directory to test>
go test ./ # This automatically adds new dependency to go.mod
# or manually add new dependency to go.mod
bazel run //:gazelle -- update-repos -from_file=go.mod
```

Also remember to run `go mod tidy` to clean up.
