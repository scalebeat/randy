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
)

// Randy represents runtime structure.
type Randy struct {
    HashFunction   string
    Hostname       string
    Port           string
    RunCounter     int
    CurrentCounter int
}

var randy *Randy

func main() {
    initRandy()

    log.Print(fmt.Sprintf("Randy initialized, running on port %s", randy.Port))

    http.HandleFunc("/", countHandler)

    log.Fatal(http.ListenAndServe(":"+randy.Port, nil))
}

func countHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, randy.ToString())
}

// ToString returns current Randy structure as a final ID string.
func (randy *Randy) ToString() string {
    uid := fmt.Sprintf("%s:%d:%d:%d", randy.Hostname, randy.RunCounter, randy.Timestamp(), randy.Increment())
    return uid
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

func initRandy() {
    port := flag.String("port", "8080", "port on which Randy will listen")
    flag.Parse()

    runCounter := loadRunCounter(*port)
    hostname, err := os.Hostname()

    if err != nil {
        log.Panic(err)
    }

    randy = &Randy{
        Hostname:       hostname,
        Port:           *port,
        RunCounter:     runCounter,
        CurrentCounter: 0,
    }
}

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
