package main

/* see ./paste.go */

import (
        "io"
        "net/http"
        "log"
        "io/ioutil"
        "strconv"
        "gopb/paste"
)

/* to run use
   $ go run *.go */

/* links opened
https://golang.org/pkg/crypto/sha256/
https://stackoverflow.com/questions/28933687/golang-random-sha256#28933817
https://golang.org/pkg/net/http/
*/

func hello(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "gopb is running\n")
}

func handleGetRequest(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "not implemented yet\n")
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
    file, _/*fileHeader*/, err := r.FormFile("c")
    if err != nil {
        io.WriteString(w, "\n")
        return
    }
    b, err := ioutil.ReadAll(file)
    if err != nil {
        log.Println(err)
        io.WriteString(w, "an error has occurred reading the file\n")
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
    pst, ok := paste.New(b, p == 1, t)
    log.Println(pst)
    if !ok {
        io.WriteString(w, "an error has occurred saving the file\n")
        return
    }
    // continue with the paste
    io.WriteString(w, "ok\n")
}

func handlePutRequest(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "not implemented yet\n")
}

func handleDeleteRequest(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "not implemented yet\n")
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
            io.WriteString(w, "is this a 404?\n")
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



