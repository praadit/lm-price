# LM Price API

A small **Go** service that fetches public gold price tables from upstream HTML pages, parses them into structured JSON, and exposes them over HTTP.

Upstreams (overridable via environment variables):

- LM: `https://emasantam.id/content/lm.txt`
- Antaremas: `https://antaremas.com/harga-emas/`
- Galeri24: `https://galeri24.co.id/harga-emas`

This project is **not** affiliated with those sites; it only reads the same public URLs you could open in a browser.

## Live deployment

A sample build is hosted on Vercel:

**Base URL:** [https://lm-price.vercel.app/](https://lm-price.vercel.app/)

| Endpoint | URL |
|----------|-----|
| Health | [https://lm-price.vercel.app/health](https://lm-price.vercel.app/health) |
| LM prices (JSON) | [https://lm-price.vercel.app/v1/prices/antam](https://lm-price.vercel.app/v1/prices/antam) |
| Antaremas “Harga Beli” (JSON) | [https://lm-price.vercel.app/v1/prices/hf](https://lm-price.vercel.app/v1/prices/hf) |
| Galeri24 “Harga ANTAM” (JSON) | [https://lm-price.vercel.app/v1/prices/galeri24](https://lm-price.vercel.app/v1/prices/galeri24) |

Example (filtered): [https://lm-price.vercel.app/v1/prices/antam?area=Area%20Jawa-Bali&location=Bandung](https://lm-price.vercel.app/v1/prices/antam?area=Area%20Jawa-Bali&location=Bandung)

## What it does

1. **Downloads** the LM document over HTTPS.
2. **Parses** embedded tables (per region / butik) into a list of locations, each with gram-based **price**, **stock**, and **sold out** flags.
3. **Serves** JSON via a Gin API, with optional **area** and **location** query filters. Unknown filter values return **400** with lists of valid **areas** and **locations** from the latest scrape.
4. **Scrapes** Antaremas' “Ukuran / Harga Beli” table into a simple list of size + buy price.
5. **Scrapes** Galeri24's “Harga ANTAM” table into weight + sell/buyback prices.
6. Caches upstream scrapes (configurable TTL) and enforces per-IP rate limiting with a higher quota for Basic-Auth requests.

## API

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Liveness check: `{"status":"ok"}`. |
| `GET` | `/v1/prices/antam` | Parsed LM prices as JSON (see below). |
| `GET` | `/v1/prices/hf` | Antaremas “Harga Beli” table as JSON (see below). |
| `GET` | `/v1/prices/galeri24` | Galeri24 “Harga ANTAM” table as JSON (see below). |

### `GET /v1/prices/antam` query parameters

| Query | Description |
|--------|-------------|
| `raw=1` | Return the raw upstream document as `text/plain` (no parsing). |
| `area` | Filter by region label (e.g. `Area Jawa-Bali`). Case-insensitive; whitespace is normalized. |
| `location` | Filter by butik / column name (e.g. `Bandung`). Same matching rules as `area`. |

If `area` or `location` does not match any value in the current scrape, the response is **400** with a JSON body that includes `code`, `message`, and `available_areas` / `available_locations` as appropriate.

### JSON shape (default response)

All endpoints return the same envelope:

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
          "buy_price": 1452500,
          "sell_price": 0,
          "stock": 44,
          "sold_out": false
        }
      ]
    }
  ]
}
```

Defaults:

- Missing `location` / `area` are set to `"Indonesia"`.
- Missing `stock` / `sell_price` are set to `0`.

`buy_price` and `sell_price` are in **IDR** (integer). When `sold_out` is `true`, `stock` is `0`.

Notes:

- `GET /v1/prices/antam`: `prices[].buy_price` comes from the upstream LM `price` column; `sell_price` is `0`. Stock/sold_out are populated.
- `GET /v1/prices/hf`: `prices[].buy_price` is the “Harga Beli” value; `stock` and `sell_price` are `0`.
- `GET /v1/prices/galeri24`: `prices[].buy_price` is the “Harga Jual” value and `sell_price` is the “Harga Buyback” value (stock is `0`).

## Configuration

| Variable | Default | Purpose |
|----------|---------|---------|
| `PORT` | `8080` | Listen port (value is used as `:{PORT}`). |
| `GIN_MODE` | unset (debug) | Set to `release` for production-style Gin logging. |
| `LM_SOURCE_URL` | `https://emasantam.id/content/lm.txt` | URL of the LM HTML document to fetch. |
| `ANTAREMAS_SOURCE_URL` | `https://antaremas.com/harga-emas/` | URL of the Antaremas page to fetch. |
| `GALERI24_SOURCE_URL` | `https://galeri24.co.id/harga-emas` | URL of the Galeri24 page to fetch. |
| `CACHE_TTL` | `60s` | Cache TTL for upstream scrapes (e.g. `30s`, `5m`). Set to `0s` to disable caching. |
| `BASIC_AUTH_USER` | (unset) | If set (with `BASIC_AUTH_PASS`), requests using HTTP Basic Auth get a higher rate limit. |
| `BASIC_AUTH_PASS` | (unset) | Basic Auth password. |
| `RATE_LIMIT_UNAUTHORIZED_PER_MINUTE` | `1` | Per-IP rate limit for unauthenticated requests (requests/min). |
| `RATE_LIMIT_AUTHORIZED_PER_MINUTE` | `100` | Per-IP rate limit for authenticated requests (requests/min). |

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
- LM prices: `http://127.0.0.1:8080/v1/prices/antam`
- Antaremas: `http://127.0.0.1:8080/v1/prices/hf`
- Galeri24: `http://127.0.0.1:8080/v1/prices/galeri24`
- Example filter: `http://127.0.0.1:8080/v1/prices/antam?area=Area%20Jawa-Bali&location=Bandung`

For a quick start, copy `.env.example` to `.env`.

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
- `internal/domain/galeri24` — entities, parsing, and the `RawSource` port.
- `internal/usecase` — orchestrates fetch → parse → filter.
- `internal/repository/lmremote` — HTTP implementation of `lm.RawSource`.
- `internal/repository/antaremasremote` — HTTP implementation of `antaremas.RawSource`.
- `internal/repository/galeri24remote` — HTTP implementation of `galeri24.RawSource`.
- `internal/delivery/http` — Gin router and handlers.

`migrations/` and `internal/pkg/` are reserved for future database migrations and shared utilities.

## License

See the repository owner for license terms (not specified in this README).
