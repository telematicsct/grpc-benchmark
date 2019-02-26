# gRPC Http mTLS Hmac benchmarking

### 100KB data

|   Benchmarks        |   N           |    Latency (Nanosecond/operation)  |    Bandwith (bytes/operation)  |   Memory (allocations/operation)  |
| ------------- |:-------------:| -----:| -----:|-----:|
|   Benchmark_MTLS_GRPC_HMAC_Protobuf-8          |   500      |     2504702 ns/op    |     128794 B/op    |     97 allocs/op |
|   Benchmark_MTLS_GRPC_HMAC_Protobuf_Stream-8   |   500      |     2533754 ns/op    |     129169 B/op    |    106 allocs/op |
|   Benchmark_MTLS_GRPC_Protobuf-8               |   500      |     2508642 ns/op    |     127944 B/op    |     91 allocs/op |
|   Benchmark_MTLS_GRPC_Protobuf_Stream-8        |   500      |     2529270 ns/op    |     128332 B/op    |    100 allocs/op |
|   Benchmark_MTLS_HTTP_HMAC_JSON-8              |   300      |     4572717 ns/op    |     249072 B/op    |    292 allocs/op |
|   Benchmark_MTLS_HTTP_JSON-8                   |   300      |     4783064 ns/op    |     262402 B/op    |    352 allocs/op |
|   Benchmark_TLS_HTTP_HMAC_JSON-8               |   300      |     4870114 ns/op    |     225274 B/op    |    219 allocs/op |


### 1MB data

|   Benchmarks        |   N           |    Latency (Nanosecond/operation)  |    Bandwith (bytes/operation)  |   Memory (allocations/operation)  |
| ------------- |:-------------:| -----:| -----:|-----:|
|   Benchmark_MTLS_GRPC_HMAC_Protobuf-8          |   100      |     17431225 ns/op   |     1030183 B/op    |   106 allocs/op |
|   Benchmark_MTLS_GRPC_HMAC_Protobuf_Stream-8   |   100      |     17682164 ns/op   |     1030704 B/op    |   120 allocs/op |
|   Benchmark_MTLS_GRPC_Protobuf-8               |   100      |     18944755 ns/op   |     1029412 B/op    |   103 allocs/op |
|   Benchmark_MTLS_GRPC_Protobuf_Stream-8        |   100      |     18849443 ns/op   |     1029891 B/op    |   112 allocs/op |
|   Benchmark_MTLS_HTTP_HMAC_JSON-8              |    50      |     28537868 ns/op   |     3990679 B/op    |   495 allocs/op |
|   Benchmark_MTLS_HTTP_JSON-8                   |    50      |     25730243 ns/op   |     3534994 B/op    |   438 allocs/op |
|   Benchmark_TLS_HTTP_HMAC_JSON-8               |    50      |     26732076 ns/op   |     3780368 B/op    |   258 allocs/op |