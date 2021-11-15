# azcosmos-example

## Setup

Setup environment variables
```bash
export DOCUMENT_DB_URI="https://<DB NAME>.documents.azure.com:443/"
export DOCUMENT_DB_PRIMARY_KEY="<PRIMARY KEY>"
```

Modify dbName/dbContainer to match database/container names. 

```bash
go build main.go
```

## Usage

Create item
```bash
./main -id 1 -create "test val"
```

Read item
```bash
./main -id 1 -read
```

Replace item
```bash
./main -id 1 -replace "replaced test val"
```

Delete item
```bash
./main -id 1 -delete
```