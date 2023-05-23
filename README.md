# Net-cat

## Description

This project consists on recreating the NetCat in a Server-Client Architecture that can run in a server mode on a specified port listening for incoming connections, and it can be used in a client mode, trying to connect to a specified port and transmitting information to the server. NetCat, nc system command, is a command-line utility that reads and writes data across network connections using TCP or UDP.

## Usage: how to run

- Server

```bash
  go run .
  go run . $port
```

- Client

```bash
  nc localhost $port
```

## Author

- [@nzharylk](https://www.github.com/nzharylk)
