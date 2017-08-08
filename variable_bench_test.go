package bearcub_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/gernest/bearcub"
)

// TODO: Add benchmarks
/*

Instructions:


*/

func BenchmarkReplaceString(b *testing.B) {
	jr, err := bearcub.NewJSONReplacer([]byte("{}"))
	if err != nil {
		b.Fatal(err)
	}
	var buf bytes.Buffer
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		s := []byte("{number}")
		e := fmt.Sprint(i)
		jr.O["number"] = i
		buf.Reset()
		b.StartTimer()
		err := bearcub.ReplaceString(&buf, s, jr.Replace)
		if err != nil {
			b.Fatal(err)
		}
		b.StopTimer()
		if buf.String() != e {
			b.Fatalf("expected %s got %s", e, buf.String())
		}
	}
}
