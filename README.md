# BIT Inc. 2020 - Bolk Assignment

## Transics TX-TANGO Data Importer & Driver Style Analysis

The program will import data from Transics directly into the database.
It will checks what fields are already present in the database before importing.

The program can also generate reports aimed at drivers for a given period of time.

### Requirements

* Go
* R
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

Install the required dependencies (on a Ubuntu/Debian based distribution) by running `config/install.sh`.

### Usage

The efficient way to import or generate report periodically data is to use **CRON** (for instance every 4h for the importer and every wednesday for the analysis).  
This can as well be done manually using `tx2db` commands.

An exaustive list of available commands can be found by running `tx2db --help`.

#### Importer

Run the import manually
```tx2db import```

Options exist for this command, more information by running `tx2db import --help`

#### Report

Generate the report manually
```tx2db gen-report```

Generate a report from specific date
```tx2db gen-report --startTime 2020-02-22```

Options exist for this command, more information by running `tx2db gen-report --help`

### Architechture

* ```analysis``` contains the driver analysis. Graphs are built with R and the different metrics in SQL via Go. The template of the report is written in `.html`. The reports are then converted to a `.png` thanks to `phantomjs`.
* ```cmd``` are the commands accessible in `tx2db`
* ```config```  are configuration files: please read [config/README.md](config/README.md).
* ```txtango``` implements the TX-TANGO API

### More Info

More info about Transics TX-TANGO API:
* [TX-TANGO API](http://integratorsprod.transics.com/OperationOverview.aspx)
