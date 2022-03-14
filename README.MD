# File 2 URL Bot

A simple bot to convert files sent to it to URLs to make direct download links.

## Features

* Fast download speeds
* Option to allow only specific users to create download links
* Range support
* Database support

## Setup

At first download the program from [releases](https://github.com/HirbodBehnam/file2url-go/releases). Then before running
the application, set these values as environment variables:

```bash
export APP_ID=1234
export APP_HASH=abcd
export BOT_TOKEN=zyxw
export ALLOWED_USERS=1,2
```

`ALLOWED_USERS` is the list of user IDs which are allowed to use the bot. Do not specify this variable in order to make
the bot public.

### Running on your server

If you are planning to run this bot on your server you also have to specify these two environment variables:

```bash
export LISTEN=0.0.0.0:12345
export URL_PREFIX=my-ip-or-domain:12345
```

`LISTEN` is the listen address of the server. `URL_PREFIX` will be as the domain/ip of the links which are sent to
users.

Then simply run the executable file.

### Running with Heroku

To run this app on Heroku, you only need to specify one environment variable other than bot id and app hash/id.

```bash
export DYNO_NAME=my-dyno-name
```

### Databases

This bot will use a database to store the file information which is mapped from a UUID. By default, bot will use an in
memory cache to store these values. This each entry has a TTL of one day by default.

You can change the in memory cache to [bbolt](https://github.com/etcd-io/bbolt) in order to make the data persistent on
disk. To enable bbolt, set the environment variable of `BOLT_DB_PATH` to a file which is the database.