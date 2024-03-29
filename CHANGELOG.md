# Release Notes

Most recent version is listed first.  


## v0.0.11
- Bugfix;direct dependencies in a require block with indirect deps should be well formatted: https://github.com/komuw/ote/pull/58

## v0.0.10
- Recognize toolchain keyword in go.mod files: https://github.com/komuw/ote/pull/50
- Exit with code 0 on success: https://github.com/komuw/ote/pull/51
- Format indirect dependencies into their own block: https://github.com/komuw/ote/pull/53
- Use generic slice functions: https://github.com/komuw/ote/pull/56

## v0.0.9
- Update to Go v1.20: https://github.com/komuw/ote/pull/47

## v0.0.8
- Bugfix, do not let build tags impact the source files that are scanned: https://github.com/komuw/ote/pull/41
- Bugfix, check whether modules are in both test and non-test files
  previsously, we were checking whether import paths(as opposed to module paths) are in both test & non-test files: https://github.com/komuw/ote/pull/43


## v0.0.7
- List test-only dependencies in their own require block: https://github.com/komuw/ote/pull/37
- The require block of test-only dependencies should only be inserted after the last require block: https://github.com/komuw/ote/pull/39


## v0.0.6
- Remove unnecessary loop, so as to improve perfomance: https://github.com/komuw/ote/pull/33


## v0.0.5
- Add version in cli: https://github.com/komuw/ote/pull/34


## v0.0.4
- rewrite how `ote` is implemented: https://github.com/komuw/ote/pull/18
  Get all the golang files in a project, parse them, get imports
- Work with nested Go modules: https://github.com/komuw/ote/pull/17
  In repositories that have nested Go modules; ignore all the nested modules and only   
  work/analyze the module that is passed in as an argument to `ote`
- Have all static analysis passes succeed: https://github.com/komuw/ote/pull/20
- dont analyze Go files inside vendor/ directory: https://github.com/komuw/ote/pull/21
- Render unformatted `//test` comment correctly: https://github.com/komuw/ote/pull/24
- Perf improvement, do not generate a new string: https://github.com/komuw/ote/pull/25
  This is an improvement of https://github.com/komuw/ote/pull/24
- Call `fetchModule` only for the import paths that are not shared between test files and non-test files: https://github.com/komuw/ote/pull/26
- Add more tests: https://github.com/komuw/ote/pull/29
- Run tests in parallel: https://github.com/komuw/ote/pull/30
- If `ote` is unable to fetch the module for an import path due to build tags, emit appropriate error : https://github.com/komuw/ote/pull/31


## v0.0.3
-  take into account files ending in _test.go: https://github.com/komuw/ote/pull/14


## v0.0.2
- add some quality metrics: https://github.com/komuw/ote/pull/9
- add ability to remove test comments: https://github.com/komuw/ote/pull/10


## v0.0.1
- Add CLI: https://github.com/komuw/ote/pull/1
- Add documentation: https://github.com/komuw/ote/pull/2
- Add e2e tests: https://github.com/komuw/ote/pull/4
- add dummy run in CI: https://github.com/komuw/ote/pull/5
- add versioning script : https://github.com/komuw/ote/pull/6
