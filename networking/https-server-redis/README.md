# HTTPS Server with Redis

Simple server that prints (newest first) a few recent access logs.

Backed by Redis. Runs inside docker containers. Deployed via docker-compose.

## Behavior to User

After writing each line of response, we attempt to flush it and then wait 300ms
so that the user can see the page load gradually.

## Source Code

Package split into appropriate sub-packages.

* `./internal` has app-specific packages
* `./internal/util` has app-specific utilities
  * `tls.Config` generator
  * error checker that quits app on error
* `./internal/platform` has database level tooling
  * Stack: an abstraction of the stack data structure (with only the methods we
    need implemented) and an implementation based on Redis.
* `./handlers` holds HTTP handlers required for app. Curently just a single
  handler that does everything.
* `main` is in the root package. This reads user flags, sets db connections,
  launches a server and runs handler. It also handles graceful shutdown of
  server. Much of the code here can be reused in future projects.

## DevOps with Docker

Application fully dockerized and deployed via docker-compose.

TLS credential sharing and redis persistence achieved via docker volumes.

Application uses compose file v2.2 that passes --init flag to Docker. Swarm
does not currently support the init flag and uses compose file v3 instead. See
[this github issue](https://github.com/docker/compose/issues/5049) for details.

## Makefile

Currently just a frontend for the docker-compose command.

Running `make` asks compose to rebuild image from Dockerfile and run app in
daemon mode.

Runing `make clean` asks compose to tear down containers and also delete built
local images (not pulled remote images) so as to not leave any unnecessary
artefacts behind.

## TODO Tasks

### Package Organization

We don't need internal/ style of packge organization. Use only a single level of
package nesting for a simple project like this one.

### Dynamic TLS Loading

Use tls.Config for dynamic loading of TLS certificates based on this
[example](https://diogomonica.com/2017/01/11/hitless-tls-certificate-rotation-in-go/)

### Test Cases

Add test-cases. Figure out the best way to unit and integration test every
single component of this app.

