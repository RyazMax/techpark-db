
## Documentation for API
Documentation for API is stated in [OpenAPI](https://ru.wikipedia.org/wiki/OpenAPI_%28%D1%81%D0%BF%D0%B5%D1%86%D0%B8%D1%84%D0%B8%D0%BA%D0%B0%D1%86%D0%B8%D1%8F%29): swagger.yml

Documantion can be read in : https://tech-db-forum.bozaro.ru/

To run docker use:

```
docker build -t forum .
docker run -p 5000:5000 --name forum -t forum
```

## Functional testing
API will be testes by automated functional testing

Method of testing:

 * build Docker;
 * run Docker;
 * run testing script in Go;
 * stop Docker.

Testing programs can be downloaded via:

 * [darwin_amd64.zip](https://bozaro.github.io/tech-db-forum/darwin_amd64.zip)
 * [linux_386.zip](https://bozaro.github.io/tech-db-forum/linux_386.zip)
 * [linux_amd64.zip](https://bozaro.github.io/tech-db-forum/linux_amd64.zip)
 * [windows_386.zip](https://bozaro.github.io/tech-db-forum/windows_386.zip)
 * [windows_amd64.zip](https://bozaro.github.io/tech-db-forum/windows_amd64.zip)

For local usage of Go testing script use:
```
go get -u -v github.com/bozaro/tech-db-forum
go build github.com/bozaro/tech-db-forum
```
After that `tech-db-forum` will be created.

### Run functional testing

To run functional testing use:
```
./tech-db-forum func -u http://localhost:5000/api -r report.html
```

Possible parammetres:

Parammetr                              | Description
---                                   | ---
-h, --help                            | Print possible parammetres
-u, --url[=http://localhost:5000/api] | set base URL of tested application
-k, --keep                            | Keep testing after first failed test
-t, --tests[=.*]                      | Mask of running tests (regular expression)
-r, --report[=report.html]            | Report file name
