This is a small helper bot that handles tasks in the bot support server
that Termora does not have the privileged intents to do itself.

Right now this is just requesting guild members for the `t;credits` command,
but once message content becomes privileged, it will probably expand to
those tasks too.

## Example .env file

```
TOKEN=yourDiscordTokenHere
GUILD_ID=755426161521328189
RPC=localhost:58952
```