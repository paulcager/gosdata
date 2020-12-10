package osgrid

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestOSGridToLatLon(t *testing.T) {
	tests := []struct {
		name    string
		osgrid    string
		lat    float64
		lon   float64
		err bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := OSGridToLatLon(tt.osgrid)
			if (err != nil) != tt.err {
				t.Errorf("OSGridToLatLon() error = %v, wantErr %v", err, tt.err)
				return
			}
			if got != tt.lat {
				t.Errorf("OSGridToLatLon() got = %v, want %v", got, tt.lat)
			}
			if got1 != tt.lon {
				t.Errorf("OSGridToLatLon() got1 = %v, want %v", got1, tt.lon)
			}
		})
	}
}

// As you'd expect, the Otto implementation is about 1000 times slower than the Go implementation.
// 	BenchmarkOttoImpl-16               13500            440363 ns/op
// 	BenchmarkGoImpl-16              12458200               468 ns/op
// This doesn't much matter for my purposes, and we need to use it until I get the WGS84 datum implemented.
// Once that's done it can be kept as a known-good implementation I can test against.
func BenchmarkOttoImpl(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_,_, err := OSGridToLatLon("TL 44982 57869")
		assert.NoError(b, err)
	}
}

func BenchmarkGoImpl(b *testing.B) {
	o, err := ParseOsGridRef("TL 44982 57869")
	assert.NoError(b, err)

	for i := 0; i < b.N; i++ {
		_,_ = o.toLatLon()
	}
}

func BenchmarkAPICall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resp, err := http.Get("http://192.168.0.123:9090/?gridRef=SJ121689")
		require.NoError(b, err)
		reply, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		require.NoError(b, err)
		assert.Contains(b, string(reply), "53")
	}
}
