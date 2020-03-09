#Driving Style Analysis in R: Maps of visited places for drivers

if (!require(tidyverse))
  install.packages("tidyverse", dep = TRUE)
if (!require(leaflet))
  install.packages("leaflet", dep = TRUE)
if (!require(mapview)) {
  install.packages("mapview", dep = TRUE)
  library(webshot)
  webshot::install_phantomjs()
}

#load packages
library(tidyverse)
library(leaflet)
library(mapview)

buildMap = function(conn, driverTransicsID) {
  #get all destinations of a given drivers
  destinations <- tbl(conn, "tours") %>%
    filter(driver_transics_id == driverTransicsID) %>%
    select(id, destination_latitude, destination_longitude) %>%
    collect()
  
  map <- leaflet(data = destinations) %>%
    addTiles() %>%  # Add default OpenStreetMap map tiles
    addMarkers(lat =  ~ destination_latitude,
               lng =  ~ destination_longitude) %>%
    addPolylines(
      lat = ~ destination_latitude,
      lng = ~ destination_longitude,
      group = ~ id,
      color = "#0984e3"
    )
  
  map_name <-  paste0("driver_", driverTransicsID, "_maps.png")
  mapshot(map, file = map_name)
}