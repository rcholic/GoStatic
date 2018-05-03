package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/sirupsen/logrus"

	"github.com/julienschmidt/httprouter"
)

func index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := ps.ByName("name")
	fmt.Fprintf(w, "hello world from index: "+name)
}
func main() {
	// r := mux.NewRouter()
	r := httprouter.New()
	r.GET("/hello/:name", index)

	imagePaths := []string{"/Users/guoliangwang/Desktop/cool_me.jpeg", "/Users/guoliangwang/Downloads/128.png"}

	ex, err := os.Executable()
	if err != nil {
		logrus.Errorf("err in getting exec path: %s\n", ex)
		panic(err)
	}
	logrus.Infof("exec path is %s\n", ex)
	for idx, imgPath := range imagePaths {
		if rel, err := filepath.Rel(ex, imgPath); err != nil {
			logrus.Errorf("err in getting relative path: %v\n", err)
		} else {
			logrus.Infof("for image: %s, rel dir is %s\n", imgPath, filepath.Dir(rel))
			r.ServeFiles(fmt.Sprintf("/static%d/*filepath", idx), rice.MustFindBox(filepath.Dir(rel)).HTTPBox())
		}
	}

	// fs := rice.MustFindBox("../../Desktop").HTTPBox()
	// r.ServeFiles("/static/*filepath", fs)
	// fs2 := rice.MustFindBox("../").HTTPBox()
	// r.ServeFiles("/static2/*filepath", fs2)
	errc := make(chan error, 1)
	go func() {
		errc <- http.ListenAndServe(":8888", r)
	}()

	// sleep 30 seconds
	time.Sleep(30 * time.Second)
	r.ServeFiles("/static10/*filepath", rice.MustFindBox("../185").HTTPBox()) // does it work after starting http server?
	<-errc

}
