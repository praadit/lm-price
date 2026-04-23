# LM Price API

A small **Go** service that fetches public gold price tables from upstream HTML pages, parses them into structured JSON, and exposes them over HTTP.

Upstreams (overridable via environment variables):

- LM: `https://emasantam.id/content/lm.txt`
- Antaremas: `https://antaremas.com/harga-emas/`

This project is **not** affiliated with those sites; it only reads the same public URLs you could open in a browser.

## Live deployment

A sample build is hosted on Vercel:

**Base URL:** [https://lm-price.vercel.app/](https://lm-price.vercel.app/)

| Endpoint | URL |
|----------|-----|
| Health | [https://lm-price.vercel.app/health](https://lm-price.vercel.app/health) |
| Prices (JSON) | [https://lm-price.vercel.app/v1/prices](https://lm-price.vercel.app/v1/prices) |
| Antaremas buy prices (JSON) | [https://lm-price.vercel.app/v1/antaremas/prices](https://lm-price.vercel.app/v1/antaremas/prices) |

Example (filtered): [https://lm-price.vercel.app/v1/prices?area=Area%20Jawa-Bali&location=Bandung](https://lm-price.vercel.app/v1/prices?area=Area%20Jawa-Bali&location=Bandung)

## What it does

1. **Downloads** the LM document over HTTPS.
2. **Parses** embedded tables (per region / butik) into a list of locations, each with gram-based **price**, **stock**, and **sold out** flags.
3. **Serves** JSON via a Gin API, with optional **area** and **location** query filters. Unknown filter values return **400** with lists of valid **areas** and **locations** from the latest scrape.
4. **Scrapes** Antaremas' “Ukuran / Harga Beli” table into a simple list of size + buy price (no stock, no area, no location).

## API

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Liveness check: `{"status":"ok"}`. |
| `GET` | `/v1/prices` | Parsed prices as JSON (see below). |
| `GET` | `/v1/antaremas/prices` | Antaremas “Harga Beli” table as JSON (see below). |

### `GET /v1/prices` query parameters

| Query | Description |
|--------|-------------|
| `raw=1` | Return the raw upstream document as `text/plain` (no parsing). |
| `area` | Filter by region label (e.g. `Area Jawa-Bali`). Case-insensitive; whitespace is normalized. |
| `location` | Filter by butik / column name (e.g. `Bandung`). Same matching rules as `area`. |

If `area` or `location` does not match any value in the current scrape, the response is **400** with a JSON body that includes `code`, `message`, and `available_areas` / `available_locations` as appropriate.

### JSON shape (default response)

The LM endpoint returns an envelope with the upstream `last_update` (RFC3339) and `data` rows.

```json
{
  "last_update": "2026-04-23T11:30:02+07:00",
  "data": [
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
}
```

`price` is in **IDR** (integer). When `sold_out` is `true`, `stock` is `0`.

### JSON shape (`GET /v1/antaremas/prices`)

Antaremas endpoint returns the closest “Terakhir Diperbarui …” timestamp (RFC3339) and the “Harga Beli” table (size + buy price):

```json
{
  "last_update": "2026-04-23T09:30:00+07:00",
  "data": [
    { "size": "0.5 gram", "buy_price": 1655000 },
    { "size": "1 gram", "buy_price": 3074000 }
  ]
}
```

## Configuration

| Variable | Default | Purpose |
|----------|---------|---------|
| `PORT` | `8080` | Listen port (value is used as `:{PORT}`). |
| `GIN_MODE` | unset (debug) | Set to `release` for production-style Gin logging. |
| `LM_SOURCE_URL` | `https://emasantam.id/content/lm.txt` | URL of the LM HTML document to fetch. |
| `ANTAREMAS_SOURCE_URL` | `https://antaremas.com/harga-emas/` | URL of the Antaremas page to fetch. |

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
- Antaremas: `http://127.0.0.1:8080/v1/antaremas/prices`
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
- `internal/domain/antaremas` — entities, parsing, and the `RawSource` port.
- `internal/usecase` — orchestrates fetch → parse → filter.
- `internal/repository/lmremote` — HTTP implementation of `lm.RawSource`.
- `internal/repository/antaremasremote` — HTTP implementation of `antaremas.RawSource`.
- `internal/delivery/http` — Gin router and handlers.

`migrations/` and `internal/pkg/` are reserved for future database migrations and shared utilities.

## License

See the repository owner for license terms (not specified in this README).
