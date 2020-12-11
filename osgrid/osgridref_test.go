package osgrid

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOsGridRef_toLatLon(t *testing.T) {
	type fields struct {
		Easting  int
		Northing int
	}
	tests := []struct {
		name        string
		gridRef     string
		expectedLat float64
		expectedLon float64
	}{
		{
			name: "SJ 92395 52997",
			// From  http://www.movable-type.co.uk/scripts/latlong-os-gridref.html:
			// 	For OSGB36, expect 53.073851°N, 002.113526°W
			//	For WGS84,  expect 53.074149°N, 002.114964°W
			expectedLat: +53.074149,
			expectedLon: -2.114964,
		},
		{
			name:        "TG 51409 13177",
			expectedLat: +52.657977,
			expectedLon: 1.716020,
		},
		{
			name:        "Movable Type Example (TL4498257869)",
			gridRef:     "TL4498257869",
			expectedLat: 52.199992,
			expectedLon: 0.119989,
		},
		{
			name:        "Cardiff (ST1784076329)",
			gridRef:     "ST1784076329",
			expectedLat: 51.479928,
			expectedLon: -3.184500,
		},
		{
			name:        "Aberdeen (NJ9439206608)",
			gridRef:     "NJ9439206608",
			expectedLat: 57.150318,
			expectedLon: -2.094323,
		},
		{
			name:        "Newlyn (SW4676028548)",
			gridRef:     "SW4676028548",
			expectedLat: 50.102910,
			expectedLon: -5.542751,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gridRef := tt.gridRef
			if gridRef == "" {
				gridRef = tt.name
			}
			o, err := ParseOsGridRef(gridRef)
			assert.NoError(t, err)
			lat, lon := o.ToLatLon()
			lat1, lon1, err := OttoGridToLatLon(gridRef)
			assert.NoError(t, err)
			fmt.Printf("%s: expected %f,%f got %f,%f (JS: %f,%f)\n", tt.name, tt.expectedLat, tt.expectedLon, lat, lon, lat1, lon1)
			assert.InDelta(t, tt.expectedLat, lat, 0.00005)
			assert.InDelta(t, tt.expectedLon, lon, 0.00005)
		})
	}
}

func TestParseOsGridRef(t *testing.T) {
	tests := []struct {
		s       string
		want    OsGridRef
		wantErr bool
	}{
		{
			s:       "651409, 313177",
			want:    OsGridRef{Easting: 651409, Northing: 313177},
			wantErr: false,
		},
		{
			s:       "TG 51409 13177",
			want:    OsGridRef{Easting: 651409, Northing: 313177},
			wantErr: false,
		},
		{
			s:       "SU 0 0",
			want:    OsGridRef{Easting: 400000, Northing: 100000},
			wantErr: false,
		},
		{
			s:       "SE095255",
			want:    OsGridRef{Easting: 409500, Northing: 425500},
			wantErr: false,
		},
		{
			s:       "SE0849025580",
			want:    OsGridRef{Easting: 408490, Northing: 425580},
			wantErr: false,
		},
		{
			s:       "SI095255",
			wantErr: true,
		},
		{
			s:       "ZZ095255",
			wantErr: true,
		},
		{
			s:       "S095255",
			wantErr: true,
		},
		{
			s:       "SJ95255",
			wantErr: true,
		},
		{
			s:       "SJ95X255",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			got, err := ParseOsGridRef(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOsGridRef() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseOsGridRef() got = %v, want %v", got, tt.want)
			}
		})
	}
}
