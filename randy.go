package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strconv"
    "time"
    "hash"
    "crypto/md5"
    "crypto/sha256"
    "crypto/sha512"
)

// Randy represents runtime structure.
type Randy struct {
    HashAlgo       string
    Hostname       string
    Port           string
    RunCounter     int
    CurrentCounter int
}

var randy *Randy
var hashes = map[string]bool {
    "md5": true,
    "sha256": true,
    "sha512": true,
    "raw": true,
}

func main() {
    initRandy()

    log.Print(fmt.Sprintf("Randy initialized, running on port %s", randy.Port))

    http.HandleFunc("/", randyHandler)

    log.Fatal(http.ListenAndServe(":"+randy.Port, nil))
}

func randyHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, randy.ToString())
}

// ToString returns current Randy structure as a final ID string.
func (randy *Randy) ToString() string {
    uid := fmt.Sprintf("%d:%d:%d:%s:%s", randy.Timestamp(), randy.RunCounter, randy.Increment(), randy.Hostname, randy.Port)
    
    hash := randy.CreateHash()

    // If hash type is "raw" - just return generated uid.
    // You should never use "raw" in production as it will expose your hostnames!
    if hash == nil {
        return uid
    }

    hash.Write([]byte(uid))
    return fmt.Sprintf("%x", hash.Sum(nil))
}

// Increment increments internal Randy counter.
func (randy *Randy) Increment() int {
    randy.CurrentCounter++
    return randy.CurrentCounter
}

// Timestamp returns current datetime in unix format in nanoseconds.
func (randy *Randy) Timestamp() int64 {
    return time.Now().UnixNano()
}

// CreateHash creates object that will be used to hash ids
func (randy *Randy) CreateHash() hash.Hash {
    switch randy.HashAlgo {
    case "raw":
        return nil
    case "md5":
        return md5.New()
    case "sha256":
        return sha256.New()
    case "sha512":
        return sha512.New()
    }
    
    return nil
}

func initRandy() {
    portFlag := flag.String("port", "8080", "port on which Randy will listen")
    hashFlag := flag.String("hash", "raw", "function to hash ids")
    flag.Parse()

    if !hashes[*hashFlag] {
        log.Panic(fmt.Sprintf("\"%s\" is not a valid hash function", *hashFlag))
    }

    runCounter := loadRunCounter(*portFlag)
    hostname, err := os.Hostname()

    if err != nil {
        log.Panic(err)
    }

    randy = &Randy{
        HashAlgo:       *hashFlag,
        Hostname:       hostname,
        Port:           *portFlag,
        RunCounter:     runCounter,
        CurrentCounter: 0,
    }
}

// Randy stores counters of process initialization in counter files.
func loadRunCounter(port string) int {
    filename := "counter/" + port + ".cnt"

    var currentCounter int

    count, err := ioutil.ReadFile(filename)

    if err != nil {
        currentCounter = -1
    } else {
        currentCounter, _ = strconv.Atoi(string(count))
    }

    currentCounter = currentCounter + 1

    err = ioutil.WriteFile(filename, []byte(strconv.Itoa(currentCounter)), 0600)

    if err != nil {
        log.Panic(err)
    }

    return currentCounter
}
