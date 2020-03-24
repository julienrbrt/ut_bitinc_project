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
if (!require(gridExtra)) install.packages("gridExtra", dep = TRUE)
if (!require(grid)) install.packages("grid", dep = TRUE)

#Load .env file
library(dotenv)
db_host <- Sys.getenv("DB_HOST")
db_name <- Sys.getenv("DB_NAME")
db_username <- Sys.getenv("DB_USERNAME")
db_password <- Sys.getenv("DB_PASSWORD")

#Connect to SQL database
library(DBI)
library(odbc)
#This is commented as the connection will be done from Go
conn <- dbConnect(odbc(),
                 Driver = "SQL Server",
                 Server = db_host,
                 Database = db_name,
                 UID = db_username,
                 PWD = db_password,
                 Port = 1433)

#Set working directory
setwd("analysis/assets/graph")

###############
#####GRAPH#####
###############

#load libraries
library(tidyverse)
theme_set(theme_bw())
library(leaflet)
library(mapview)
library(gridExtra)
library(grid)

#parse dates
parseDate = function(datecolumn) {
  if (datecolumn == "0001-01-01 00:00:00.0000000 +00:00") {
    return(as.character(Sys.time()))
  } else {
    yEnd = str_sub(datecolumn, 1,-16)
    zEnd = str_sub(datecolumn, 29,-7)
    end_time = paste(yEnd, zEnd, sep = "")
    return(end_time)
  }
}

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
    labs(x = "", y = "Verhouding stationair draaien tov totale rijtijd (%)") +
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
    labs(x = "", y = "Gemiddelde sneilheid") +
    theme(axis.text.x = element_text(vjust = 0.5))
  
  #save it to file
  graph_name <-  paste0("driver_", driverTransicsID, "_high_speed.png")
  ggsave(graph_name)
}

#Build list of activities
buildActivityList = function(conn, driverTransicsID, startTime, endTime) {
  #select driver ids and tour ids from tours
  activityList = tbl(conn, "tours") %>%
    filter(start_time>= startTime && end_time <= endTime && driver_transics_id == driverTransicsID) %>%
    select(tour_id = id, driver_transics_id) %>%
    # join tours and activities to connect driver _ids to activities
    inner_join(tbl(conn,"truck_activity_reports") %>% 
               filter(start_time>= startTime && end_time <= endTime) %>%
               select(tour_id, activity, start_time, end_time), by = "tour_id") %>%
    collect() %>%
    #format end and start time
    mutate(start_time = sapply(start_time, parseDate)) %>%
    mutate(end_time = sapply(end_time, parseDate)) %>%
    #create duration column
    mutate(duration = difftime(end_time,start_time, units = "secs")) %>%
    filter(duration > 0) %>%
    group_by(activity) %>%
    summarize(duration = sum(as.numeric(duration)))

  #build grid
  data = as.data.frame(paste(round(activityList$duration / sum(activityList$duration) * 100, digits = 2), "%",sep = ""), row.names = activityList$activity)
  tt3 <- ttheme_minimal(core=list(bg_params = list(fill = blues9[4:1], col=NA), fg_params=list(fontface=3)),colhead=list(fg_params=list(col="#003580", fontface=4L)), rowhead=list(fg_params=list(col="#003580", fontface=3L)), base_size = 28)
  
  #save it to file
  graph_name <-  paste0("driver_", driverTransicsID, "_activity.png")
  png(graph_name)
  tableGrob(data, cols = "Total Duration", theme = tt3) %>%
    grid.arrange()
  dev.off()
}

###################
#####GENERTATE#####
###################

#get arguments
args <- commandArgs(trailingOnly = TRUE)
args <- as.Date(args)

#get list of report to generate
getReport = function(startTime, endTime) {
  tours <- tbl(conn, "tours") %>%
  filter(start_time >= startTime && end_time <= endTime) %>%
  select(id, driver_transics_id) %>%
  inner_join(
    tbl(conn, "driver_eco_monitor_reports") %>% 
      select(tour_id, distance),
    by = c("id" = "tour_id")
  ) %>%
  filter(distance > 0) %>%
  distinct(driver_transics_id) %>%
  collect()
  
  return(tours$driver_transics_id)
}

for (driverTransicsID in getReport(args[1], args[2])){
  buildMap(conn, driverTransicsID, args[1], args[2])
  buildIdling(conn, driverTransicsID, args[1] - 7, args[2])
  buildFuelConsumption(conn, driverTransicsID, args[1] - 7, args[2])
  buildHighSpeed(conn, driverTransicsID, args[1] - 7, args[2])
  buildActivityList(conn, driverTransicsID, args[1], args[2])
}