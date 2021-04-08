# Administration commands

Terms are added, edited, and removed with the bot's administration commands.

All of these commands are subcommands of `t;admin`; for a list, check `t;admin help`.

## Director commands

These commands can be used by anyone with a director role (`bot.support.staff_roles`), as well as bot admins and owners.

### `t;admin import`

Adds a term from a correctly formatted message. The bot will automatically react with a ✅ emoji if the import was successful.

**Usage:** `t;admin import [-c category] [-r] <message link|ID>`

**Available flags:**

- `-c`/`--category`: override the category, useful if it was misspelled or wasn't detected correctly.
- `-r`/`--raw-source`: this command prepends "Coined by" to the source by default. This flag disables that.

**Examples:**  
`t;admin import https://discord.com/channels/793828572472803358/793832038091325461/829092754444648449`  
`t;admin import -r https://discord.com/channels/793828572472803358/793832038091325461/825365217973239848`  
`t;admin import -c lgbtq+ https://discord.com/channels/793828572472803358/794302373602263060/821961310194630656`

### `t;admin editterm`

This command can be used to edit an existing term.

**Usage:** `t;admin editterm <part> <ID> <new|-clear>`
​
Available parts to edit are:
- `title`
- `desc` (description)
- `source` ("Coined by")
- `aliases` ("Synonyms")
- `tags`

For `aliases` and `tags`, you can use "-clear", without quotes, to clear them.

- `title`  
  The term's new title
- `desc`  
  The term's new description. Note that this should be wrapped in "quotes" to preserve newlines.
- `source`  
  The term's new source.
- `aliases`  
  The term's new synonyms. Synonyms should be space separated; if a synonym has a space in it, wrap it in "quotes".
- `tags`  
  The term's new tags, space separated, like `aliases`.

**Examples:**

```
t;admin editterm title 1 Plural

t;admin editterm desc 410 "Name or term for someone in a wavership."

t;admin editterm source 2 "Unknown; already in circulation."

t;admin editterm aliases 7 NV Non-Verbal "Non Verbal"

t;admin editterm tags 8 Plurality "Member Type"
```

### `t;admin setcw`

Sets a term's content warning.  
A term with a content warning is automatically spoiler-tagged in the bot.  
Use with `-clear` as an argument to remove the term's content warning.

**Usage:** `t;admin setcw <id> <content warning>`

**Examples:**  
`t;admin setcw 1 Religion`  
`t;admin setcw 1 -clear`

### `t;admin setnote`

Like `t;admin setcw`, but the note does *not* automatically spoiler-tag a term.

**Usage:** `t;admin setcw <id> <note>`

**Examples:**  
`t;admin setcw 1 Please be careful not to stigmatise persecutors.`  
`t;admin setcw 1 -clear`