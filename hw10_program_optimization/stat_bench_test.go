package hw10programoptimization

import (
	"archive/zip"
	"testing"
)

// go test -bench=BenchmarkGetDomainStat -benchmem -benchtime 10s
// benchstat -delta-test none old new
// go test -bench=BenchmarkGetDomainStat -benchmem -benchtime 10s -cpuprofile=cpu.out .
// go test -bench=BenchmarkGetDomainStat -benchmem -benchtime 10s -memprofile=mem.out .
// go tool pprof -http=":8090" hw10_program_optimization.test mem.out
// go tool pprof -http=":8091" hw10_program_optimization.test cpu.out
func BenchmarkGetDomainStat(b *testing.B) {
	r, _ := zip.OpenReader("testdata/users.dat.zip")
	defer r.Close()
	data, _ := r.File[0].Open()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetDomainStat(data, "biz")
	}
}
