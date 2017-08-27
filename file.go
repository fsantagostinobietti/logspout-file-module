package file

import (
	"bytes"
	//"errors"
	"log"
	//"net"
	"os"
	//"reflect"
	"text/template"
	"time"
	"strconv"

	"github.com/gliderlabs/logspout/router"
)

//
// file route exaple:
//   file:///var/log/sample.log?maxfilesize=102400
//

func init() {
	router.AdapterFactories.Register(NewFileAdapter, "file")
}

// NewRawAdapter returns a configured raw.Adapter
func NewFileAdapter(route *router.Route) (router.LogAdapter, error) {
	// route.Address, route.Options
	
	// get 'filename' from route.Address
	filename := route.Address
	
	tmplStr := "{{.Data}}\n"
	tmpl, err := template.New("file").Parse(tmplStr)
	if err != nil {
		return nil, err
	}
	
	// default max size (100Mb)
	maxfilesize := 1024*1024*100
	if route.Options["maxfilesize"] != "" {
		szStr := route.Options["maxfilesize"]
		sz, err := strconv.Atoi(szStr)
		if err == nil {
		    maxfilesize = sz
		}
	}
	
	a := Adapter{
		route: route,
		filename:  filename,
		maxfilesize: maxfilesize,
		tmpl:  tmpl,
	}
	
	// rename if exists, otherwise create it
	err = a.Rotate()
	if err != nil {
	    return nil, err
	}
	return &a, nil
}

// Adapter is a simple adapter that streams log output to a connection without any templating
type Adapter struct {
	filename  string
	filesize  int
	maxfilesize   int
	fp  *os.File
	route *router.Route
	tmpl  *template.Template
}

// Stream sends log data to a connection
func (a *Adapter) Stream(logstream chan *router.Message) {
	for message := range logstream {
		buf := new(bytes.Buffer)
		err := a.tmpl.Execute(buf, message)
		if err != nil {
			log.Println("err:", err)
			return
		}
		//log.Println("debug:", buf.String())
		_, err = a.fp.Write(buf.Bytes())
		if err != nil {
			log.Println("err:", err)
		}
		
		// update file size
		a.filesize = a.filesize+len(buf.Bytes())
		
		// rotate file if size exceed max size 
		if a.filesize > a.maxfilesize {
		    a.Rotate()
		}
	}
}


// Perform the actual act of rotating and reopening file.
func (a *Adapter) Rotate() (err error) {
	// Close existing file if open
    if a.fp != nil {
        err = a.fp.Close()
        a.fp = nil
        if err != nil {
            return err
        }
    }
    // Rename dest file if it already exists
    _, err = os.Stat(a.filename)
    if err == nil {
        err = os.Rename(a.filename, a.filename+"."+time.Now().Format(time.RFC3339))
        if err != nil {
            return err
        }
    }
    // Create a file.
    a.fp, err = os.Create(a.filename)
    if err != nil {
        return err
    }
    a.filesize = 0
    return nil
}
