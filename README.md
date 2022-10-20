[![Go Tests](https://github.com/bitstonks/go-adt/actions/workflows/go.yml/badge.svg)](https://github.com/bitstonks/go-adt/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/bitstonks/go-adt)](https://goreportcard.com/report/github.com/bitstonks/go-adt)
[![Go Reference](https://pkg.go.dev/badge/github.com/bitstonks/go-adt.svg)](https://pkg.go.dev/github.com/bitstonks/go-adt)
[![Go Version](https://img.shields.io/github/go-mod/go-version/bitstonks/go-adt.svg)](https://github.com/bitstonks/go-adt/blob/main/go.mod)
[![Module Version](https://img.shields.io/github/v/tag/bitstonks/go-adt.svg)](https://github.com/bitstonks/go-adt/tags)

# go-adt

Go implementations of different abstract data types using generics.
Requires Go 1.18+.

 * `./set`: [generic set](https://pkg.go.dev/github.com/bitstonks/go-adt/set)
 * `./broadcast`: [one to many broadcast service](https://pkg.go.dev/github.com/bitstonks/go-adt/broadcast)
     * [NoSyncBroadcaster](https://pkg.go.dev/github.com/bitstonks/go-adt/broadcast#NoSyncBroadcaster) - subscribe, unsibscribe and send actions have to be synchronised externally
     * [SyncBroadcaster](https://pkg.go.dev/github.com/bitstonks/go-adt/broadcast#SyncBroadcaster) - actions are synchronised using an internal mutex
     * [ChanBroadcaster](https://pkg.go.dev/github.com/bitstonks/go-adt/broadcast#ChanBroadcaster) - actions are synchronised using channels and processed in an eventloop