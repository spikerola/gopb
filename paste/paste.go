package paste

import (
        "fmt"
        "encoding/hex"
        "database/sql"
        "crypto/sha256"
        "github.com/google/uuid"
        _ "github.com/mattn/go-sqlite3"
       )

type Paste struct {
    Data        []byte
    Uuid        uuid.UUID
    Hash        [32]byte
    ShortHash   [4]byte
}

func (p Paste) String() string {
    return fmt.Sprintf("%s : %s : %x : %x", string(p.Data[:10]), p.Uuid, p.Hash, p.ShortHash)
}

func GetPaste(idstring string) ([]byte, error) {
    db, err := sql.Open("sqlite3", "./pastes.db")
    if err != nil {
        return nil, err
    }
    defer db.Close()

    var rows *sql.Rows
    var paste string

    if uid, err := uuid.Parse(idstring); err == nil { // we have the uuid
        rows, err = db.Query(fmt.Sprintf("SELECT data FROM paste WHERE uuid = '%s'", uid))
        if err != nil {
            return nil, err
        }
    } else if id, err := hex.DecodeString(idstring); err == nil { // we have the hex
        if len(id) == 32 { // complete
            rows, err = db.Query(fmt.Sprintf("SELECT data FROM paste WHERE hash = '%x'", id))
        } else { // short
            rows, err = db.Query(fmt.Sprintf("SELECT data FROM paste WHERE shortHash = '%x'", id))
        }
        if err != nil {
            return nil, err
        }
    } else {
        return nil, fmt.Errorf("not found")
    }


    if rows.Next() == false {
        return nil, fmt.Errorf("not found")
    }

    err = rows.Scan(&paste)
    if err != nil {
        return nil, err
    }

    return []byte(paste), nil
}

func UpdatePaste(possibleUuid []byte, data []byte) (error) {
    uid, err := uuid.Parse(string(possibleUuid))
    if err != nil {
        return err
    }

    if len(data) == 0 {
        return fmt.Errorf("no data")
    }

    db, err := sql.Open("sqlite3", "./pastes.db")
    if err != nil {
        return err
    }
    defer db.Close()

    stmt, err := db.Prepare("UPDATE paste SET data = ? WHERE uuid = ?")
    if err != nil {
        return err
    }

    res, err := stmt.Exec(data, fmt.Sprintf("%s", uid))
    if err != nil || res == nil {
        return err
    }

    return nil
}

func New(data []byte, private bool, timer int) (*Paste, error) {
    uid     := uuid.New()
    sha     := sha256.Sum256([]byte(fmt.Sprintf("%s", uid)))
    var short [4]byte
    copy(short[:], sha[28:])

    db, err := sql.Open("sqlite3", "./pastes.db")
    if err != nil {
        return nil, err
    }

    stmt, err := db.Prepare("INSERT INTO paste(uuid, data, hash, shortHash) values(?, ?, ?, ?)")
    if err != nil {
        return nil, err
    }

    res, err := stmt.Exec(fmt.Sprintf("%s", uid), data, fmt.Sprintf("%x", sha), fmt.Sprintf("%x", short))
    if err != nil || res == nil {
        return nil, err
    }

    db.Close()

    return &Paste{data, uid, sha, short}, nil
}

