# gRPC Http mTLS Hmac benchmarking

## Development

Run `make servers` and then `make test`

PS: Running locally both the server and benchmarking test will skew the results. Please change the `TARGET_HOST` to your remote server before running the code.

## Benchmark Setup

![Test Setup](/setup.png)


CY17 max diagnostic data upload size is 512MB. It is expected to be higher for CY17+. So 1MB is chosen.

## Benchmark Runs (using go test)

### Shinkansen Test (1 MB)

|   Benchmarks        |   N           |    Latency (Nanosecond/operation)  |    Bandwith (bytes/operation)  |   Memory (allocations/operation)  |
| ------------- |:-------------:| -----:| -----:|-----:|
|   Benchmark_MTLS_GRPC_HMAC_Protobuf-8          |   10      |     232329677 ns/op  |     128794 B/op    |    108 allocs/op  |
|   Benchmark_MTLS_GRPC_HMAC_Protobuf_Stream-8   |   10      |     223537533 ns/op  |     129587 B/op    |    112 allocs/op  |
|   Benchmark_MTLS_GRPC_Protobuf-8               |   10      |     192687629 ns/op  |     127694 B/op    |    103 allocs/op  |
|   Benchmark_MTLS_GRPC_Protobuf_Stream-8        |   10      |     214624319 ns/op  |     128323 B/op    |    110 allocs/op  |
|   Benchmark_MTLS_HTTP_HMAC_JSON-8              |   1       |    1214743491 ns/op  |     634212 B/op    |   1232 allocs/op  |
|   Benchmark_MTLS_HTTP_JSON-8                   |   1       |    1173860561 ns/op  |     604352 B/op    |   1373 allocs/op  |
|   Benchmark_TLS_HTTP_HMAC_JSON-8               |   1       |    1194234338 ns/op  |     603127 B/op    |   1367 allocs/op  |


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

## Load Test Runs (using modified fortio tool)

### RUN1 (Data: 100KB)

| Name | count | qps | conns | duration(s) | min(ms) | avg(ms) | p50(ms) | p75(ms) | p90(ms) | p99(ms) | p99.9(ms) | max(ms) |
| :--- | :---:| ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | 
| HTTP_TLS_HMAC | 1000 | 6.3 | 8 | 3.8 | 1132.282 | 1229.648 | 1256.46 | 1324.19 | 1364.83 | 1389.21 | 1391.65 | 1391.920 |
| HTTP_MTLS_HMAC | 1000 | 7.4 | 8 | 3.9 | 876.322 | 985.183 | 980.68 | 1063.61 | 1121.87 | 1156.82 | 1160.32 | 1160.709 |
| HTTP_MTLS | 1000 | 5 | 8 | 4.8 | 1035.332 | 1225.442 | 1517.67 | 1780.76 | 1938.61 | 2245.35 | 2315.08 | 2322.829 |
| GRPC_MTLS | 1000 | 21.3 | 8 | 3.4 | 134.099 | 358.557 | 287.5 |433.33 | 613.64 | 721.94 | 732.77 | 733.969 |
| GRPC_MTLS_HMAC | 1000 | 18.1 | 8 | 3.3 | 217.111 | 431.567 | 431.82 | 500 | 645.81 | 733.3 | 742.04 | 743.016 |

See [Report](https://github.com/telematicsct/grpc-benchmark/issues/2)

### RUN2 (Data: 100KB)

| Name | count | qps | conns | duration(s) | min(ms) | avg(ms) | p50(ms) | p75(ms) | p90(ms) | p99(ms) | p99.9(ms) | max(ms) |
| :--- | :---:| ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | 
| HTTP_TLS_HMAC | 1000 | 8 | 8 | 6 | 857.031 | 958.272 | 950.98 | 1000 | 1036.94 | 1059.1 | 1061.32 | 1061.564 |
| HTTP_MTLS_HMAC | 1000 | 8.3 | 8 | 5.8 | 850.158 | 901.073 | 925.08 | 964.17 | 987.62 | 1151.1 | 1276.63 | 1290.58 |
| HTTP_MTLS | 1000 | 5 | 8.3 | 5.8 | 866.123 | 933.91 | 945.08 | 986.27 | 1023.19 | 1054.49 | 1057.62 | 1057.963 |
| GRPC_MTLS | 1000 | 44.1 | 8 | 5.1 | 121.162 | 179.044 | 171.62 | 197.07 | 276.11 | 344.97 | 352.76 | 353.625 | 
| GRPC_MTLS_HMAC | 1000 | 37.8 | 8 | 5.3 | 123.326 | 205.094 | 180.4 | 234.78 | 300 | 750 | 813.94 | 821.039 |

See [Report](https://github.com/telematicsct/grpc-benchmark/issues/3)

### RUN3 (Data: 1mb)

| Name | count | qps | conns | duration(s) | min(ms) | avg(ms) | p50(ms) | p75(ms) | p90(ms) | p99(ms) | p99.9(ms) | max(ms) |
| :--- | :---:| ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | 
| HTTP_TLS_HMAC | 1000 | 8 | 8 | 5 | 850.321 | 955.098 | 957.9 | 1020.27 | 1068.92 | 1098.11 | 1101.03 | 1101.35 |
| HTTP_MTLS_HMAC | 1000 | 8.5 | 8 | 5 | 865.712 | 922.695 | 926.41 | 958.09 | 977.09 | 988.49 | 989.63 | 989.756 |
| HTTP_MTLS | 1000 | 7.2 | 8 | 5 | 929.99 | 1020.397 | 1013.5 | 1087.78 | 1132.34 | 1159.08 | 1161.76 | 1162.054 |
| GRPC_MTLS | 1000 | 36.6 | 8 | 5 | 136.585 | 215.997 | 204.82 | 261.45 | 295.42| 377.05 | 386.99 | 388.093 |
| GRPC_MTLS_HMAC | 1000 | 37.6 | 8 | 5 | 129.592 | 209.796 | 206.59 | 259.89 | 291.87 | 406.76 | 410.62 | 411.051 |

See [Report](https://github.com/telematicsct/grpc-benchmark/issues/4)

### Summary

As per the reports above
- for 1mb data, gRPC with mTLS (avg: 200ms) is faster than https with mTLS (avg: 850ms)
- for 100kb data, gRPC with mTLS (avg: 120ms) is faster than https with mTLS (avg: 850ms)
- 1mb or 100kb data, https has same average response time.
- gRPC mTLS with HMAC or without HMAC -- doesn't make a big difference.
- bandwidth and memory is quite stable for gRPC irrespective of network
