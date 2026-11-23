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
| `GET /v1/auth/check` | Redirect with a fresh OTT or to login |
| `GET /v1/auth/me` | Current session |
| `POST/DELETE /v1/auth/logout` | Revoke the session |
| `GET /v1/user/sessions` | List my sessions |
| `DELETE /v1/user/sessions/:id` | Terminate one session |
| `DELETE /v1/user/sessions` | Terminate `{ids: [...]}` or all when the body is empty |

Admin panel (server-rendered, Persian):

- `/admin` – dashboard with per-app stats, daily chart and status breakdown
- `/admin/apps` – list/create/edit/delete apps with search and pagination
- `/admin/sessions` – list/filter/revoke sessions with search and pagination

Per-app configuration: OTP channel (`sms`, `email`, `both`), OTP length
(4–32) and charset (`numeric`, `alpha`, `alphanumeric`), OTP/OTT expiry, session
duration or infinite sessions, per-app session cookie (name, HttpOnly, Secure,
SameSite), same-origin return redirect, callback/webhook URLs, and the full auth
page appearance: layout, theme, font, text alignment, colors, logo mask and
position, captcha, ordered multi-banner carousel, custom copy and a live preview
in the admin form.

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
      notification/   SMS + email delivery (+ mock.go)
      otpgen/         OTP codes (+ mock.go)
      security/       JWT + auth hash (+ mock.go)
      webhook/        HTTP webhook delivery (+ mock.go)
```

Services are concrete structs: they receive outbound adapters and do not
implement interfaces themselves.

## Errors and 404

The app does not return JSON on errors. It redirects to `/not-found`, which
nginx answers with the custom `nginx/404.html` page.

## Running

```bash
cp src/.env.example src/.env
docker compose up --build
```

Swagger is only registered outside `production` (`APP_ENV != production`).

OTP delivery falls back to logging in development. Set `SMTP_*` for email and
`SMS_GATEWAY_URL` for SMS in production.
