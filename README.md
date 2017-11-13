# go-api-auth

`Install dep`

```sh
$ go get -v github.com/golang/dep/cmd/dep
```

`Install the project's dependencies`

```sh
$ dep ensure
```

`Build docker API`

```sh
$ CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/api
```

`Start API`

```sh
$ go run main.go
```

```sh
$ docker build --no-cache -t img-auth-go .
```