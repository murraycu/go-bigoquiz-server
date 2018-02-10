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
Unfortunately, this project only builds when it is in a parent src/github.com/murraycu/ directory.
That is incredibly annoying, but hopefully one day the dep tool will not require it:
https://github.com/golang/dep/issues/911#issuecomment-318516931

For instance, first place the sources at a suitable path:

    $ mkdir -p ~/bigoquiz-build/github.com/murraycu/src
    $ cd ~/bigoquiz-build/github.com/murraycu/src
    $ git clone git@github.com:murraycu/go-bigoquiz-server.git

Then build the project:

    $ cd go-bigoquiz-server
    $ export GOPATH=~/bigoquiz-build/github.com/murraycu
    $ dep ensure
    $ go build

    ./config_use_prod.sh
    $ gcloud app deploy

### Running locally

    Change the configuration by applying the patch:
    $ patch -p1 < ./0001-debugging-Use-localhost.patch
      (Don't git push this)

    Use an appropriate oauth2 config file:
    (These cannot be added to the git repository.)
    ./config_use_local.sh

    Then start the local server:
    $ dev_appserver.py .

[1]: https://developers.google.com/appengine
[2]: https://golang.org
[3]: https://developers.google.com/appengine/docs/python/ndb/
[4]: https://developers.google.com/appengine/docs/python/users/
