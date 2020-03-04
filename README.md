# BIT Inc. 2020 - Bolk Assignment

## TX-TANGO Data Importer & Driver Style Analysis

The program will import data from Transics directly into the database.
It will checks what fields are already present in the database before importing.

The program can also generate report from the data for a specific driver and truck for a given period of time.

### Requirements

* Go
* R
* Redis
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

#### Dependencies

Install all the required dependencies (on a Ubuntu/Debian based distribution) by running `config/install.sh`.

### Usage

The efficient way to import or generate report periodically data is to use **CRON**.  
This can as well be done manually using `tx2db` commands.

#### Importer

Run the import manually
```tx2db import```

#### Analysis

`tx2db` will call that analysis script to start the generation of the reports.
The analysis is performed in R.

Generate the report manually
```tx2db gen-report```

### Architechture

* ```analysis``` contains the analysis performed in R and the results of that analysis
* ```cmd``` are the commands accessible in `tx2db`
* ```config``` are configuration files used to setup the project (installation, drivers...). The files might require modification depending on which system `tx2db` is deployed.
* ```template``` contains logic building reports given templates
* ```test``` contains test of the program
* ```txtango``` implements the TX-TANGO API


### More Info

More info about Transics TX-TANGO API:
* [TX-TANGO API](http://integratorsprod.transics.com/OperationOverview.aspx)