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
 
