# gRPC Http mTLS Hmac benchmarking


GO111MODULE=on go test -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: github.com/telematicsct/grpc-benchmark
|            Benchmarks                        |   N        | Latency (Nanosecond/operation) | Bandwith (bytes/operation) | Memory (allocations/operation) |
|----------------------------------------------|:----------:|:------------------------------:|:--------------------------:|:------------------------------:|
| Benchmark_MTLS_GRPC_HMAC_Protobuf-8          |   500      |     2504702 ns/op    |     128794 B/op    |     97 allocs/op |
| Benchmark_MTLS_GRPC_HMAC_Protobuf_Stream-8   |   500      |     2533754 ns/op    |     129169 B/op    |    106 allocs/op |
| Benchmark_MTLS_GRPC_Protobuf-8               |   500      |     2508642 ns/op    |     127944 B/op    |     91 allocs/op |
| Benchmark_MTLS_GRPC_Protobuf_Stream-8        |   500      |     2529270 ns/op    |     128332 B/op    |    100 allocs/op |
| Benchmark_MTLS_HTTP_HMAC_JSON-8              |   300      |     4572717 ns/op    |     249072 B/op    |    292 allocs/op |
| Benchmark_MTLS_HTTP_JSON-8                   |   300      |     4783064 ns/op    |     262402 B/op    |    352 allocs/op |
| Benchmark_TLS_HTTP_HMAC_JSON-8               |   300      |     4870114 ns/op    |     225274 B/op    |    219 allocs/op |
PASS
ok      github.com/telematicsct/grpc-benchmark  15.653s
