# zserve

Read-only API for Zendesk ticket data mirrored by [zmirror](https://github.com/perbu/zmirror).

## Requirements

- Go 1.25+
- PostgreSQL with a populated `zendesk` schema (via zmirror)

## Setup

Copy `.env.example` to `.env` and set `DATABASE_URL`:

```
DATABASE_URL=postgres://user@localhost:5432/zmirror?sslmode=disable
```

Edit `config.yaml` to set the listen address and optional API key:

```yaml
server:
  address: ":8080"
  api_key: ""
```

## Build and run

```
go build -o zserve ./cmd/zserve
./zserve -config config.yaml
```

## API

See `openapi.yaml` for the full spec.

- `GET /v1/tickets` — list tickets (filterable by status, type, priority, assignee, requester, organization, tag, search, date range; paginated)
- `GET /v1/tickets/{id}` — single ticket with tags, collaborators, followers, email CCs, and followup IDs
- `GET /v1/tickets/{id}/tags`
- `GET /v1/tickets/{id}/collaborators`
- `GET /v1/tickets/{id}/followers`
- `GET /v1/tickets/{id}/email-ccs`
- `GET /v1/tickets/{id}/followups`

If `api_key` is set in config, all requests must include `X-API-Key` header.

## Code generation

```
go tool oapi-codegen --config oapi-codegen.yaml openapi.yaml
go tool sqlc generate
```

## License

BSD 2-Clause. See [LICENSE](LICENSE).
