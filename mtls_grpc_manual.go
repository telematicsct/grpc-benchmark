package main

// import (
// 	"context"
// 	"crypto/rand"
// 	"crypto/tls"
// 	"crypto/x509"
// 	"flag"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials"
// 	"io/ioutil"

// 	"log"

// 	"github.com/montanaflynn/stats"
// 	pb "github.com/telematicsct/grpc-benchmark/dcm"
// 	"strings"
// 	"sync"
// 	"sync/atomic"
// 	"time"
// )

// var (
// 	concurrency = flag.Int("c", 10, "concurrency")
// 	total       = flag.Int("n", 1000, "total requests for all clients")
// 	host        = flag.String("s", "localhost:7900", "grpc host port")
// )

// func getClient() *grpc.ClientConn {
// 	certificate, err := tls.LoadX509KeyPair(
// 		"certs/client.crt",
// 		"certs/client.key",
// 	)

// 	certPool := x509.NewCertPool()
// 	bs, err := ioutil.ReadFile("certs/ca.crt")
// 	if err != nil {
// 		log.Fatalf("failed to read ca cert: %s", err)
// 	}

// 	ok := certPool.AppendCertsFromPEM(bs)
// 	if !ok {
// 		log.Fatal("failed to append certs")
// 	}

// 	transportCreds := credentials.NewTLS(&tls.Config{
// 		ServerName:   "localhost",
// 		Certificates: []tls.Certificate{certificate},
// 		RootCAs:      certPool,
// 	})

// 	opts := []grpc.DialOption{
// 		grpc.WithTransportCredentials(transportCreds),
// 		// grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")),
// 	}
// 	conn, err := grpc.Dial("localhost:7900", opts...)
// 	if err != nil {
// 		log.Fatalf("failed to dial server: %s", err)
// 	}

// 	return conn
// }

// func getPayload() []byte {
// 	//100000 - 100kb
// 	payload := make([]byte, 100000)
// 	if _, err := rand.Read(payload); err != nil {
// 		log.Fatalf("payload error %v", err)
// 	}
// 	return payload
// }

// func testManual() {
// 	flag.Parse()
// 	n := *concurrency
// 	m := *total / n

// 	selected := -1
// 	servers := strings.Split(*host, ",")
// 	sNum := len(servers)

// 	log.Printf("Servers: %+v\n\n", servers)

// 	log.Printf("concurrency: %d\nrequests per client: %d\n\n", n, m)

// 	payload := getPayload()
// 	data := &pb.DiagRecorderData{CanId: 123456789, Payload: &pb.Payload{Body: payload}}

// 	var wg sync.WaitGroup
// 	wg.Add(n * m)

// 	var trans uint64
// 	var transOK uint64

// 	d := make([][]int64, n, n)

// 	//it contains warmup time but we can ignore it
// 	totalT := time.Now().UnixNano()
// 	for i := 0; i < n; i++ {
// 		dt := make([]int64, 0, m)
// 		d = append(d, dt)
// 		selected = (selected + 1) % sNum

// 		go func(i int, selected int) {
// 			conn := getClient()
// 			c := pb.NewDCMServiceClient(conn)

// 			//warmup
// 			for j := 0; j < 5; j++ {
// 				c.DiagnosticData(context.Background(), data)
// 			}

// 			for j := 0; j < m; j++ {
// 				t := time.Now().UnixNano()
// 				reply, err := c.DiagnosticData(context.Background(), data)
// 				t = time.Now().UnixNano() - t

// 				d[i] = append(d[i], t)

// 				if err == nil && reply.Code == 200 {
// 					atomic.AddUint64(&transOK, 1)
// 				}

// 				atomic.AddUint64(&trans, 1)
// 				wg.Done()
// 			}

// 			conn.Close()

// 		}(i, selected)

// 	}

// 	wg.Wait()
// 	totalT = time.Now().UnixNano() - totalT
// 	totalT = totalT / 1000000
// 	log.Printf("took %d ms for %d requests\n", totalT, n*m)

// 	totalD := make([]int64, 0, n*m)
// 	for _, k := range d {
// 		totalD = append(totalD, k...)
// 	}
// 	totalD2 := make([]float64, 0, n*m)
// 	for _, k := range totalD {
// 		totalD2 = append(totalD2, float64(k))
// 	}

// 	mean, _ := stats.Mean(totalD2)
// 	median, _ := stats.Median(totalD2)
// 	max, _ := stats.Max(totalD2)
// 	min, _ := stats.Min(totalD2)
// 	p99, _ := stats.Percentile(totalD2, 99.9)

