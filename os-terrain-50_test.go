package gosdata

import (
	"math"
	"testing"

	"github.com/paulcager/osgridref"
	"github.com/prometheus/client_golang/prometheus"
	prommodel "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	NY21 = osgridref.OsGridRef{
		Easting:  320000,
		Northing: 510000,
	}

	ts = NewTileServer("test-data")
)

func Test_toTile(t *testing.T) {
	tile, err := getTile(NY21, "test-data/ny/ny21_OST50GRID_20200303.zip")
	require.NoError(t, err)

	_ = tile
}

func Test_LoadFile(t *testing.T) {
	tile, err := ts.LoadTile(NY21)
	require.NoError(t, err)

	assert.EqualValues(t, 4028, tile[0][0])
	assert.EqualValues(t, 2606, tile[199][199])
}

func TestErrors(t *testing.T) {
	_, err := ts.Height("SY 21725 68352") // Valid GR, but in middle of the sea
	assert.Error(t, err)

	_, err = ts.Height("NY")
	assert.Error(t, err)

	_, err = ts.Height("NY123X456")
	assert.Error(t, err)
}

func Test_sanity(t *testing.T) {
	// Test against known heights. Specific files have been loaded into `test-data` directory.

	assertBetween(t, 175, 180, ts.MustHeight("NY 30898 17869"))
	assertBetween(t, 5, 40, ts.MustHeight("NY 03165 02149"))
	assertBetween(t, 35, 40, ts.MustHeight("NY 09462 02315"))
	assertBetween(t, 890, 900, ts.MustHeight("NY 21108 10343")) // Great Gable, 899
	assertBetween(t, 340, 350, ts.MustHeight("SE 70166 94892"))
	assertBetween(t, 540, 550, ts.MustHeight("SK 11062 83392")) // Lord's Seat
	assertBetween(t, 210, 220, ts.MustHeight("SK 14319 83727"))
	assertBetween(t, 195, 205, ts.MustHeight("TQ 37430 07467"))
}

func assertBetween(t *testing.T, min, max, actual float64) {
	assert.GreaterOrEqual(t, actual, min)
	assert.LessOrEqual(t, actual, max)
}

func Test_trunc10k(t *testing.T) {
	assert.EqualValues(t, 320000, trunc10k(321234))
	assert.EqualValues(t, 510000, trunc10k(519876))
}

func Test_caching(t *testing.T) {
	initReq := readCounter(t, requestCount)
	initRead := readCounter(t, readCount)
	initBad := readCounter(t, badCount)

	// Dedicated ts, so other tests can't interfere
	ts := NewTileServer("test-data")

	// Bad request should increment all 3 counters
	_, err :=ts.Height("SY 21725 68352") // Middle of sea - no tile
	assert.Error(t, err)
	assert.EqualValues(t, initReq+1, readCounter(t, requestCount))
	assert.EqualValues(t, initRead+1, readCounter(t, readCount))
	assert.EqualValues(t, initBad+1, readCounter(t, badCount))

	// Subsequent bad request should not re-read
	_, err =ts.Height("SY 21725 68352")
	assert.Error(t, err)
	assert.EqualValues(t, initReq+2, readCounter(t, requestCount))
	assert.EqualValues(t, initRead+1, readCounter(t, readCount))
	assert.EqualValues(t, initBad+1, readCounter(t, badCount))

	// Good request should increment requests and reads
	_, err =ts.Height("SK 11062 83392")
	assert.NoError(t, err)
	assert.EqualValues(t, initReq+3, readCounter(t, requestCount))
	assert.EqualValues(t, initRead+2, readCounter(t, readCount))
	assert.EqualValues(t, initBad+1, readCounter(t, badCount))

	// Subsequent good request should not re-read
	_, err =ts.Height("SK 11099 83344")
	assert.NoError(t, err)
	assert.EqualValues(t, initReq+4, readCounter(t, requestCount))
	assert.EqualValues(t, initRead+2, readCounter(t, readCount))
	assert.EqualValues(t, initBad+1, readCounter(t, badCount))
}

func readCounter(t *testing.T, c prometheus.Counter) int {
	var m prommodel.Metric
	require.NoError(t, c.Write(&m))
	return int(math.Round(m.Counter.GetValue()))
}
