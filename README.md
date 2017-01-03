# gosdata
Decoding OS Elevation data in Go

xllcorner goes from [0] 10000-650000
y                   [0] 10000-1220000

ncols, nrows 200
cellsize 50
missing -1


40000 per file, ~ 3000 files -> 120 million points.

--------

Great Gable is  321132,510350 or NY211103 (NY2113210350), 54.482161,-3.219135
Height: 899m

NW corner: 321000,511000 or NY210110
Height: ~550

NE corner: 322000, 511000, NY220110
Height: ~640

SW corner: 321000, 510000, NY210100
Height: ~~675 (on Great Napes so very approx_.

SE corner: 322000, 510000 or NY220100
Height: 490

-----------

File os-terrain-50/ny/ny21_OST50GRID_20160726.zip contains ny21.asc:

  ncols 200
  nrows 200
  xllcorner 320000
  yllcorner 510000
  cellsize 50


NY21.gml contains:

    <gml:lowerCorner>320000 510000</gml:lowerCorner>
    <gml:upperCorner>330000 520000</gml:upperCorner>
    origin <gml:pos>320000 510000</gml:pos>   - SW corner
    <gml:offsetVector>0 50</gml:offsetVector>
