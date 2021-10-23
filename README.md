Fizzbuzz Http Server
====================

A smal project for me to learn more practical go and also show what I can do when learning new things.

Quick documentation
-------------------
### Usage
```bash
$ go build
$ ./fizzbuzz-go
$ ./fizzbuz-go -host domin.tld # default localhost
$ ./fizzbuzz-go -port 8081 # default 8080
```

```bash
$ curl -X POST http://localhost:8080/fizzbuzz -d "{\"int1\": 3, \"int2\": 5, \"limit\": 100, \"str1\": \"fizz\", \"str2\": \"buzz\"}"
```

To post a new fizzbuzz request. Every parameter in the json body are optionnal and will take these values if not set.

```bash
$ curl http://localhost:8080/stats
```

Prints the stats and requests used (sorted by mostly used)
