# Application Details
**What does your application do? Please be as detailed as possible, and feel free to include links to image or video examples.**  
Termora is a glossary bot for terms coined by and used by the plural and LGBTQ+ communities. It allows users to search the database for terms and pronouns that might suit them.

# Data Collection

## What Discord data do you store?
The bot stores server IDs and channel IDs for a command blacklist. It can also store user IDs and channel IDs for error logging: when Sentry integration is disabled, it stores these locally along with the error text.

## For what purpose(s) do you store it?
The command blacklist is used to let server moderators control in which channels the bot responds to commands, without having to completely hide it from those channels. Errors are logged for debugging purposes.

## For how long do you store it?
Server and channel IDs for the command blacklist are stored as long as the bot is in a server. When Sentry integration is disabled, local error logs are stored indefinitely, but this isn't used in practice, as the production version of the bot always logs errors to Sentry.

## What is the process for users to request deletion of their data?
To remove server and channel IDs from the database, they can simply remove the bot from their server. To delete error logs, they should contact us (the developers), we will delete those on request.

# Infrastructure

## What systems and infrastructure do you use?
The bot and connected services are all run from a single Hetzner cloud server, along with a couple of other projects. They are separated from those projects by being run as a separate local user.

## How have you secured access to your systems and infrastructure?
The server is located remotely, and only accessed through SSH. Root login is blocked, as well as password login: all logins must be done with a public/private key pair. The bot's code is uploaded directly to the server via Gitea, as a mirror of the upstream repository hosted on GitHub.

## How can users contact you with security issues?
Users can contact us by joining the support server (linked in the bot's help command, and on the website) and opening a private ticket channel, or by DMing us at starshine system 🌠✨#0001, or by emailing an address linked on the bot's website.

## Does your application utilize other third-party auth services or connections? If so, which, and why?
Errors are logged to Sentry (https://sentry.io/) for debugging purposes. The database the bot uses also serves a website and a REST API, but there is no other connection between the bot and those two services.