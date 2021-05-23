
# upl

A dropbox backed by any S3 bucket

## Features

- User interface to create new dropboxes
- Store bucket credentials and settings in a Redis KV store
- Drop stored credentials when the dropbox link expires
- Upload files using S3 multipart uploads, powered by [Uppy](https://uppy.io)
- Remembers previously created dropboxes and uploaded files
- Single fat binary

## Deploying

To deploy as a single binary, build it from source and run:

```
export REDIS_CONNECTION=simple:redis-hostname:6379
export LISTEN=:8080
./upl
```


<!-- vim: set conceallevel=2 et ts=2 sw=2: -->
