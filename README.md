# Some insights on Maps in Golang

Examples used in the article https://hackernoon.com/some-insights-on-maps-in-golang-rm5v3ywh

The article covers the internal structure of the Golang maps implementation.
This implementation can be changed from version to version.
That's why this code works only on a limited number of Golang versions.

The easiest way to run them is with Docker:

```
docker run --rm -it -v $(pwd):/app golang:1.13.10 go run /app/random.go
```

```
docker run --rm -it -v $(pwd):/app golang:1.13.10 go run /app/grow.go
```

```
docker run --rm -it -v $(pwd):/app golang:1.13.10 go run /app/spread.go
```

```
docker run --rm -it -v $(pwd):/app golang:1.13.10 go run /app/presizes.go
```

```
docker run --rm -it -v $(pwd):/app golang:1.13.10 go run /app/delete.go
```

Each example in detail explained in the article. 