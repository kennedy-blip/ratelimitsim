# RateLimitSim

Concurrent API rate-limit, token-bucket and burst-protection simulator. Vanilla HTML/CSS/JavaScript frontend; Go standard-library Vercel functions backend.

## Included
- Seeded payments, public-search and heavy-export policies.
- Per-client token buckets with proportional elapsed-time refill.
- Concurrent request workers with mutex-serialized token consumption.
- 200/429 traces, remaining tokens, retry-after estimates and aggregate metrics.
- Policy tuning, bucket resets and offline browser fallback.

## Deploy
Import `kennedy-blip/ratelimitsim` into Vercel with the repository root as the project root. The project uses `vercel.json` routes to serve `public/` and expose `/api/policies` and `/api/simulate`. No Node build step is required.

## Local checks
```sh
gofmt -w api/*.go
go vet ./api/...
go test ./...
```

## API
- `GET /api/policies` — seeded policies and bucket snapshots.
- `POST /api/policies` — update `id`, `capacity`, `refillPerSec`, `costPerReq`; set `reset:true` to clear buckets, optionally scoped by `clientKey`.
- `POST /api/simulate` — body: `{"policyId":"payments-stk-api","clientKey":"tenant_ke_01","requestCount":30,"concurrency":8,"intervalMs":0}`.

Simulation limits: 1–100 requests, 1–16 workers and 0–10000 ms interval.

## Important deployment limitation
The included store is process-local memory. Vercel serverless instances can scale independently and restart, so state is not globally consistent across instances. This is appropriate for a simulator/demo; use a shared atomic store such as Redis for real distributed enforcement.
