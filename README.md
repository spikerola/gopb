
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

