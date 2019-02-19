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


2019/02/19 23:46:45 Simulations started...
2019/02/19 23:46:45 concurrency:        20
2019/02/19 23:46:45 spins:              1e+09
2019/02/19 23:46:45 bet:                1
2019/02/19 23:46:45 lines:              20
2019/02/20 00:13:26 Simulation finished
2019/02/20 00:13:26 Total Bets:         20000000000
2019/02/20 00:13:26 Total Wins:         19405710383
2019/02/20 00:13:26 Total RTP:          0.970286
2019/02/20 00:13:26 Line Pays RTP:      0.634474
2019/02/20 00:13:26 Scatter RTP:        0.069751
2019/02/20 00:13:26 Bonus RTP:          0.266060
2019/02/20 00:13:26 Sum RTP:            0.970286
```

## Benchmarking
```
go test -benchmem -run=XXX ./pkg/atkinsdiet -bench BenchmarkSpin -v
```