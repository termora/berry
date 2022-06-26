# Berry

Berry is a searchable glossary bot for Discord.

## Requirements

The easiest way to run Berry is with Docker. To run it on bare metal, you need the following:

- PostgreSQL (only 12.5 tested)
- Typesense 0.23 or later
- Go 1.16 or later
- For the site and API: a reverse proxy (such as Caddy or nginx)
- Optionally, but strongly recommended: Redis, to not spam users with deprecation warnings for text commands

## Configuration

All services are configured with a `config.toml` file in the root of this repository.
An example can be seen in `config.toml.example`.

The following keys are required or strongly recommended to be set (using `.` to indicate nesting):
- `bot.token`: the Discord bot token
- `bot.prefixes`: the prefixes the bot will respond to. Will also respond to mentions
- `bot.bot_owners`: users that have absolute control over the bot
- `bot.admins`: roles that can do most administrative tasks, such as adding, editing, or deleting terms
- `bot.directors`: roles that can add and edit terms and pronouns, but cannot delete terms

If you're not using Docker, the following keys are also required or recommended:
- `core.database_url`: the DSN for the PostgreSQL database (required)
- `core.typesense_url`: the URL for the Typesense server (required)
- `core.typesense_key`: the API key for the Typesense server (required)
- `core.redis`: the URL for the Redis instance

## Running

The easiest way to get Berry running is with Docker.

- Clone this repository: `git clone https://github.com/termora/berry`
- Create a `config.toml` file in the same directory as `docker-compose.yml`, containing at least a `bot.token` field
- Build the bot: `docker-compose build`
- Run the bot: `docker-compose up` (or `docker-compose up -d` to run in the background)

The site will listen on `localhost:2839`, and the API will listen on `localhost:2838`.

If you want to run it on the bare metal, it's a lot more involved and you're mostly on your own.  
If you get stuck, feel free to ask for help on the [support server](https://termora.org/server).

## License

Copyright (C) 2022, Starshine System

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
