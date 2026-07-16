# AMJ HUB — Professional Cinematography & Visual Production Website

A complete, production-ready Go web application for AMJ HUB, a professional
cinematography, photography, videography, and drone services brand.

---

## Tech Stack

| Layer    | Technology                                     |
|----------|------------------------------------------------|
| Backend  | Go (Golang) · `net/http` standard library      |
| Frontend | HTML5 · CSS3 (Variables, Grid, Flexbox)        |
| Fonts    | Google Fonts — Playfair Display + Inter        |
| No deps  | Zero external Go packages required             |

---

## Project Structure

```
amjhub/
├── main.go                   # Server entry point & router
├── go.mod                    # Go module definition
├── backend/
│   └── handlers.go           # Route handlers (home page, contact form)
├── static/
│   ├── css/
│   │   └── style.css         # Full design system (variables, components)
│   └── js/
│       └── main.js           # Nav, animations, AJAX form submission
└── templates/
    └── index.html            # Main landing page template
```

---

## Getting Started

### Prerequisites
- Go 1.21+ installed ([https://go.dev/dl/](https://go.dev/dl/))

### Run Locally

```bash
# Clone / navigate to project
cd amjhub

# Run the server
go run main.go
```

Open your browser at: **http://localhost:8080**

### Build for Production

```bash
# Build binary
go build -o amjhub-server .

# Run binary
./amjhub-server
```

---

## Features

### Frontend
- ✅ Sticky, scroll-aware navigation bar with mobile hamburger menu
- ✅ Cinematic hero section with letterbox overlay, ambient grid, animated headline
- ✅ About section with company story and value proposition
- ✅ Services grid — 3 core services + 8 additional service tiles
- ✅ "Why Choose Us" section with feature highlights
- ✅ Full contact & inquiry form with inline AJAX response (no page reload)
- ✅ Responsive footer with links and company info
- ✅ Scroll-reveal animations using IntersectionObserver
- ✅ Animated stat counters
- ✅ CSS custom properties for easy theme editing
- ✅ Fully responsive down to 320px
- ✅ `prefers-reduced-motion` respected

### Backend (Go)
- ✅ Serves HTML template via `html/template`
- ✅ Serves all static assets (CSS/JS) from `/static/`
- ✅ `POST /contact` endpoint:
  - Parses `application/x-www-form-urlencoded` form data
  - Validates: full name (required, ≥2 chars) + email (required, regex) + message (required)
  - Logs structured inquiry to terminal
  - Returns JSON `{ success: true/false, message/error }` for AJAX handling
- ✅ Health check endpoint at `GET /health`
- ✅ 404 handling

---

## Customising the Theme

All design tokens live in `static/css/style.css` at the top under `:root {}`:

```css
:root {
  --color-gold:    #C9A84C;   /* ← Change brand accent color here */
  --color-black:   #0A0A0A;   /* ← Background */
  --font-display:  'Playfair Display', serif;
  --font-body:     'Inter', sans-serif;
  /* ... */
}
```

---

## Email Configuration (Contact Form Delivery)

Contact form submissions are emailed via SMTP using only Go's standard
library (`net/smtp` — no external dependencies). All settings, including
**the recipient's inbox**, are read from environment variables so they can
be changed at any time without touching the code.

### Required Environment Variables

| Variable                  | Description                                              | Example                     |
|----------------------------|------------------------------------------------------------|------------------------------|
| `CONTACT_RECIPIENT_EMAIL`  | **The inbox that receives inquiries.** Change this anytime. | `info@amjhub.com`           |
| `SMTP_HOST`                 | SMTP server hostname                                       | `smtp.gmail.com`            |
| `SMTP_PORT`                 | SMTP server port                                            | `587`                       |
| `SMTP_USERNAME`             | Account used to authenticate and send the email             | `bookings@amjhub.com`       |
| `SMTP_PASSWORD`             | Password / app-password for that account                    | `xxxxxxxxxxxxxxxx`          |
| `EMAIL_FROM` *(optional)*  | "From" address shown to the recipient (defaults to `SMTP_USERNAME`) | `info@amjhub.com`   |

A template is provided at `.env.example` — copy it to `.env` and fill in real values.

### Setting Environment Variables Locally

Go does **not** auto-load `.env` files. Either export variables manually,
or use a tool like `direnv` / `godotenv`. Simplest option for local testing:

```bash
export CONTACT_RECIPIENT_EMAIL="info@amjhub.com"
export SMTP_HOST="smtp.gmail.com"
export SMTP_PORT="587"
export SMTP_USERNAME="bookings@amjhub.com"
export SMTP_PASSWORD="your-app-password"

go run main.go
```

Or inline for a single run:

```bash
CONTACT_RECIPIENT_EMAIL="info@amjhub.com" \
SMTP_HOST="smtp.gmail.com" \
SMTP_PORT="587" \
SMTP_USERNAME="bookings@amjhub.com" \
SMTP_PASSWORD="your-app-password" \
go run main.go
```

### In Production

Set these as environment variables in your hosting platform's dashboard
(Render, Railway, Fly.io, a Docker `--env-file`, or a systemd `Environment=`
directive). **Never commit real credentials to source control** — `.env`
should stay in `.gitignore`.

### Changing the Recipient Later

To redirect inquiries to a new inbox, the client only needs to update
**one value** — `CONTACT_RECIPIENT_EMAIL` — in their hosting platform's
environment settings, then restart the server. No code or redeploy required
beyond a restart.

### Gmail Note

Gmail requires an **App Password** (not your normal login password) when
authenticating via SMTP. Generate one at: Google Account → Security →
2-Step Verification → App Passwords.

### Fallback Behavior

If SMTP environment variables are missing or incomplete, the server will
**not crash** — it logs a clear warning to the terminal and still records
the full inquiry in the console log, so no submission is lost during setup.

---

## Production Checklist

- [ ] Set `CONTACT_RECIPIENT_EMAIL`, `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD` as real environment variables
- [ ] Replace placeholder phone / email in `templates/index.html`
- [ ] Add real social media URLs in the contact section
- [ ] Add real portfolio images to `/static/images/`
- [ ] Set `PORT` via environment variable (update `main.go`)
- [ ] Enable HTTPS (via reverse proxy like Nginx or Caddy)
- [ ] Add CSRF protection for the contact form
- [ ] Add rate limiting to `POST /contact`

---

## License

© AMJ HUB. All rights reserved.
