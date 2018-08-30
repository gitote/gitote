Execute following command in ROOT directory when anything is changed:

```sh
$ make bindata

or

$ ./go-bindata -o=$@ -ignore="\\.DS_Store|README.md|TRANSLATORS|auth.d" -pkg=bindata conf/...
```