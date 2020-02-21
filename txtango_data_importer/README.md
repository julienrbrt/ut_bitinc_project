# TX-TANGO Data Importer

The program will import data from Transics directly into the database.
It will checks what fields are already present in the database before importing.

The efficient way to import periodically data is to use **CRON**.

## Requirements

* Go
* SQL Server
* Redis

## Configuration

### .env

The credentials of used services must be filled in the `.env` file. You can find an example of what information to fill-in in [.env.example](.env.example).

### MSSQL

It is necessary to set a default schema in the database prior to use the program so as following:

```sql
ALTER USER [DATABASE_USER_NAME] WITH DEFAULT_SCHEMA=[dbo]
GO
```

## More Info

More info about Transics TX-TANGO API:
* [TX-TANGO API](http://integratorsprod.transics.com/OperationOverview.aspx)