# Requirements
## Must (M)

- User authentication (register / login / logout) using bcrypt-hashed passwords and server-side sessions.

- CRUD for workouts (type, duration, distance, calories, notes, timestamp).

- CRUD for daily metrics: weight, resting heart rate, steps, sleep_hours.

- Use Gin for HTTP API + server-side rendered pages (multitemplate) for web UI.

- Persist data 

- Data validation and input sanitation on server side.

- Migrations and schema versioning (go-migrate or embed migrations).

- Unit and integration tests for core model and handler logic.

## Should (S)

- Sessions securely stored (cookie-based with HTTPOnly, Secure flags) and session expiry.

- User profile (height, gender, DOB) to compute derived metrics (BMI, age-based targets).

- Basic reporting endpoints (weekly totals, best runs) and CSV export.

- Device import endpoint to accept CSV/JSON from phone or wearable vendors (manual import MVP).

- Rate limiting and brute-force protection on login.

## Could (C)

- Background job to summarize daily activity and send digest (simple goroutine + cron-like scheduler).

- Social features: share workouts, public profile (opt-in).

- Cloud sync option (post-MVP) — server exposes JSON API for mobile clients.

## Won't (W)

`Won't integrate with third-party OAuth or native mobile SDKs in MVP.`

`Won't support realtime websockets or live telemetry in MVP.`