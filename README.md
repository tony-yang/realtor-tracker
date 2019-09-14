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
Run `make test` and `cd build && make` to test and format the code before commit.

To run a long-live dev environment
```
docker run -v <absolute path to>/realtor-tracker:/go/src/github.com/<org>/realtor-tracker -p 9999:80 -itd --rm realtor-tracker bash
```

To develop in the long-live dev container, make sure to add any new dependency in the source code, and run the `build` directory's Makefile to update the go module and the workspace. Then run the build script that updates bazel and protos. Finally, run the application for manual testing.
```
cd realtor-tracker/build
make
./update-bazel.sh
./update-protos.sh

# Run the application for manual testing
```

Note: The build Makefile will automatically run the `go mod tidy` to clean up any unused go package from the go module.
