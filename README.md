# Cloud Foundry CLI PSQL Plugin
[![Build Status](https://travis-ci.org/jaecktec/cf-psql-plugin.svg?branch=master)](https://travis-ci.org/jaecktec/cf-psql-plugin)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/jaecktec/cf-psql-plugin/blob/master/LICENSE)

cf-psql-plugin makes it easy to connect the `psql` command line client to any PSQL-compatible database used by
Cloud Foundry apps. Use it to

* inspect databases for debugging purposes
* manually adjust schema or contents in development environments

## Contents

* [Usage](#usage)
* [Removing service keys](#removing-service-keys)
* [Installing and uninstalling](#installing-and-uninstalling)
* [Building](#building)
* [Details](#details)

## Usage

```bash
$ cf psql -h
NAME:
   psql - Connect to a PSQL database service

USAGE:
   Open a psql client to a database:
   cf psql <service-name> [psql args...]
```

### Connecting to a database

Passing the name of a database service will open a PSQL client:

```bash
$ cf psql my-db
psql (10.3, server 9.6.3)
SSL connection (protocol: TLSv1.2, cipher: DHE-RSA-AES256-GCM-SHA384, bits: 256, compression: off)
Type "help" for help.

db_aaaaaaaa_bbbb_cccc_dddd_eeeeeeeeeeeee=>
```

## Removing service keys

The plugin creates a service key called 'cf-psql' for each service instance a user connects to. The keys are reused
when available and never deleted. Keys need to be removed manually before their service instances can be removed:

```bash
$ cf delete-service -f somedb
Deleting service somedb in org DevProjects / space constantin as constantin@dilab.com...
FAILED
Cannot delete service instance. Service keys, bindings, and shares must first be deleted.
```
Deleting the service failed. The CLI hints at service keys and app bindings that might still exist.
```bash
$ cf service-keys somedb
Getting keys for service instance somedb as constantin@dilab.com...

name
cf-psql
```
A key called 'cf-psql' is found for the service instance 'somedb', because we have used the plugin with 'somedb'
earlier. After removing the key, the service instance can be deleted:

```bash
$ cf delete-service-key -f somedb cf-psql
Deleting key cf-psql for service instance somedb as constantin@dilab.com...
OK

$ cf delete-service -f somedb
Deleting service somedb in org DevProjects / space constantin as constantin@dilab.com...
OK
```

This behavior might change in the future as it's not optimal to leave a key around.

## Installing and uninstalling

### Pre steps
You need to have the postgres client in your path. On Mac simply run `brew install postgresql`.

### Installing the plugin

The plugin is currently not available on the cloudfoundry marketplace. I am still waiting for some feedback. 

You can download a binary release or build yourself by running `go build`. Then, install the plugin with

```bash
$ cf install-plugin /path/to/cf-psql-plugin
```

The plugin can be uninstalled with:

```bash
$ cf uninstall-plugin psql
```

## Building

```bash
# download dependencies
go get -v ./...
go get github.com/onsi/ginkgo
go get github.com/onsi/gomega
go install github.com/onsi/ginkgo/ginkgo

# run tests and build
ginkgo -r
go build
```

## Details

### Obtaining credentials

cf-psql-plugin creates a service key called 'cf-psql' to obtain credentials. It no longer retrieves credentials from
application environment variables, because with the introduction of [CredHub](https://github.com/cloudfoundry-incubator/credhub/blob/master/docs/secure-service-credentials.md),
service brokers can decide to return a CredHub reference instead.

The service key is currently not deleted after closing the connection. It can be deleted by running:

```
cf delete-service-key service-instance-name cf-psql
```

A started application instance is still required in the current space for setting up an SSH tunnel. If you don't
have an app running, try the following to start an nginx app:

```bash
TEMP_DIR=`mktemp -d`
pushd $TEMP_DIR
touch Staticfile
cf push static-app -m 128M --no-route
popd
rm -r $TEMP_DIR
```
