#!/bin/bash

# General dependencies
sudo apt-get -y install build-essential curl

# Go dependencies
echo "Installing Go depencencies..."
sudo apt-get -y install snapd
sudo snap install go

# R dependencies
echo "Installing R dependencies..."
curl https://packages.microsoft.com/keys/microsoft.asc | sudo apt-key add -
### https://docs.microsoft.com/en-us/sql/connect/odbc/linux-mac/installing-the-microsoft-odbc-driver-for-sql-server?view=sql-server-ver15
sudo su -c "curl https://packages.microsoft.com/config/ubuntu/18.04/prod.list > /etc/apt/sources.list.d/mssql-release.list"
sudo apt-get update
ACCEPT_EULA=Y sudo apt-get -y install msodbcsql17 r-base r-base-dev unixodbc unixodbc-dev libcurl4-openssl-dev libssl-dev libudunits2-dev libfontconfig1-dev libcairo2-dev libgdal-dev pandoc
ln -sr odbcinst.ini ~/.odbcinst.ini
sudo chmod 777 -R /usr/local/lib/R/site-library