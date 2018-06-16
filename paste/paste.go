package paste

import (
        "fmt"
        "github.com/google/uuid"
        "crypto/sha256"
       )

type Paste struct {
    paste       []byte
    uuid        uuid.UUID
    hash        [32]byte
    shortHash   [4]byte
}

func (p *Paste) String() string {
    return fmt.Sprintf("%s : %s : %x : %x", string(p.paste[:10]), p.uuid, p.hash, p.shortHash)
}

func New(data []byte, private bool, timer int) (*Paste, bool) {
    uid     := uuid.New()
    sha     := sha256.Sum256(data)
    var short [4]byte
    copy(short[:], sha[28:])
    // save paste on db
    return &Paste{data, uid, sha, short}, false
}

