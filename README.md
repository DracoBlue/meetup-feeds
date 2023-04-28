# meetup-feeds

This is an experimental little service makes meetup.com events fetched from the rest api available as feeds.

Add to `.env` file the keys to protect your service from unauthenticated access

```shell
USERNAME=asdf
PASSWORD=idk
```

this will limit the access to requiring the `Authentication` header:

```shell
$ curl --user asdf:idk http://0.0.0.0:8090/rss/example
```

and you will fetch an rss feed of the example group

If you need a specific port, you can add the environment variable for port:

```shell
PORT=8081
```

and it will be accessible via port 8081.

# License

Meetup-feeds is licensed under MIT License and Copyright by DracoBlue.
