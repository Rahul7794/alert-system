# Alert System

### Introduction

An "alerting" service which will consume a file of currency conversion rates and produce alerts for a number of situations.

- When the spot rate for a currency pair changes by more than 10% from the 5-minute average for that currency pair.
- When the spot rate has been rising/falling for 15 minutes. This alert should be
  throttled to only output once per minute and should report the length of time of the rise/fall in seconds.
  
  
### Steps to build

`make build`

### Steps to run

```
./alert-system alert -i {input_path} -o {outpath}
```
The application expects input path where the currency conversion rates is available and output path where it wants to write the alerts.

  
  

