# gRPC Http mTLS Hmac benchmarking

## Benchmark Runs (using go test)

### Data: 100KB

|   Benchmarks        |   N           |    Latency (Nanosecond/operation)  |    Bandwith (bytes/operation)  |   Memory (allocations/operation)  |
| ------------- |:-------------:| -----:| -----:|-----:|
|   Benchmark_MTLS_GRPC_HMAC_Protobuf-8          |   500      |     2504702 ns/op    |     128794 B/op    |     97 allocs/op |
|   Benchmark_MTLS_GRPC_HMAC_Protobuf_Stream-8   |   500      |     2533754 ns/op    |     129169 B/op    |    106 allocs/op |
|   Benchmark_MTLS_GRPC_Protobuf-8               |   500      |     2508642 ns/op    |     127944 B/op    |     91 allocs/op |
|   Benchmark_MTLS_GRPC_Protobuf_Stream-8        |   500      |     2529270 ns/op    |     128332 B/op    |    100 allocs/op |
|   Benchmark_MTLS_HTTP_HMAC_JSON-8              |   300      |     4572717 ns/op    |     249072 B/op    |    292 allocs/op |
|   Benchmark_MTLS_HTTP_JSON-8                   |   300      |     4783064 ns/op    |     262402 B/op    |    352 allocs/op |
|   Benchmark_TLS_HTTP_HMAC_JSON-8               |   300      |     4870114 ns/op    |     225274 B/op    |    219 allocs/op |


### Data: 1MB

|   Benchmarks        |   N           |    Latency (Nanosecond/operation)  |    Bandwith (bytes/operation)  |   Memory (allocations/operation)  |
| ------------- |:-------------:| -----:| -----:|-----:|
|   Benchmark_MTLS_GRPC_HMAC_Protobuf-8          |   100      |     17431225 ns/op   |     1030183 B/op    |   106 allocs/op |
|   Benchmark_MTLS_GRPC_HMAC_Protobuf_Stream-8   |   100      |     17682164 ns/op   |     1030704 B/op    |   120 allocs/op |
|   Benchmark_MTLS_GRPC_Protobuf-8               |   100      |     18944755 ns/op   |     1029412 B/op    |   103 allocs/op |
|   Benchmark_MTLS_GRPC_Protobuf_Stream-8        |   100      |     18849443 ns/op   |     1029891 B/op    |   112 allocs/op |
|   Benchmark_MTLS_HTTP_HMAC_JSON-8              |    50      |     28537868 ns/op   |     3990679 B/op    |   495 allocs/op |
|   Benchmark_MTLS_HTTP_JSON-8                   |    50      |     25730243 ns/op   |     3534994 B/op    |   438 allocs/op |
|   Benchmark_TLS_HTTP_HMAC_JSON-8               |    50      |     26732076 ns/op   |     3780368 B/op    |   258 allocs/op |

## Load Test using modified fortio

### RUN1 (Data: 100KB)

| Name | count | qps | conns | duration(s) | min(ms) | avg(ms) | p50(ms) | p75(ms) | p90(ms) | p99(ms) | p99.9(ms) | max(ms) |
| :--- | :---:| ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | 
| HTTP_TLS_HMAC | 1000 | 6.3 | 8 | 3.8 | 1132.282 | 1229.648 | 1256.46 | 1324.19 | 1364.83 | 1389.21 | 1391.65 | 1391.920 |
| HTTP_MTLS_HMAC | 1000 | 7.4 | 8 | 3.9 | 876.322 | 985.183 | 980.68 | 1063.61 | 1121.87 | 1156.82 | 1160.32 | 1160.709 |
| HTTP_MTLS | 1000 | 5 | 8 | 4.8 | 1035.332 | 1225.442 | 1517.67 | 1780.76 | 1938.61 | 2245.35 | 2315.08 | 2322.829 |
| GRPC_MTLS | 1000 | 21.3 | 8 | 3.4 | 134.099 | 358.557 | 287.5 |433.33 | 613.64 | 721.94 | 732.77 | 733.969 |
| GRPC_MTLS_HMAC | 1000 | 18.1 | 8 | 3.3 | 217.111 | 431.567 | 431.82 | 500 | 645.81 | 733.3 | 742.04 | 743.016 |

### RUN2 (Data: 100KB)

