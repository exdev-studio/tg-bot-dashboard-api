# Telegram Bot Dashboard API

An API that will interact with Telegram/Telegram Bot API

## Requirements

Currently, the api supports only Linux like systems.  
All the checks have been made using Ubuntu 20.04 via WSL 2; so in case of any troubles please add a GitHub issues to the
repo

In order to build and run the api you need to meet the following requirements:

* Go >= 1.15

## Build

To build the api please run the following:

```shell
make build
```

Results files are about to be in `bin` folder

## Run

To run the api you need to build it before. Please see [Build](#build) section.  
Once you have the api built please run the following:

```shell
./bin/apiserver
```
