package main

import (
  "fmt"
  "io"
  "os"
  "net/http"
  "html/template"
  "strconv"
  "bytes"
  "log"
  "time"
  "strings"

//  "github.com/BobBurns/particle"
  "github.com/gorilla/mux"
  "github.com/tarm/serial"
)

var t *template.Template
var s *serial.Port
var serialData = ""

type PData struct {
  Name  string
  Data  int
}

func routeOutput (w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  q := vars["q"]

  fmt.Println(q)

  // handle button press
  if q != "" {


    switch(q) {
    case "on":
	    fmt.Println("got on")
	    b := []byte{0xb0, 0}
	    _, err := s.Write(b)
	    if err != nil {
		    log.Println("error on: ", err)
	    }
	    // set a timer
	    timer := time.NewTimer(10 * time.Minute)
	    go func() {
		    //wait for timer to expire, then redirect to off
		log.Println("starting timer")
		<-timer.C
		log.Println("ending timer")
	    	b := []byte{0xb1}
	    	_, err := s.Write(b)
	        if err != nil {
		    log.Println("error off: ", err)
	        }
	    }()

//      event.Data.Data = "on"
    case "off":
 //     event.Data.Data = "off"
	    b := []byte{0xb1}
	    _, err := s.Write(b)
	    if err != nil {
		    log.Println("error off: ", err)
	    }
    }
  }


  serialData = strings.Trim(serialData, "\r\n")
  intdata, _ := strconv.Atoi(serialData)

  d := PData{
    Name: "Moisture",
    Data: intdata,
  }
  var b bytes.Buffer

  h := "html-template.html"
  err := t.ExecuteTemplate (&b, h, d)

  if err != nil {
    fmt.Fprintf(w, "Error with template: %s ", err)
    return
  }
  b.WriteTo(w)
}


func main() {
  if len(os.Args) < 2 {
	  log.Fatal("Usage: sudo ./water <server IP address>")
  }
  ipaddr := os.Args[1]

  // 
  // parse html template
  t = template.Must(template.ParseFiles("html/html-template.html"))

  // subscribe to particle events

  c := &serial.Config{
	Name: "/dev/ttyUSB0",
	Baud: 9600,
	ReadTimeout: time.Second,
	}
  port, err := serial.OpenPort(c)
  if err != nil {
	log.Fatal(err)
  }
  s = port

  // start reading from serial port
  go func() {
  	for {
  	  buf := make([]byte, 128)
	  n, err := s.Read(buf)
	  if err != nil {
          // just print it
	    if err==io.EOF {
 	      // handle eof 
	      // try sleep and then flush
	      fmt.Println("got eof")
	      log.Fatal(err)
	  }
	}
	serialData = fmt.Sprintf("%s", buf[:n])
    }
  }()



	router := mux.NewRouter()
	/* change this to IP addr !! */
	sub := router.Host(ipaddr).Subrouter()
	sub.PathPrefix("/html/").Handler(http.StripPrefix("/html/", http.FileServer(http.Dir("html"))))
	sub.HandleFunc("/data", routeOutput)
	sub.HandleFunc("/data/{q}", routeOutput)

	// IdleTimeout requires go1.8
	server := http.Server{
		Addr:         ":8082",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      router,
	}
	fmt.Printf("Server started at %s:8082\n", ipaddr)
	log.Fatal(server.ListenAndServe())

}
