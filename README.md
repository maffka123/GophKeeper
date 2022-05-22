# GophKeeper

Service store your data.

Has two components: client ans server.

### Table of Contents

1. [Client](#client)
    1. [Build client](#build-client)
2. [Server](#server)
3. [Requests for tests](#requests-for-tests)

## Client

API with the followind endpoints:

* **/api/user/register**

    request:

```json
{"login": ..., "password": ...}
```

    response: cookie with jwt token

* **/api/user/login**

    request:

```json
{"login": ..., "password": ...}
```

    response: cookie with jwt token

* **/api/user/insert**

    request:

```json
{"data": 
    {"id": <int>, 
    "data": {
        "authdata": {
            "login": ..., 
            "password": ...}, 
        "text": ..., 
        "bytes": ..., 
        "bankcard": {
            "cardnumber": ..., 
            "expiry": ..., 
            "holdername": ..., 
            "address": ..., 
            "bankname": ...}}, 
    "metadata": <text>
    }}
```

    response: id of the new data entry

* **/api/user/search**

    request: same data structure as for insert, just with fileds that person remembers

    response: full data entry/entries
* **/api/user/delete**

    request: same data structure as for insert, just with fileds that person remembers

    response: full data entry/entries that were deleted

### Build client

```bash
GOOS=darwin 
GOARCH=amd64
go build -ldflags "-X 'main.BuildCommit=$(git rev-list -1 HEAD)' -X 'main.BuildDate=$(date)' -X 'main.Version=1.0'" -o client'-'$GOOS'-'$GOARCH
```

```bash
GOOS=windows 
GOARCH=amd64
go build -ldflags "-X 'main.BuildCommit=$(git rev-list -1 HEAD)' -X 'main.BuildDate=$(date)' -X 'main.Version=1.0'" -o client'-'$GOOS'-'$GOARCH'.exe'
```

```bash
GOOS=linux 
GOARCH=amd64
go build -ldflags "-X 'main.BuildCommit=$(git rev-list -1 HEAD)' -X 'main.BuildDate=$(date)' -X 'main.Version=1.0'" -o client'-'$GOOS'-'$GOARCH
```

## Server

It is a gRPC server that has the same endpoints as above

## Requests for tests

```bash
curl -d '{"login":"test1","password":"mypass"}' -H "Content-Type: application/json" -X POST http://localhost:8081/api/user/register
```

```bash
curl -d '{"login":"test1","password":"mypass"}' -H "Content-Type: application/json" -X POST http://localhost:8081/api/user/login
```

```bash
curl -v --cookie "jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyfQ.XUjieZQLFHd61t9ZjifbQ6c1BGB6ANYD1Xo-aog249U" -d '{"Data": {"BankCard":{"CardNumber": 123456789}}, "Metadata": "this is my card"}' -H "Content-Type: application/json" -X POST http://localhost:8081/api/user/insert
```

```bash
curl -v --cookie "jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyfQ.XUjieZQLFHd61t9ZjifbQ6c1BGB6ANYD1Xo-aog249U" -d '{"Data": {"AuthData":{"login": "login1", "password": "pass1"}}, "Metadata": "this is my login"}' -H "Content-Type: application/json" -X POST http://localhost:8081/api/user/insert
```

```bash
curl -v --cookie "jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyfQ.XUjieZQLFHd61t9ZjifbQ6c1BGB6ANYD1Xo-aog249U" -d '{"Metadata": "this is my card"}' -H "Content-Type: application/json" -X GET http://localhost:8081/api/user/search
```

```bash
curl -v --cookie "jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyfQ.XUjieZQLFHd61t9ZjifbQ6c1BGB6ANYD1Xo-aog249U" -d '{"ID": 1, "Metadata": "this is my card"}' -H "Content-Type: application/json" -X GET http://localhost:8081/api/user/search
```

```bash
curl -v --cookie "jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyfQ.XUjieZQLFHd61t9ZjifbQ6c1BGB6ANYD1Xo-aog249U" -d '{"Data": {"AuthData":{"login": "login1"}}}' -H "Content-Type: application/json" -X GET http://localhost:8081/api/user/search
```
