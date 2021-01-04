# go-bigoquiz-server

This is the go-based backend server implementation for bigoquiz.com.

## Products Used

- [App Engine][1]

## Language

- [Go][2]

## APIs Used

- [NDB Datastore API][3]
- [Users API][4]

## Build / Deploy

    $ git clone git@github.com:murraycu/go-bigoquiz-server.git
    $ cd go-bigoquiz-server
    $ go build

    $ gcloud app deploy .

### Running locally

    Change the configuration by applying the patch:
    $ patch -p1 < ./0001-debugging-Use-localhost.patch
      (Don't git push this)

    Then start the local server:
    $ make local_run

[1]: https://developers.google.com/appengine
[2]: https://golang.org
[3]: https://developers.google.com/appengine/docs/python/ndb/
[4]: https://developers.google.com/appengine/docs/python/users/
