# garchive

Ultra-simple front-end for a web archive.

[Demo](https://archive.gpanders.com)

## Description

This tool provides a web front-end for a repository of archived web pages,
typically generated via `wget`. For example, you can use `wget` to mirror a
website using

    $ wget -kmpc https://gpanders.com

This will create a folder `gpanders.com` in your current directory with a copy
of `gpanders.com` that you can open locally. This is a very simple way to
maintain archives of websites and fight [link rot][].

## Installation

The easiest method is to use [Docker](#using-docker). If you can't or don't
want to use Docker, you can build `garchive` yourself:

    $ git clone git://git.gpanders.com/garchive.git
    $ cd garchive
    $ go build

## Usage

You must have a tab-delimited CSV file with each line containing a link
title/description and the corresponding URL, e.g.

    Personal website of Greg Anders	https://gpanders.com
                                   ^ literal tab character (\t)

A tab character (`\t`) is used as the delimiter as it can be reasonably
expected not to appear in the title of any webpage (which cannot be said for
more traditional CSV delimiters such as commas).

You can use the included [`fetch`][fetch] utility to parse this CSV file and
clone all of the URLs into a local directory:

    $ bin/fetch links.csv data

This will use `wget` to archive all of the URLs in `links.csv` under the `data`
directory.

Once you have your `links.csv` file and your archived websites, use `garchive`
to serve up a simple front-end to access those archives:

    $ garchive --links links.csv --archive data

By default, `garchive` will bind to address `0.0.0.0` and port `8080`. Use the
`-addr` and `-port` commandline flags to change those settings.

Your archive will now be available on `localhost` at port `8080`.

## Importing from Pinboard

You can easily import bookmarks from Pinboard into your `links.csv` file using
`curl` and `jq`:

    $ curl -s "https://api.pinboard.in/v1/posts/all?auth_token=$PINBOARD_API_KEY&format=json" | jq -r '.[] | "\(.description)\t\(.href)"' > links.csv

## Using Docker

First, build the Docker image. You only need to do this once.

    $ git clone https://github.com/gpanders/garchive
    $ cd garchive
    $ docker build -t garchive .

Provide the path to your `links.csv` and archive directory as volumes to the
Docker container:

    $ docker run -v /path/to/links.csv:/app/links.csv -v /path/to/archive:/app/data -p 8080:8080 garchive

You can create a `docker-compose.yml` file to easily generate a `garchive`
container:

```yaml
version: "3"

services:
  garchive:
    container_name: garchive
    image: garchive
    build: ./
    ports:
      - 8080
    volumes:
      - /path/to/links.csv:/app/links.csv
      - /path/to/archive/data:/app/data
```

[link rot]: https://en.wikipedia.org/wiki/Link_rot
[go]: https://golang.org/doc/install
[fetch]: ./bin/fetch
[release]: https://github.com/gpanders/garchive/releases
