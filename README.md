# Aman / CAS

A small central authentication service (CAS) in Go. It authenticates users
with one-time codes and never manages users: an identity is just the mobile
number or email address that received the code.

## How the flow works

```
app ──► /auth/check?app_id=...                 (browser)
          ├─ valid CAS cookie ──► 302 callback?ott=...
          └─ no/expired cookie ──► /auth/login

/auth/login ──► POST /auth/send-otp ──► /auth/verify-otp
/auth/verify-otp ──► POST /auth/verify-otp
          ├─ creates session + OTT
          └─ 302 callback?ott=...

app backend ──► POST /v1/auth/verify-ott {ott}
          └─ credentials (session token, identifier) + CAS/app cookies
```

- **OTT** is a short-lived, single-use ticket. Its lifetime is configured per
  app in the admin panel.
- The CAS session cookie (and, when the app shares a parent domain, the same
  cookie for the app domain) is set after the OTP is verified.

## Authentication modes

- **Managed** (default): Aman serves `/auth/check` and its branded login pages.
  The optional hosted-page setting redirects `/auth/login` to the app's own
  page while retaining the managed browser flow.
- **Custom** (`custom_auth_enabled=true`): Aman returns HTTP 404 for
  `/auth/**` requests for that app. The app owns its UI and calls the `/v1/auth`
  API from its backend. Verify the returned OTT server-side and create the
  app's own session cookie; do not expose `session_token` to browser code.

The API server does not configure CORS. Cross-origin browser code should call
the app backend rather than Aman directly.

## Routes

Browser / server-rendered:

| Route | Purpose |
| --- | --- |
| `GET /auth/check` | Reuse an existing CAS session or send the user to login |
| `GET /auth/login` | Branded login page |
| `POST /auth/send-otp` | Create and deliver an OTP |
| `GET/POST /auth/verify-otp` | Verify the code, then redirect to the callback with `?ott=` |
| `GET/POST /auth/logout` | Revoke the session |

JSON API (requires the `APP-ID` header):

| Route | Purpose |
| --- | --- |
| `POST /v1/auth/send-otp` | `{identifier}` → `{otp_hash, channel, expires_at}` |
| `POST /v1/auth/verify-otp` | `{identifier, otp_code, otp_hash}` → `{session_token, ott, redirect_url}` |
| `POST /v1/auth/verify-ott` | `{ott}` → credentials |
| `GET /v1/auth/check` | Managed mode redirects; custom mode returns JSON |
| `GET /v1/auth/check-login` | JSON login check; accepts the session cookie or `API-TOKEN` |
| `GET /v1/auth/me` | Current session |
| `POST/DELETE /v1/auth/logout` | Revoke the session |
| `GET /v1/user/sessions` | List my sessions |
| `DELETE /v1/user/sessions/:id` | Terminate one session |
| `DELETE /v1/user/sessions` | Terminate `{ids: [...]}` or all when the body is empty |

Admin panel (server-rendered, Persian):

- `/admin` – dashboard with per-app stats, daily chart and status breakdown
- `/admin/apps` – list/create/edit/delete apps with search and pagination
- `/admin/sessions` – list/filter/revoke sessions with search and pagination

Per-app configuration: managed/custom auth mode, OTP channel (`sms`, `email`, `both`), OTP length
(4–32) and charset (`numeric`, `alpha`, `alphanumeric`), OTP/OTT expiry, session
duration or infinite sessions, per-app session cookie (name, HttpOnly, Secure,
SameSite), same-origin return redirect, callback/webhook URLs, and the full auth
page appearance: layout, theme, font, text alignment, colors, logo mask and
position, selectable ArCaptcha or Google reCAPTCHA v2, ordered multi-banner
carousel, custom copy and a live preview in the admin form. Apps can also set a
comma-separated ISO 3166-1 alpha-2 country allowlist; an empty list permits all
countries.

Geo restriction and session location enrichment use a local MaxMind GeoLite2
City MMDB configured with `GEOIP_DATABASE_PATH`. Obtain and update the database
from MaxMind under its applicable terms. No IP lookup traffic leaves Aman. If a
country allowlist is set and no database is configured or the request IP cannot
be resolved, authentication is denied; a configured but unreadable database
prevents startup. Without an allowlist, a lookup failure only leaves unavailable
location fields empty. Google reCAPTCHA
uses checkbox v2 keys and Google's `siteverify` endpoint; only the CAPTCHA
response and configured secret are sent there. A CAPTCHA enabled without both
keys is rejected rather than bypassed.

When Aman is behind a reverse proxy, set `TRUSTED_PROXIES` to the proxy IPs or
CIDRs (comma-separated), for example `10.0.0.10,10.0.1.0/24`. The default is
empty, which disables forwarded-IP trust; untrusted `X-Forwarded-For` values
cannot influence geo-fencing.

## Architecture

```
internal/
  domains/            entities + services (business logic)
    services/         AuthService, AppService, SessionService, WebhookService
  ports/
    inbound/          contracts used by the HTTP adapters
    outbound/         contracts implemented by adapters
  adapters/
    inbound/http/     gin routes, handlers, middlewares, page components
    outbound/
      postgres/       gorm repositories (+ mock.go)
      geoip/          optional local MaxMind MMDB lookup
      notification/   SMS + email delivery (+ mock.go)
      otpgen/         OTP codes (+ mock.go)
      security/       JWT + auth hash (+ mock.go)
      webhook/        HTTP webhook delivery (+ mock.go)
```

Services are concrete structs: they receive outbound adapters and do not
implement interfaces themselves.

## Cookies and errors

Session cookies use the configured name (default `cas_cookie`), `Path=/`,
expiry, `HttpOnly`, `Secure` and `SameSite` settings. Aman can set a shared
parent-domain cookie when CAS and app hosts share a parent domain. For unrelated
domains the browser cannot accept Aman’s cookie for the app; exchange the OTT
server-side and set an app-owned cookie. API check-login accepts the token in
`API-TOKEN` for server-to-server checks. Invalid API requests still redirect to
`/not-found`, which nginx answers with the custom `nginx/404.html` page.

## Running

```bash
cp src/.env.example src/.env
docker compose up --build
```

Swagger is only registered outside `production` (`APP_ENV != production`).

OTP delivery falls back to logging in development. Set `SMTP_*` for email and
`SMS_GATEWAY_URL` for SMS in production.
