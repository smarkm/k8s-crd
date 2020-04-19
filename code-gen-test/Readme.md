# code-generator test

## Init API
- you need manuly init `doc.go`, `types.go`,`register.go` in versioned folder under `apis` and `register.go` under `group` folder, please ref the `Makefile init:`
- init a `header` file will be used by `gen` phase

## generate code
```bash
make gen
```
