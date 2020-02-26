# BIT Inc. 2020 - Bolk Assignment

## TX-TANGO Data Importer & Analysis

The program will import data from Transics directly into the database.
It will checks what fields are already present in the database before importing.

The program can also generate report from the data for a specific driver and truck for a given period of time.

### Requirements

* Go
* SQL Server

### Configuration

#### .env

The credentials of used services must be filled in the `.env` file. You can find an example of what information to fill-in in [.env.example](.env.example).

#### MSSQL

It is necessary to set a default schema in the database prior to use the program so as following:

```sql
ALTER USER [DATABASE_USER_NAME] WITH DEFAULT_SCHEMA=[dbo]
GO

```

### Usage

#### Importer

Run the import
```tx2db import```

The efficient way to import periodically data is to use **CRON**.

#### Analysis


### More Info

More info about Transics TX-TANGO API:
* [TX-TANGO API](http://integratorsprod.transics.com/OperationOverview.aspx)