The program parses transactions from transactions.txt file which is included in the repository and outputs the
prioritized transactions into prioritized-transactions.txt file, which will be created under the rood directory.
The paths of files are hardcoded in main.go

The following command compiles the program:
```
make build
```

The program can be executed by running the following command from the root directory:
```
bin/mempool
```

The following command runs the unit tests:
```
make test
```
