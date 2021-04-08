# Installing

Termora is open source, and you can run your own instance of the bot, website, and API.  
The source code is available [here](https://github.com/termora/berry).

## Requirements

All components require the following:

- A working Go 1.16 installation
- A working PostgreSQL installation (only 12.5 has been tested)

Additionally, the website and API will most likely need to be used behind a reverse proxy.

## Bot

The bot's code resides in the `cmd/bot` directory.  
To build it, run `go build`.

The bot is configured using `config.json`, a sample of which is available in the bot's directory.
It isn't complete, and the [BotConfig type](https://github.com/termora/berry/blob/main/structs/botConfig.go)
in `structs/botConfig.go` is authorative.

The bot can be sharded using `termora.service`, by enabling `termora@0`, `termora@1` etc.

The following configuration keys are *required*:

```
- auth:
  - token (string): Discord bot token
  - database_url (string): dsn for the Postgres database
- bot:
  - prefixes: ([]string): default prefixes used
  - bot_owners: ([]int): bot owner IDs, these users can use all commands including admin commands
```

## Site

The website's code resides in the `cmd/site` directory, and can also be built using `go build`.
It uses `config.yaml` for its configuration, a sample of which is available as `config.sample.yaml`. All keys are required.

## API

The api's code resides in the `cmd/api` directory, and can also be built using `go build`.
It uses `config.yaml` for its configuration, a sample of which is available as `config.sample.yaml`. All keys are required.