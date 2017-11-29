# go-api-task

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

`Start Postgres`

```sh
$ docker run --name postgres -d -p 5432:5432 --env 'DB_USER=postgres' --env 'DB_PASS=postgres' --env 'DB_NAME=api-auth' sameersbn/postgresql:9.6-2
```

`Start API`

```sh
$ go run main.go
```

`Build docker`

```sh
$ docker build --no-cache -t img-auth-go .
```

`Start docker`

```sh
$ docker run --net=host -it -p 3000:3000 img-auth-go
```
