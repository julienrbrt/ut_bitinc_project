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
theme_set(theme_bw())
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
  
  graph_name <-  paste0("driver_", driverTransicsID, "_maps.png")
  mapshot(map, file = graph_name)
}

#Build total idling percentage barplot
buildIdling = function(conn, driverTransicsID, startTime, endTime) {
  idling <- tbl(conn, "tours") %>%
    filter(driver_transics_id == driverTransicsID) %>%
    filter(start_time >= startTime && end_time <= endTime) %>%
    select(id) %>%
    inner_join(
      tbl(conn, "driver_eco_monitor_reports") %>% 
        select(tour_id, start_time, duration_idling, duration_driving),
      by = c("id" = "tour_id")
    ) %>%
    collect()
  
  #convert the time to R date object (Warning, we are losing the actual time)
  idling$start_time <- as.Date(idling$start_time)
  #get week number
  idling$week_number <- paste("Week", strftime(idling$start_time, format="%V"))
  
  #calculating ratio idling time / driving
  idling <- idling %>% group_by(week_number) %>% summarize(idling = sum(duration_idling) / sum(duration_idling+duration_driving) * 100)

  #building histogram
  idling %>% ggplot(aes(x=week_number, y=idling)) +
    geom_bar(stat="identity", fill = "#003580", alpha = 0.8) +
    labs(x = "Week nummer", y = "Verhouding stationair draaien tov totale rijtijd (%)") +
    theme(axis.text.x = element_text(vjust = 0.5))
  
  #save it to file
  graph_name <-  paste0("driver_", driverTransicsID, "_idling.png")
  ggsave(graph_name)
}

#Build Fuel Consumption barplot
buildFuelConsumption = function(conn, driverTransicsID, startTime, endTime) {
  consumption <- tbl(conn, "tours") %>%
    filter(driver_transics_id == driverTransicsID) %>%
    filter(start_time >= startTime && end_time <= endTime) %>%
    select(id) %>%
    inner_join(
      tbl(conn, "driver_eco_monitor_reports") %>% 
        select(tour_id, fuel_consumption, start_time, distance),
      by = c("id" = "tour_id")
    ) %>%
    collect()
  
  #convert the time to R date object (Warning, we are losing the actual time)
  consumption$start_time <- as.Date(consumption$start_time)
  
  #summing the fuel consumption per day
  consumption <- consumption %>% group_by(start_time) %>% summarize(fuel_consumption = sum(distance) / sum(fuel_consumption))

  #building histogram
  consumption %>% ggplot(aes(x=start_time, y=fuel_consumption)) +
    geom_bar(stat="identity", fill = "#003580", alpha = 0.8) +
    scale_x_date(date_breaks = "2 days", date_labels = "%d %b %Y") +
    labs(x = "", y = "Verbruik (Km / L)") +
    theme(axis.text.x = element_text(angle = 75, vjust = 0.5))
  
  #save it to file
  graph_name <-  paste0("driver_", driverTransicsID, "_fuel_consumption.png")
  ggsave(graph_name)
}

#Build total high speed percentage barplot
buildHighSpeed = function(conn, driverTransicsID, startTime, endTime) {
  speed <- tbl(conn, "tours") %>%
    filter(driver_transics_id == driverTransicsID) %>%
    filter(start_time >= startTime && end_time <= endTime) %>%
    select(id) %>%
    inner_join(
      tbl(conn, "driver_eco_monitor_reports") %>% 
        select(tour_id, start_time, speed_average, distance),
      by = c("id" = "tour_id")
    ) %>%
    filter(distance > 0 && speed_average > 0) %>%
    collect()
  
  #convert the time to R date object (Warning, we are losing the actual time)
  speed$start_time <- as.Date(speed$start_time)
  #get week number
  speed$week_number <- paste("Week", strftime(speed$start_time, format="%V"))
  
  #calculating speed average per week
  speed <- speed %>% group_by(week_number) %>% summarize(speed_average = mean(speed_average))
  
  #building histogram
  speed %>% ggplot(aes(x=week_number, y=speed_average)) +
    geom_bar(stat="identity", fill = "#003580", alpha = 0.8) +
    labs(x = "Week nummer", y = "Gemiddelde sneilheid") +
    theme(axis.text.x = element_text(vjust = 0.5))
  
  #save it to file
  graph_name <-  paste0("driver_", driverTransicsID, "_high_speed.png")
  ggsave(graph_name)
}