#!/bin/bash

# Dependencies
curl https://packages.microsoft.com/keys/microsoft.asc | sudo apt-key add -
### https://docs.microsoft.com/en-us/sql/connect/odbc/linux-mac/installing-the-microsoft-odbc-driver-for-sql-server?view=sql-server-ver15
sudo su -c "curl https://packages.microsoft.com/config/ubuntu/18.04/prod.list > /etc/apt/sources.list.d/mssql-release.list"
sudo apt-get update
ACCEPT_EULA=Y sudo apt-get -y install build-essential curl snapd msodbcsql17 r-base r-base-dev unixodbc unixodbc-dev libcurl4-openssl-dev libssl-dev libudunits2-dev libfontconfig1-dev libcairo2-dev libgdal-dev pandoc
sudo snap install go
ln -sr odbcinst.ini ~/.odbcinst.ini
sudo chmod 777 -R /usr/local/lib/R/site-library