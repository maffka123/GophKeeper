# GophKeeper

Service store your data.

Has two components: client ans server.

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
