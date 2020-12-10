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
		easting     int
		northing    int
		expectedLat float64
		expectedLon float64
	}{
		{
			name:    "SJ 92395 52997",
			easting: 392395, northing: 352997,
			// From  http://www.movable-type.co.uk/scripts/latlong-os-gridref.html:
			// 	For OSGB36, expect 53.073851°N, 002.113526°W
			//	For WGS84,  expect 53.074149°N, 002.114964°W
			expectedLat: +53.073851,
			expectedLon: -2.113526,
		},
		{
			name:    "TG 51409 13177",
			easting: 651409, northing: 313177,
			expectedLat: +52.657568,
			expectedLon: 001.717908,
		},
		{
			name:    "OS Worked Example",
			easting: 651410, northing: 313177,
			expectedLat: 52.0 + (39.0 / 60) + (27.2531 / 3600),
			expectedLon: 1 + (43.0 / 60) + (4.5177 / 3600),
		},
		{
			name:    "Movable Type Example",
			easting: 544982, northing: 257869,
			expectedLat: 52.199558,
			expectedLon: 0.121654,
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := OsGridRef{
				Easting:  tt.easting,
				Northing: tt.northing,
			}
			lat, lon := o.toLatLon()
			fmt.Printf("%s: expected %f,%f got %f,%f\n", tt.name, tt.expectedLat, tt.expectedLon, lat, lon)
			assert.InEpsilon(t, tt.expectedLat, lat, 1.0/(50*60*60))
			assert.InEpsilon(t, tt.expectedLon, lon, 1.0/(50*60*60))
		})
	}
}

func TestParseOsGridRef(t *testing.T) {
	tests := []struct {
		s string
		want    OsGridRef
		wantErr bool
	}{
		{
			s:       "651409, 313177",
			want:    OsGridRef{Easting:  651409, Northing: 313177},
			wantErr: false,
		},
		{
			s:       "TG 51409 13177",
			want:    OsGridRef{Easting:  651409, Northing: 313177},
			wantErr: false,
		},
		{
			s:       "SU 0 0",
			want:    OsGridRef{Easting:  400000, Northing: 100000},
			wantErr: false,
		},
		{
			s:       "SE095255",
			want:    OsGridRef{Easting:  409500, Northing: 425500},
			wantErr: false,
		},
		{
			s:       "SE0849025580",
			want:    OsGridRef{Easting:  408490, Northing: 425580},
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