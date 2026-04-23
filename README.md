# LM Price API

A small **Go** service that **fetches public ANTAM logam mulia (LM) price tables** from a fixed upstream HTML document, **parses** them into structured JSON, and exposes them over **HTTP** with optional **filters** and **validation**.

The upstream content is the HTML snippet published at `https://emasantam.id/content/lm.txt` (overridable via environment variable). This project is **not** affiliated with that site; it only reads the same public URL you could open in a browser.

## Live deployment

A sample build is hosted on Vercel:

**Base URL:** [https://lm-price.vercel.app/](https://lm-price.vercel.app/)

| Endpoint | URL |
|----------|-----|
| Health | [https://lm-price.vercel.app/health](https://lm-price.vercel.app/health) |
| Prices (JSON) | [https://lm-price.vercel.app/v1/prices](https://lm-price.vercel.app/v1/prices) |

Example (filtered): [https://lm-price.vercel.app/v1/prices?area=Area%20Jawa-Bali&location=Bandung](https://lm-price.vercel.app/v1/prices?area=Area%20Jawa-Bali&location=Bandung)

## What it does

1. **Downloads** the LM document over HTTPS.
2. **Parses** embedded tables (per region / butik) into a list of locations, each with gram-based **price**, **stock**, and **sold out** flags.
3. **Serves** JSON via a Gin API, with optional **area** and **location** query filters. Unknown filter values return **400** with lists of valid **areas** and **locations** from the latest scrape.

## API

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Liveness check: `{"status":"ok"}`. |
| `GET` | `/v1/prices` | Parsed prices as JSON (see below). |

### `GET /v1/prices` query parameters

| Query | Description |
|--------|-------------|
| `raw=1` | Return the raw upstream document as `text/plain` (no parsing). |
| `area` | Filter by region label (e.g. `Area Jawa-Bali`). Case-insensitive; whitespace is normalized. |
| `location` | Filter by butik / column name (e.g. `Bandung`). Same matching rules as `area`. |

If `area` or `location` does not match any value in the current scrape, the response is **400** with a JSON body that includes `code`, `message`, and `available_areas` / `available_locations` as appropriate.

### JSON shape (default response)

Each element is one **location** (butik) with nested **prices**:

```json
[
  {
    "location": "Bandung",
    "product": "Harga Emas ANTAM Certicard Fine Gold Bar 999.9",
    "area": "Area Jawa-Bali",
    "prices": [
      {
        "gramasi": 0.5,
        "price": 1452500,
        "stock": 44,
        "sold_out": false
      }
    ]
  }
]
```

`price` is in **IDR** (integer). When `sold_out` is `true`, `stock` is `0`.

## Configuration

| Variable | Default | Purpose |
|----------|---------|---------|
| `PORT` | `8080` | Listen port (value is used as `:{PORT}`). |
| `GIN_MODE` | unset (debug) | Set to `release` for production-style Gin logging. |
| `LM_SOURCE_URL` | `https://emasantam.id/content/lm.txt` | URL of the LM HTML document to fetch. |

HTTP timeouts are defined in code (`internal/config`): **15s** for the upstream client, **20s** per `/v1/prices` request context.

## Run locally

Requires **Go 1.25+** (see `go.mod`).

```bash
go run ./cmd/api
```

Or:

```bash
make run-api
```

Then:

- Health: `http://127.0.0.1:8080/health`
- Prices: `http://127.0.0.1:8080/v1/prices`
- Example filter: `http://127.0.0.1:8080/v1/prices?area=Area%20Jawa-Bali&location=Bandung`

## Build

```bash
go build -o bin/api ./cmd/api
./bin/api
```

## Project layout

The code follows a **layered** layout (domain, use case, HTTP delivery, remote repository):

- `cmd/api` — process entrypoint.
- `internal/config` — environment-backed settings.
- `internal/domain/lm` — entities, parsing, filtering, and the `RawSource` port.
- `internal/usecase` — orchestrates fetch → parse → filter.
- `internal/repository/lmremote` — HTTP implementation of `lm.RawSource`.
- `internal/delivery/http` — Gin router and handlers.

`migrations/` and `internal/pkg/` are reserved for future database migrations and shared utilities.

## License

See the repository owner for license terms (not specified in this README).
