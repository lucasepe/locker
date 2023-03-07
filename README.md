# Locker

Store secrets on your local file system.

A Locker is a store on your file system (built on top of the amazing [bbolt](https://github.com/etcd-io/bbolt)).

- create as many lockers as you need

## Secret

Secrets are credentials, tokens, secure notes, credit cards, and any info you want.

- a secret has a label and a content
- create an unlimited number of secrets
- organize secrets into boxes
- secrets are encrypted and decrypted automatically
  - using the environment variable `LOCKER_SECRET` with your master secret phrase
  - encryption will be done using [AES-256-CFB](https://it.wikipedia.org/wiki/Advanced_Encryption_Standard)

## Box

Boxes are used to group and organize your secrets.

```txt
┬  ┌─┐┌─┐┬┌─┌─┐┬─┐
│  │ ││  ├┴┐├┤ ├┬┘
┴─┘└─┘└─┘┴ ┴└─┘┴└─
Store secrets on your local file system.

Usage:
   locker <command>

Commands:
   delete   Delete one secret from a box or a whole box.
   get      Get one or all secrets from a box.
   help     Show a list of all commands or describe a specific command.
   import   Import secrets.
   info     Print build information and list all existing lockers.
   list     List all boxes or all secrets in a box.
   put      Put a secret into a box.
```

# How To Install

## MacOs

```sh
brew tap lucasepe/locker
brew install locker
```

or if you have already installed memo using brew, you can upgrade it by running:

```sh
brew upgrade locker
```

## From [binary releases](https://github.com/lucasepe/locker/releases) (macOS, Windows, Linux)

memo currently provides pre-built binaries for the following:

- macOS (Darwin)
- Windows
- Linux

1. Download the appropriate version for your platform from [locker releases](https://github.com/lucasepe/locker/releases).

2. Once downloaded unpack the archive (zip for Windows; tarball for Linux and macOS) to extract the executable binary. 

3. If you want to use from any location you must put the binary executable to your `Path` or add the directory where is it to the environment variables.

## Using [`Go`](https://go.dev/dl/) toolchain

```sh
go install github.com/lucasepe/locker@latest
```
