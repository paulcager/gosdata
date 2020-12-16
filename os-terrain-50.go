package gosdata

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	_ "io/ioutil"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/paulcager/osgridref"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestCount = promauto.NewCounter(prometheus.CounterOpts{
		Subsystem: "terr50",
		Name:      "request_counter",
		Help:      "Number of times requested a terrain-50 tile",
	})
	readCount = promauto.NewCounter(prometheus.CounterOpts{
		Subsystem: "terr50",
		Name:      "retrieval_counter",
		Help:      "Number of times read a terrain-50 tile from a zip file",
	})
	badCount = promauto.NewCounter(prometheus.CounterOpts{
		Subsystem: "terr50",
		Name:      "bad_counter",
		Help:      "Number of times requested an unavailable terrain-50 tile",
	})
)

type Tile [200][200]int16

type cacheEntry struct {
	once sync.Once
	tile *Tile
	err  error
}

type TileServer struct {
	DataDirectory string
	mutex         sync.Mutex
	cache         map[string]*cacheEntry
}

func NewTileServer(DataDirectory string) *TileServer {
	return &TileServer{
		DataDirectory: DataDirectory,
		cache:         make(map[string]*cacheEntry),
	}
}

const scaleFactor = 10

func (ts *TileServer) Height(gridref string) (float64, error) {
	gr, err := osgridref.ParseOsGridRef(gridref)
	tile, err := ts.LoadTile(gr)
	if err != nil {
		return 0, err
	}

	offsetX := gr.Easting % 10000 / 50
	offsetY := gr.Northing % 10000 / 50

	return float64(tile[offsetY][offsetX] / scaleFactor), nil
}

func (ts *TileServer) MustHeight(gridref string) float64 {
	h, err := ts.Height(gridref)
	if err != nil {
		panic(err)
	}
	return h
}

func (ts *TileServer) LoadTile(gr osgridref.OsGridRef) (*Tile, error) {
	tileName := ts.tileZip(gr)

	ts.mutex.Lock()

	if ts.cache == nil {
		ts.cache = make(map[string]*cacheEntry)
	}
	entry, ok := ts.cache[tileName]
	if !ok {
		entry = &cacheEntry{}
		ts.cache[tileName] = entry
	}

	ts.mutex.Unlock()

	requestCount.Inc()
	entry.once.Do(func() {
		readCount.Inc()
		entry.tile, entry.err = getTile(gr, tileName)
		if entry.err != nil {
			badCount.Inc()
		}
	})

	return entry.tile, entry.err
}

func (ts *TileServer) ClearCache() {
	ts.mutex.Lock()
	ts.cache = make(map[string]*cacheEntry)
	ts.mutex.Unlock()
}

func (ts *TileServer) tileZip(gr osgridref.OsGridRef) string {
	tileName := strings.ToLower(gr.StringNCompact(2))
	dirName := strings.ToLower(tileName[:2])
	// TODO - this is crap
	fileName := strings.ToLower(tileName[:4]) + "_OST50GRID_20200303.zip"
	return path.Join(ts.DataDirectory, dirName, fileName)
}

func getTile(gridref osgridref.OsGridRef, zipFile string) (*Tile, error) {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".asc") {
			src, err := f.Open()
			if err != nil {
				return nil, err
			}

			tile, err := readTile(gridref, src)
			src.Close()

			return tile, err
		}
	}

	return nil, fmt.Errorf("no .asc file in %s", zipFile)
}

func readTile(gridref osgridref.OsGridRef, r io.Reader) (*Tile, error) {
	// Format is:
	// 	ncols 200
	// 	nrows 200
	// 	xllcorner 310000
	// 	yllcorner 500000
	// 	cellsize 50
	// 	200 rows, each with 200 columns.
	//		northern-most row is at the top.

	scanner := bufio.NewScanner(r)

	xllCorner := fmt.Sprintf("xllcorner %d", trunc10k(gridref.Easting))
	yllCorner := fmt.Sprintf("yllcorner %d", trunc10k(gridref.Northing))
	err := expect(scanner, "ncols 200", "nrows 200", xllCorner, yllCorner, "cellsize 50")
	if err != nil {
		return nil, err
	}

	var values Tile
	for i := 0; i < 200; i++ {
		ok := scanner.Scan()
		if !ok {
			return nil, fmt.Errorf("short file after %d records for %q", i, gridref.StringNCompact(2))
		}

		row := strings.Split(scanner.Text(), " ")
		if len(row) != 200 {
			return nil, fmt.Errorf("%d values in row %d for %q", len(row), i, gridref.StringNCompact(2))
		}

		valueRow := 199 - i // File has most northerly at the top; we want first row in returned Tile to be most southerly
		for x := 0; x < len(row); x++ {
			f, err := strconv.ParseFloat(row[x], 16)
			if err != nil {
				return nil, fmt.Errorf("row %d for %q: %s", i, gridref.StringNCompact(2), err)
			}
			values[valueRow][x] = int16(f * scaleFactor)
		}
	}

	if scanner.Scan() {
		return nil, fmt.Errorf("overlong file for %q: %s", gridref.StringNCompact(2), scanner.Text())
	}

	return &values, err
}

func expect(scanner *bufio.Scanner, lines ...string) error {
	for _, line := range lines {
		ok := scanner.Scan()
		if !ok {
			return fmt.Errorf("EOF before %q", line)
		}
		if scanner.Text() != line {
			return fmt.Errorf("expected %q but got %q", line, scanner.Text())
		}
	}
	return nil
}

func trunc10k(n int) int {
	return (n / 10000) * 10000
}
