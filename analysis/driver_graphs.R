#Driving Style Analysis in R: Setup

#Install/Load required libraries
if (!require(dotenv)) install.packages("dotenv", dep = TRUE)
if (!require(odbc)) install.packages("odbc", dep = TRUE)
if (!require(tidyverse)) install.packages("tidyverse", dep = TRUE)
if (!require(leaflet)) install.packages("leaflet", dep = TRUE)
if (!require(mapview)) {
  install.packages("mapview", dep = TRUE)
  library(webshot)
  webshot::install_phantomjs()
}

#Load .env file
library(dotenv)
db_host <- Sys.getenv("DB_HOST")
db_name <- Sys.getenv("DB_NAME")
db_username <- Sys.getenv("DB_USERNAME")
db_password <- Sys.getenv("DB_PASSWORD")

#Connect to SQL database
library(DBI)
library(odbc)
conn <- dbConnect(odbc(),
                 Driver = "SQL Server",
                 Server = db_host,
                 Database = db_name,
                 UID = db_username,
                 PWD = db_password,
                 Port = 1433)

#Set working directory
setwd("analysis")

###############
#####GRAPH#####
###############

#load libraries
library(tidyverse)
library(leaflet)
library(mapview)

#build a map of visited places of a drivers
buildMap = function(conn, driverTransicsID, startTime, endTime) {
  #get all destinations of a given drivers
  destinations <- tbl(conn, "tours") %>%
    filter(driver_transics_id == driverTransicsID) %>%
    filter(start_time >= startTime && end_time <= endTime) %>%
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
      color = "#003580",
      opacity = 0.2
    )
  
  map_name <-  paste0("driver_", driverTransicsID, "_maps.png")
  mapshot(map, file = map_name)
}