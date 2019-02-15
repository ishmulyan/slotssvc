# slotssvc

## API
```
go build ./cmd/api && ./api

curl -X POST 0.0.0.0:3000/api/machines/atkins-diet/spins \
    -d "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJ1c2VyMSIsImNoaXBzIjoxMDAwMCwiYmV0IjoxMDAwfQ.XyHvupzW6VnSUiZSC99nwRz_naB9WV5XDms_qje8Stk"
```

## Simulation
```
go build ./cmd/sim && ./sim
```

## Benchmarking
```
go test -benchmem -run=XXX ./pkg/atkinsdiet -bench BenchmarkSpin -v
```