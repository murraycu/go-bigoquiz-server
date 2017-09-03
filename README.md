# go-bigquiz-server

This is a go-based server implementation for bigoquiz.com.
It is not yet live at bigoquiz.com.

## Products
- [App Engine][1]

## Language
- [Go][2]

## APIs
- [NDB Datastore API][3]
- [Users API][4]

## Build / Deploy

(TODO: Restructure/reconfigure the project so we don't need to specify bigoquiz and src/bigoquiz in these commands.)

    $ git clone git@github.com:murraycu/go-bigoquiz-server.git
    $ cd go-bigoquiz-server
    $ export GOPATH=`pwd`
    $ go get github.com/julienschmidt/httprouter
    $ go build bigoquiz

    $ gcloud app deploy src/bigoquiz

### Running locally

    $ dev_appserver.py app.yaml

[1]: https://developers.google.com/appengine
[2]: https://golang.org
[3]: https://developers.google.com/appengine/docs/python/ndb/
[4]: https://developers.google.com/appengine/docs/python/users/
