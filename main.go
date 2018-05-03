package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/sirupsen/logrus"

	"github.com/julienschmidt/httprouter"
)

type Model struct {
	Title string
	Name  string
}

var (
	templateMap = template.FuncMap{
		"Upper": func(s string) string {
			return strings.ToUpper(s)
		},
	}
	templates   = template.New("").Funcs(templateMap)
	templateBox *rice.Box
	port        string
)

func newTemplate(path string, _ os.FileInfo, _ error) error {
	if path == "" {
		return nil
	}
	templateString, err := templateBox.String(path)
	if err != nil {
		log.Panicf("Unable to extract: path=%s, err=%s", path, err)
	}
	if _, err = templates.New(filepath.Join("./templates", path)).Parse(templateString); err != nil {
		log.Panicf("Unable to parse: path=%s, err=%s", path, err)
	}
	return nil
}

// Render a template given a model
func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	err := templates.ExecuteTemplate(w, tmpl, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	title := ps.ByName("title")
	name := ps.ByName("name")
	model := Model{Title: title, Name: name}
	renderTemplate(w, "templates/index.html", &model)
}

func init() {
	flag.StringVar(&port, "port", "8888", "On what port number to run http server")
}

func main() {

	flag.Parse()

	templateBox = rice.MustFindBox("./templates")
	templateBox.Walk("", newTemplate)

	r := httprouter.New()
	r.GET("/hello/:title/:name", index)

	imagePaths := []string{"./images_path2/galaxy.jpeg", "./images_path3/apollo13.jpg"}
	pathMap := make(map[string]string)

	ex, err := os.Executable()
	if err != nil {
		logrus.Errorf("err in getting exec path: %s\n", ex)
		panic(err)
	}
	logrus.Infof("exec path is %s\n", ex)
	for _, imgPath := range imagePaths {
		absPath, er := filepath.Abs(imgPath)
		if er != nil {
			continue
		}
		if rel, err := filepath.Rel(ex, absPath); err != nil {
			logrus.Errorf("err in getting relative path: %v\n", err)
		} else {
			pathMap[rel] = imgPath
		}
	}

	idx := 0
	for k := range pathMap {
		r.ServeFiles(fmt.Sprintf("/static%d/*filepath", idx), rice.MustFindBox(filepath.Dir(k)).HTTPBox())
		idx++
	}

	errc := make(chan error, 1)
	go func() {
		logrus.Infof("http server listening on port %s\n", port)
		errc <- http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	}()

	// sleep 30 seconds before adding a new static path
	time.Sleep(30 * time.Second)

	// add new static path after starting the http server
	r.ServeFiles("/static/*filepath", rice.MustFindBox("./images").HTTPBox())
	<-errc
}
