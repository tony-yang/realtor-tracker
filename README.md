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
docker run -v <absolute path to>/realtor-tracker:/go/src -p 9999:80 -it --rm realtor-tracker bash
```
