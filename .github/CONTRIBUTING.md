## To contribute:             

- fork this repo.
- make the changes you want on your fork.
- your changes should have backward compatibility in mind unless it is impossible to do so.
- add release notes to CHANGELOG.md
- add tests and benchmarks
- format your code using gofmt:                                          
- run tests(with race flag) and make sure everything is passing:
```shell
 go test -race -cover -v ./...
```
- run benchmarks and make sure that they havent regressed. If you have introduced any regressions, fix them unless it is impossible to do so:
```shell
go test -race -run=XXXX -bench=. ./...
```
- open a pull request on this repo.          
          
NB: I make no commitment of accepting your pull requests.                 
