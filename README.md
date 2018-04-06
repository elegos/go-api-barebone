# API barebone project written in golang

This project aims to help developers bootstrap an API application with integrated
OAuth2/JWT [grant type password](https://aaronparecki.com/oauth-2-simplified/#password)
authentication.

## Installation
It can be cloned and then strip away the `.git` folder, like this:

```bash
# Setup GOPATH, if not specified (to the default value)
if [ "$GOPATH" == "" ]; then export GOPATH="$HOME/go"; fi

git clone git@github.com:elegos/go-api-barebone.git "$GOPATH/src/my-awesome-app"
rm -rf "$GOPATH/src/my-awesome-app/.git"
```

This will install your new application in the default GOPATH. If you want though
you can also develop your application in wherever folder you want, using the
`unpack.sh` helper:

```bash
cd wherever/you/want
git clone git@github.com:elegos/go-api-barebone.git .
rm -rf .git
./unpack.sh
```

This will move the sources in a local go's environment, rename the packages where
needed, download the required dependencies and also creating a handy `.env` file
to use with your editor (if you're using Atom, I suggest installing `go-plus` and
my dotenv loader plugin `load-dotenv-variables`)

## Project structure
```
project-dir
|- handlers     // handlers ("controllers") folder
|  |- routes.go // where routes are defined
|- src          // services (called by the handlers) folder
|  |- bbAuth    // oauth2 library (barebone auth)
|- types        // where go's types (interfaces, structs, types) are stored
|- main.go      // the main file
|- unpack.sh    // the unpack utility
```

## Nomenclature
`/src/bb[whatever]` is reserved. This means that you can use the `/src` directory
as far as you don't prefix any folder with `bb`. The `bb` folders are used
for the barebone libraries, so if you use the prefix, you may encounter
future breaking changes.