// 	log.Printf("sent     requests    : %d\n", n*m)
// 	log.Printf("received requests    : %d\n", atomic.LoadUint64(&trans))
// 	log.Printf("received requests_OK : %d\n", atomic.LoadUint64(&transOK))
// 	log.Printf("throughput  (TPS)    : %d\n", int64(n*m)*1000/totalT)
// 	// log.Printf("mean: %.f ns, median: %.f ns, max: %.f ns, min: %.f ns, p99: %.f ns\n", mean, median, max, min, p99)
// 	log.Printf("mean: %d ms, median: %d ms, max: %d ms, min: %d ms, p99: %d ms\n", int64(mean/1000000), int64(median/1000000), int64(max/1000000), int64(min/1000000), int64(p99/1000000))

// }

// func testManualStream() {
// 	flag.Parse()
// 	n := *concurrency
// 	m := *total / n

// 	selected := -1
// 	servers := strings.Split(*host, ",")
// 	sNum := len(servers)

// 	log.Printf("Servers: %+v\n\n", servers)

// 	log.Printf("concurrency: %d\nrequests per client: %d\n\n", n, m)

// 	payload := getPayload()
// 	data := &pb.DiagRecorderData{CanId: 123456789, Payload: &pb.Payload{Body: payload}}

// 	var wg sync.WaitGroup
// 	wg.Add(n * m)

// 	var trans uint64
// 	var transOK uint64

// 	d := make([][]int64, n, n)

// 	//it contains warmup time but we can ignore it
// 	totalT := time.Now().UnixNano()
// 	for i := 0; i < n; i++ {
// 		dt := make([]int64, 0, m)
// 		d = append(d, dt)
// 		selected = (selected + 1) % sNum

// 		go func(i int, selected int) {
// 			conn := getClient()
// 			c := pb.NewDCMServiceClient(conn)

// 			//warmup
// 			for j := 0; j < 5; j++ {
// 				s, err := c.DiagnosticDataStream(context.Background())
// 				if err != nil {
// 					log.Fatalf("%v.DiagnosticData(_) = _, %v", c, err)
// 				}

// 				if err := s.Send(data); err != nil {
// 					log.Fatalf("send error %v", err)
// 				}

// 				reply, err := s.CloseAndRecv()
// 				if err != nil {
// 					log.Fatalf("Warmup: %v.CloseAndRecv() got error %v, want %v", s, err, nil)
// 				}
// 				if reply.Code != 200 || reply.Message != "Done" {
// 					log.Fatalf("Warmup:  grpc response is wrong: %v", reply)
// 				}
// 			}

// 			for j := 0; j < m; j++ {
// 				t := time.Now().UnixNano()
// 				stream, err := c.DiagnosticDataStream(context.Background())
// 				if err != nil {
// 					log.Fatalf("%v.DiagnosticData(_) = _, %v", c, err)
// 				}

// 				if err := stream.Send(data); err != nil {
// 					log.Fatalf("send error %v", err)
// 				}

// 				reply, err := stream.CloseAndRecv()
// 				if err != nil {
// 					log.Fatalf("Real: %v.CloseAndRecv() got error %v, want %v", stream, err, nil)
// 				}
// 				t = time.Now().UnixNano() - t

// 				d[i] = append(d[i], t)

// 				if reply.Code != 200 || reply.Message != "Done" {
// 					log.Fatalf("Real:  grpc response is wrong: %v", reply)
// 				} else {
// 					atomic.AddUint64(&transOK, 1)
// 				}

// 				atomic.AddUint64(&trans, 1)
// 				wg.Done()
// 			}

// 			conn.Close()

// 		}(i, selected)

// 	}

// 	wg.Wait()
// 	totalT = time.Now().UnixNano() - totalT
// 	totalT = totalT / 1000000
// 	log.Printf("took %d ms for %d requests\n", totalT, n*m)

// 	totalD := make([]int64, 0, n*m)
// 	for _, k := range d {
// 		totalD = append(totalD, k...)
// 	}
// 	totalD2 := make([]float64, 0, n*m)
// 	for _, k := range totalD {
// 		totalD2 = append(totalD2, float64(k))
// 	}

// 	mean, _ := stats.Mean(totalD2)
// 	median, _ := stats.Median(totalD2)
// 	max, _ := stats.Max(totalD2)
// 	min, _ := stats.Min(totalD2)
// 	p99, _ := stats.Percentile(totalD2, 99.9)

// 	log.Printf("sent     requests    : %d\n", n*m)
// 	log.Printf("received requests    : %d\n", atomic.LoadUint64(&trans))
// 	log.Printf("received requests_OK : %d\n", atomic.LoadUint64(&transOK))
// 	log.Printf("throughput  (TPS)    : %d\n", int64(n*m)*1000/totalT)
// 	// log.Printf("mean: %.f ns, median: %.f ns, max: %.f ns, min: %.f ns, p99: %.f ns\n", mean, median, max, min, p99)
// 	log.Printf("mean: %d ms, median: %d ms, max: %d ms, min: %d ms, p99: %d ms\n", int64(mean/1000000), int64(median/1000000), int64(max/1000000), int64(min/1000000), int64(p99/1000000))

// }
