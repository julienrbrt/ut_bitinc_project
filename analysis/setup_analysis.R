#Setup Driving Style Analysis in R

#Path
rLib <- "./analysis/rlib"
result <- "./analysis/results"
#Add folder to library
.libPaths(c(rLib, .libPaths()))

#Load required libraries
if (!require(tidyverse)) {
    install.packages("readr", dep = TRUE, INSTALL_opts = c('--no-lock'))
    install.packages("tidyverse", dep = TRUE, INSTALL_opts = c('--no-lock'))
}
library(tidyverse)

if (!require(dotenv)) install.packages("dotenv", dep = TRUE)
if (!require(odbc)) install.packages("odbc", dep = TRUE)

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