# Spolyr - a private Spotify lyrics index

Ever had the lyrics of a song stuck in your head for weeks but couldn't remember the title of the track? Did it nearly drive you crazy knowing the track is somewhere buried in your Spotify library? Fear not, Spolyr to the rescue!

Spolyr is a side project I've been working on that helps you index and retrieve your favorite songs on Spotify by querying a full-text index. 

## Features
- Sign in using your Spotify account and download all tracks in your library
- Import Spotify playlists
- Automatically fetch lyrics from different providers
- Find a specific song by querying a full-text search index

## Prerequisites
- go to https://developer.spotify.com/dashboard/applications and register a new app 
- go to http://genius.com/api-clients and get a new Genius API token
- MongoDB or preferably Docker

## How to get started?

1. Set secrets as environment variables inside your `docker-compose.yml`
2. run `docker-compose up`
3. Open [localhost:8080](http://localhost:8080)

## Configuration options

### Environment variables

`SPOTIFY_ID`: unique identifier of your Spotify application - **required**

`SPOTIFY_SECRET`: secret Spotify key - **required**

`GENIUS_API_TOKEN`: API token to use when communicating with genius.com - **required**

`SUPPORTED_LANGUAGES` Used for language-specific query preprocessing and language detection. The more languages you
enable, the more RAM is required. If you tend to listen only to say English and German songs, you can reduce the
resource usage by limiting the selection to the subset "english,german". (
default: [all languages currently supported by MongoDB](https://www.mongodb.com/docs/manual/reference/text-search-languages/#std-label-text-search-languages))

`SESSION_KEY`: key used for signing cookies. (default: a random key)

`HTTP_PORT`: Specifies the http port to bind Spolyr to (default: `8080`)

`PROTOCOL`: Http protocol (default: `http`)

`DOMAIN`: Domain name of this server. (default: `localhost`)

`HTTP_PUBLIC_PORT`: Specifies the public-facing http port. Set this to `443` or `80` if you are running Spolyr with a
reverse proxy (default: value of `HTTP_PORT`)

`DATABASE_HOST`: (default: `127.0.0.1`)

`DATABASE_USER` default: `root`)

`DATABASE_PASSWORD` (default: `example`)

### Configuration file

Alternatively, all configuration options can be set by using a `config.yaml`:

```yaml
debug: 1
genius_api_token: "genius_api_token"
http_public_port: 8081
protocol: "https"
session_key: "a-secret"
spotify_id: "spotify_oauth_id"
spotify_secret: "spotify_oauth_secret"
supported_languages: "german,english,french,russian"
```

## Screenshots

![home page](doc/preview-1.png "Import and query your Spotify library.")

![full-text search](doc/preview-2.png "Automatically import lyrics from different providers")

![track details page](doc/preview-3.png "Search for songs by parts of the lyrics, title, album name, and artists")

![import of lyrics](doc/preview-4.png "View end edit lyrics of tracks.")

## Development

1. Install `node` and `docker`
2. Install frontend dependencies `npm i`
3. Generate the api and clients stubs: `make openapi-spec`
4. Start database: `docker compose -f docker-compose.dev.yml up -d`
5. Set the following environment variables:
   - DATABASE_HOST=127.0.0.1
   - DATABASE_PASSWORD=example
   - DATABASE_USER=root
   - DOMAIN=localhost
   - HTTP_PUBLIC_PORT=8081
   - PROTOCOL=https
   - SESSION_KEY=dev
   - SPOTIFY_ID=YOUR_ID
   - SPOTIFY_SECRET=YOUR_SECRET
   - GENIUS_API_TOKEN=YOUR_TOKEN
6. Start api: `go run main.go web`
7. Start webpack dev server: `npm run serve`
8. Open [localhost:8080](https://localhost:8080) in your preferred browser

### Tests and linting

```bash
# frontend
make lint-frontend
make test-frontend

# backend
make test

# e2e
make test-e2e
```