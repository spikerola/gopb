# Go Pastebin

## Usage

upload
```bash
cmd | curl -F c=@- https://yourpbserver/
```

get
```
curl https://yourpbserver/pastehash > yourdata
```

get from browser
```
https://yourpbserver/pastehash/web
```

updating
```bash
curl -X PUT -F c=@- https://yourpbserver/paste-uuid < yourdata
```

deleting
```bash
curl -X DELTE https://yourpbserver/paste-uuid
```

## Installation

dependencies
```bash
go install github.com/google/uuid
go install github.com/mattn/go-sqlite3
```

compile
```bash
cd whatever/dir/gopb/
go build
```

create db
```bash
touch pastes.db
sqlite3 pastes.db "create table paste(uuid text primary key not null, data blob not null, hash text not null, shorthash text not null);"
```

run
```bash
./gopb &>>logs &
```

## TODO

- sunset - destroy paste after x seconds

