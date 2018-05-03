### Multiple static paths served by Go HTTP Server from different directories

1. This is an experimentation with go http server serving static files/images from different directories on disk. I also tested adding a new static path after starting the server, and surprisingly it works as well :)

2. To run this project:
  - `go get`
  - `go run main.go -port 8888`
  - then visit `http://localhost:8888/hello/Mr/Tony` in your browser to see;
    you could see the static images at `http://localhost:8888/static0/galaxy.jpeg` and `http://localhost:8888/static1/apollo13.jpg`. These two images are served from the directories `images_path2` and `images_path3`, respectively. 


**TODO:** use websocket to have the web server push the static paths to frontend to update image srcattributes