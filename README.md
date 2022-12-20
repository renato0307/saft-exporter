# SAFT exporter for Invoicexpress

## How it works

![how it works](how.png)

## Flags

* `-year`, `-month`: defines the month of the SAFT file; if not set fetches last month.
* `-destination`: where to copy the file
* `-month-year-pattern`: used to append the year and month to the destination folder; for example, the default value is `%s-%s` which will make the file being copied to `/home/pi/accounting/2022-10`.

## Environment variables

It uses the following environment variables to access the Invoicexpress API:

* `INVOICE_ACCOUNT_NAME`
* `INVOICE_API_KEY`

## Build for Raspberry PI

```shell
env GOOS=linux GOARCH=arm GOARM=7 go build
```

## Copy to Raspberry PI

```shell
scp saft-exporter pi@mypi:~
```

## Adding to crontab in Raspberry PI

Add the following variables to `/home/pi/.profile`:

```bash
export INVOICE_ACCOUNT_NAME="your-invoicexpress-account-anme"
export INVOICE_API_KEY="your-invoicexpress-api-key"
```

Create the `export-saft.sh` file with the following contents:

```bash
#!/bin/bash
./saft-exporter --destination /home/pi/accounting
```

Run `crontab -e` and add the following:

```
# Runs at 00:00 on day-of-month 2.
0 0 2 * * BASH_ENV=/home/pi/.profile /home/pi/export-saft.sh >> /home/pi/export-saft.log 2>&1
```