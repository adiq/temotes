#  ![](https://emotes.adamcy.pl/v1/channel/adiq/emotes/7tv/proxy?emote=Harambe&size=1x) tEmotes API
 
Easy to use API for Twitch emotes

[![Documentation](https://img.shields.io/badge/docs-see_how_to_use-brightgreen?style=for-the-badge&logo=readthedocs)](https://adiq.stoplight.io/docs/temotes/YXBpOjMyNjU2ODIx-t-emotes-api)

We support:
* Twitch
* 7TV
* BetterTTV
* FrankerFaceZ

## Setup

> Note: _Keep in mind that you don't need to install anything if you just want to consume the API._
> 
> You can use the public API that we expose on [emotes.adamcy.pl](https://adiq.stoplight.io/docs/temotes/YXBpOjMyNjU2ODIx-t-emotes-api)

### Requirements

* Golang
* Redis
* Twitch API Access (Client ID and Client Secret)

### Configure

Configuration is as easy as defining the environment variables from the `.env` file.

### Run & Build

Running and building the application is as simple as in any other Go project.

### Health Endpoints

The API provides two health check endpoints for monitoring:

- **Liveness Probe** – checks if the server is running (`/livez`)
- **Readiness Probe** – checks if the application and its dependencies are ready to handle requests (`/readyz`)

# License

This project is licensed under the terms of the [AGPL-3.0 license](agpl-3.0.md).
