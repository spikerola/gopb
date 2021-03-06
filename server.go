package main

import (
        "fmt"
        "html"
        "net/http"
        "log"
        "io/ioutil"
        "strconv"
        "gopb/paste"
)

func hello(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "gopb is running\n")
}

func handleGetRequest(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Path[:]
    if len(id) < 8 {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    id = id[1:]
    web := id[len(id)-4:] == "/web"
    if web {
        id = id[:len(id)-4]
    }
    p, err := paste.Get([]byte(id))

    if err != nil {
        log.Println("Error has occurred getting", id, ":", err)
        w.WriteHeader(http.StatusNotFound)
        return
    }

    if web {
        p = []byte(fmt.Sprintf("<!DOCTYPE html><html><head><meta charset=\"utf-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1, shrink-to-fit=no\"></head><body><pre>%s</pre></body></html>", html.EscapeString(string(p))))
    }

    fmt.Fprintf(w, "%s", p)
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
    file, _/*fileHeader*/, err := r.FormFile("c")
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    b, err := ioutil.ReadAll(file)
    if err != nil {
        log.Println("an error has occurred reading the file:", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    // save pb
    p := len(r.Form["p"])
    t := -1
    if len(r.Form["t"]) > 0 {
        if t2, err := strconv.Atoi(r.Form["t"][0]); err == nil {
            t = t2
        }
    }
    pst, err := paste.New(b, p == 1, t)
    if err != nil {
        log.Println("an error has occurred saving the file:", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    // continue with the paste
    fmt.Fprintf(w, "uuid: %s\nlong: %x\nshort: %x\n", pst.Uuid, pst.Hash, pst.ShortHash)
}

func handlePutRequest(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Path[:]
    if len(id) < 8 {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    id = id[1:]

    file, _/*fileHeader*/, err := r.FormFile("c")
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    data, err := ioutil.ReadAll(file)
    if err != nil {
        log.Println("an error has occurred reading a paste:", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    err = paste.Update([]byte(id), data)

    if err != nil {
        log.Println("an error has occurred updating a paste:", err)
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    fmt.Fprintf(w, "ok\n")
}

func handleDeleteRequest(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Path[:]
    if len(id) < 8 {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    id = id[1:]

    err := paste.Delete([]byte(id))

    if err != nil {
        log.Println("an error has occurred deleting a paste:", err)
    }

    fmt.Fprintf(w, "ok\n")
}

type httpRouter map[string]func(http.ResponseWriter, *http.Request)

var mux httpRouter

func main() {
    server := http.Server{
        Addr:    ":8000",
        Handler: &hndlr{},
    }

    mux = make(httpRouter)
    mux["/"] = hello

    server.ListenAndServe()
}

type hndlr struct{}

func (*hndlr) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        if h, ok := mux[r.URL.String()]; ok {
            h(w, r)
        } else {
            handleGetRequest(w, r)
        }
    case "POST":
        handlePostRequest(w, r)
        break
    case "PUT":
        handlePutRequest(w, r)
        break
    case "DELETE":
        handleDeleteRequest(w, r)
        break
    default:
        w.WriteHeader(http.StatusBadRequest)
    }
}



