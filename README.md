# gosdata

# Grid Refs

(0,0) is SW of Scilly Isles at SV000000.

AA is near Iceland and ZZ in Northern Italy, but really there only a few useful squares for mainland GB.

First letter defines 500km x 500km square:

Origins (SW)

   H (0, 1000000); 
   N (0,  500000);  O (500000, 500000)
   S (0,       0);  T (500000, 0) 

Square O is practically empty.

Second letter is

   A B C D E
   F G H J K
   L M N O P
   Q R S T U
   V W X Y Z 

## Elevation Data

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

For the 1000m grid square containing GG:

NW corner: 321000,511000 or NY210110
Height: ~550

NE corner: 322000, 511000, NY220110
Height: ~640

SW corner: 321000, 510000, NY210100
Height: ~~675 (on Great Napes so very approx_.

SE corner: 322000, 510000 or NY220100
Height: 490


         NY2111                NY2211
         
         NY2110                NY2210


-----------

File os-terrain-50/ny/ny21_OST50GRID_20160726.zip contains ny21.asc:

    ncols 200
    nrows 200
    xllcorner 320000
    yllcorner 510000
    cellsize 50

* Each line has 200 (ncols) values per line, and there are 200 lines.
* Each file is a 10km² square (10000). Therefore 10x10 of the 1km² above.


NY21.gml contains:

    <gml:lowerCorner>320000 510000</gml:lowerCorner>
    <gml:upperCorner>330000 520000</gml:upperCorner>
    origin <gml:pos>320000 510000</gml:pos>   - SW corner
    <gml:offsetVector>0 50</gml:offsetVector>

Height at origin (320000,510000) is ~415m, which matches with first column of last line. SE corner matches with 474 - approx right.

