# Binary REST API server
## Description

This is a REST API server for finding integer values 
specified in input file. Input file is parsed into memory
and then search is done using binary search algorithm.

## Usage

App can be started by running
```bash
make run
```
or building and launching binary manually
```bash
make build # default output path is ./build/app
./build/app
```
```bash 
make build your/path/to/output # or override server binary output path 
# and then launch binary 
./your/path/to/output
```

After launch, the only endpoint is available at `http://localhost:8080/endpoint/{value}`.
For example,
```bash
curl http://localhost:8080/endpoint/123
```
Sample output:
```json
{"index":10,"value":1000}
```
```json
{"message":"index not found"}
```
```json
{"message":"invalid value passed"}
```

## Configuration

Configuration is done through dotenv file `.env` in the root of the project.
Example file is given at `.env.example`. 

```dotenv
PORT=9999 # number of port on which service is running, 8080 by default
LOG_LEVEL=DEBUG # one of following options - DEBUG, INFO, ERROR, fallback is INFO
INPUT_FILE=input.txt # input file name, input.txt by default
```

## Project structure

Project structure at top is divided to idiomatic `cmd` and `internal` folders,
where `cmd` contains main package, which is responsible for application launching 
and `internal` contains non-exported packages of server app.
`internal` is divided into `app`, `repo` and `api` folders.
- `app` contains main application logic, including server initialization, 
    configuration and routing.
- `repo` contains repository layer, which is responsible for domain logic. 
    In this case, it is only file parsing and binary search.
- `api` contains API layer, which is responsible for handling requests and 
    responses from outer world. In this case, it is only one endpoint handler.

Logging is done using newly introduced `log/slog` package. Global logger instance
is used over project, which is initialized during configuration parsing. 
Due to project size, just global logger is used, but in case of project growth, 
if separate logging is needed, it can be done by instantiating
new logger instance during struct creation and passing it as a dependency.
 
## Dependencies

There are three direct deps in this project:
- [chi](https://github.com/go-chi/chi) - lightweight router. Standard net/http handlers 
    does not provide convenient way of routing with path params, and writing it from scratch
    requires a lot of boilerplate code. Among analogs, chi is chosen since it does not 
    require any additional dependencies, is responsible only for routing and
    compatible with net/http API. Unlike go-gin or go-echo, which make own abstractions
    replacing net/http, chi is just a router on top of standard API.
- [chi/render](https://github.com/go-chi/render) - library for convenient JSON rendering. 
    Obviously was chosen since chi is used as a router, but still can be used separately.
- [koanf](https://github.com/knadh/koanf) - library for configuration parsing. 
    It was chosen over [joho/godotenv](https://github.com/joho/godotenv) and
    [spf13/viper](https://github.com/spf13/viper) because it is not bounded to single dotenv config
    parsing implementation and can be easily extended to parse other formats, in case config format
    changes. Also, in comparison to viper, config parsers are available as standalone packages,
    which is more convenient for loading only parsers which are used in project, without loading unnecessary code.