# gosdata

## Elevation Data

> **TODO** flesh this out, add docs for the 2 commands and fix bodges. 

Decoding OS Elevation data in Go

Grid data: https://omseprd1stdstordownload.blob.core.windows.net/downloads/Terrain50/2020-07/allGB/ASCII%20Grid%20and%20GML/terr50_gagg_gb.zip
Contour data: https://omseprd1stdstordownload.blob.core.windows.net/downloads/Terrain50/2020-07/allGB/ESRI%20Shapes/terr50_cesh_gb.zip

User guide: https://www.ordnancesurvey.co.uk/documents/os-terrain-50-user-guide.pdf

* Map Header

Post spacing of 50m; vertical interval of 10m.
Tile size: 10km  x 10km

40000 per file, ~ 3000 files -> 120 million points. Each point could be held in a 16-bit signed.

--------
