# Playhead tracker

## Purpose

To track a where a user is in an episode in relationship to the Series total

## TODO 

[ ] Build out the listener for new/changed Series and Episodes    


`make build` to build the project

## Migrations 

`make build` then `./playhead migrate`

## Web Server

`./playhead serve`

## GDPR Listener

`./playhead gdpr_listener`

## All

```
mike@MacBook-Pro-6 ~/go/src/playhead (master) $ ./playhead 
Playhead Web Application

Usage:
  playhead [command]

Available Commands:
  gdpr_listener Start the GDPR queue listener
  help          Help about any command
  migrate       
  routes        Print the routes
  serve         serves the api
  version       Print the version number

Flags:
      --config string   config file (default is config.yaml)
  -h, --help            help for playhead
  -v, --verbose         make output more verbose

Use "playhead [command] --help" for more information about a command.

```


