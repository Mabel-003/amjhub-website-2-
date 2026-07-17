# Build and run with Docker

Build the image locally (from project root):

```bash
docker build -t amjhub:latest .
```

Run with environment variables (example):

```bash
docker run -p 8000:8000 \
  -e CONTACT_RECIPIENT_EMAIL=info@amjhub.com \
  -e SMTP_HOST=smtp.example.com \
  -e SMTP_PORT=587 \
  -e SMTP_USERNAME=send@example.com \
  -e SMTP_PASSWORD=supersecret \
  amjhub:latest
```

Recommended production steps:
- Use a secrets manager or platform-provided env vars (Render, Railway, ECS Task Definitions).
- Run behind a reverse proxy / load balancer with TLS (Traefik, nginx, or cloud load balancer).
- Configure health checks to point to `/health` for container orchestration.
- Limit `Access-Control-Allow-Origin` to your domain in `backend.ContactHandler`.
