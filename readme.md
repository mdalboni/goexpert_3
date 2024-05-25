# GO Expert - Deploy com Cloud Run

You can use the live demo on: https://goexpert-3-icvwowoova-uc.a.run.app

## Run locally

In the project root execute:
```shell
docker-compose up
```

## APIs

### GET /weather/{zip_code}

200:
```json
{"temp_c":16,"temp_f":60.8,"temp_k":289.1}
```

422:
```
invalid zipcode
```

404:
```
can not find zipcode
```

## Run tests

go tests ./...