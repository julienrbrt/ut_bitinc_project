#!/bin/bash

# Go dependencies
echo "Installing Go depencencies..."
apt install snapd
snap install Go

# R dependencies
echo "Installing R dependencies..."
curl https://packages.microsoft.com/keys/microsoft.asc | apt-key add -
### https://docs.microsoft.com/en-us/sql/connect/odbc/linux-mac/installing-the-microsoft-odbc-driver-for-sql-server?view=sql-server-ver15
curl https://packages.microsoft.com/config/ubuntu/18.04/prod.list > /etc/apt/sources.list.d/mssql-release.list
apt-get update
ACCEPT_EULA=Y apt-get install msodbcsql17 build-essential r-base r-base-dev unixodbc unixodbc-dev libcurl4-openssl-dev
ln -sr odbcinst.ini ~/.odbcinst.ini

# Redis dependencies
apt install redis