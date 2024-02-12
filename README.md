# go-rest-shell

This is a small project for interviews. It's a REST API that allows you to execute shell commands.

## Prerequisites

- Go 1.22.0 or later

## Building the Project

1. Clone the repository:

```sh
git clone https://github.com/ersonp/go-rest-shell.git
```

2. Navigate to the project directory:
```sh
cd go-rest-shell
```

3. Build the project:
```sh
go build -o rest-shell ./cmd
```
This will create an executable named `rest-shell`.

## Running the Project

After building the project, you can run it with:
```sh
./rest-shell
```

Check the usage with the `-h` flag:
```sh
./rest-shell -h
Usage of ./rest-shell:
  -host string
        Host to run the server on (default "localhost")
  -port int
        Port to run the server on (default 8080)
```

## Using the API

Once the server is running, you can send HTTP requests to execute shell commands. For example, to execute the `ls` command, you could send a POST request with the command in the request body:
```sh
curl -X POST -H "Content-Type: application/json" -d '{"command":"ls"}' http://localhost:8080/api/cmd
```
or use Postman.
Replace `"ls"` with the command you want to execute.

## Testing the Project

To run the tests, use the go test command:
```sh
go test ./...
```
This will run all tests in the project.

## Docker

To build use:
```sh
docker build -t rest-shell .
```
Tu run use:
```sh
docker run -p 8080:8080 rest-shell
```