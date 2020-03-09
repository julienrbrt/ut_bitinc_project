#Driving Style Analysis in R: Setup

#Set working directory
setwd("./analysis")

#Load required libraries
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