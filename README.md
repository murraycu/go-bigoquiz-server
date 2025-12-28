# go-bigoquiz-server

This is the go-based backend server implementation for bigoquiz.com.

## Products Used

- [App Engine][1]

## Language

- [Go][2]

## APIs Used

- [NDB Datastore API][3]
- [Users API][4]

## Build

    $ git clone git@github.com:murraycu/go-bigoquiz-server.git
    $ cd go-bigoquiz-server
    $ go build

Also available via "make build".

## Deploy

    $ gcloud app deploy .

Also available via "make deploy".

But prefer the [GitHub Deployment
workflow](https://github.com/murraycu/go-bigoquiz-server/actions/workflows/deploy_to_prod.yaml),
via the GitHub "Actions" tab.

Then see the [deployed versions in
AppEngine](https://console.cloud.google.com/appengine/versions?serviceId=api&project=bigoquiz).

### Running locally

    Start the local server:
    $ make local_run

[1]: https://developers.google.com/appengine
[2]: https://golang.org
[3]: https://developers.google.com/appengine/docs/python/ndb/
[4]: https://developers.google.com/appengine/docs/python/users/
