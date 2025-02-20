# artifactsmmo-3-season

This is my repository for the `3rd` season of [artifactsmmo](https://artifactsmmo.com). It was my second season that I participated in.

## Main goal

The main goal was to get all the achievements before the end of the season.

Of the additional goals were:
- use of observability, such as logs and metrics
- automated use of buy/sell mechanics
- more beautiful code structure
- write a simulator for fights

## Structure

- `gen/` - Generated code for [openapi scheme](https://docs.artifactsmmo.com/api_guide/openapi_spec)
- `internal/`
    - `internal/api` - helpers for authrization on server. I don't know why ogen not provide it
    - `internal/generic` - wrappers for actions
    - `internal/game` - wrappers (and caches) for game entities
    - `internal/live` - my characters
    - `internal/macro` - combinations of actions
    - `internal/simulator` - fight simulator
    - `internal/strategy` - ready sets of actions (one more level of abstraction). Strategies used directly by characters

## Strategis

They are the same with [2-nd season](https://github.com/Sinketsu/artifactsmmo-2-season#strategis)

## Observations about this season
- All the achievements were achieved only by the end of the season. The first leaders got all the achievements a little earlier, but still at the end
- Using Yandex.Cloud for observability is certainly possible, but it is very inconvenient and sparse compared to the `Grafana` installation
- The GE has been pretty useless to me this season. All the players who started from the very beginning of the game (including of me) achieved everything approximately evenly, which means none of the higher-level characters could help with rare resources. That's why I didn't invest in automating this process
- In the end, quite a lot of gold accumulated, which had nowhere to spend (

## Ideas (and plans) for the next season

- Try another openapi generator, than `ogen`. It's a good place to start, but for some reason it doesn't provide convenient helpers for work and you have to write a lot of code yourself.
- Collect various statistics. Initially, I thought to use logs for this, but in Yandex.Cloud you can not make SQL-like queries on logs to collect statistics.
- Use some more thoughtful code structure. Now many components depend on each other directly, which is bad.


