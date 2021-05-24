
# upl

A dropbox backed by any S3 bucket

## Features

- User interface to create new dropboxes
- Store bucket credentials and settings in a Redis KV store
- Drop stored credentials when the dropbox link expires
- Upload files using S3 multipart uploads, powered by [Uppy](https://uppy.io)
- Remembers previously created dropboxes and uploaded files
- Single fat binary

<img alt="Screenshot 1" width=400 src="https://user-images.githubusercontent.com/1705906/119283410-842f0d80-bc6f-11eb-9126-fdfce8d44dd5.png"><img alt="Screenshot 2" width=400 src="https://user-images.githubusercontent.com/1705906/119283425-914bfc80-bc6f-11eb-97b0-c7b74fecb192.png">

## Building

You'll need:
- Node.js
- Go
- `make`

```
make TAGS=production
```

Alternatively, `docker build .` this project.

## Deploying

To deploy as a single binary, build it from source and run:

```
export REDIS_CONNECTION=simple:redis-hostname:6379
export LISTEN=:8080
./upl
```

For example Kubernetes manifests or Docker Compose files, see the [deployments](./deployments) folder.

<!-- vim: set conceallevel=2 et ts=2 sw=2: -->
