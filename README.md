# realtime-securities
### Go package for retrieving and analyzing realtime stock and option data.

The Provider interface allows new data providers to be added.  Currently, the [Tradier API](https://documentation.tradier.com/brokerage-api) has been partially implemented.

Please see the cmd directory for sample executables.  Sample data is provided with this repo to run the showhistory command:
```
[realtime-securities] (master)$ cd cmd/showhistory/

[showhistory] (master)$ go build showhistory.go

[showhistory] (master)$ ./showhistory
Retrieving daily stock prices...

BAC:
              Open   Close    High     Low       Volume
2020-02-13   34.76   34.91   35.03   34.55   31,775,115
2020-02-14   34.88   34.85   34.96   34.70   26,447,711
2020-02-18   34.77   34.21   34.83   34.01   23,962,799

F:
              Open   Close    High     Low       Volume
2020-02-13    8.21    8.25    8.36    8.21   67,648,837
2020-02-14    8.27    8.10    8.27    8.08   46,359,668
2020-02-18    8.12    8.03    8.15    8.02   42,003,280

GOOG:
              Open   Close    High     Low       Volume
2020-02-13 1512.69 1514.66 1527.18 1504.60      929,730
2020-02-14 1515.60 1520.74 1520.74 1507.34    1,197,836
2020-02-18 1515.00 1520.38 1531.63 1512.59      718,331

MSFT:
              Open   Close    High     Low       Volume
2020-02-13  183.08  183.71  186.23  182.87   35,295,834
2020-02-14  183.25  185.35  185.41  182.65   23,149,516
2020-02-18  185.60  186.84  187.60  184.34   18,622,620

NFLX:
              Open   Close    High     Low       Volume
2020-02-13  376.96  381.40  385.37  376.51    4,485,383
2020-02-14  381.47  380.40  385.15  379.43    3,736,266
2020-02-18  379.30  388.67  388.98  379.19    3,612,165
```

Other commands require a Tradier authorization token to be placed in file resources/provider-auth/tradier to retrieve data from the provider.  The format is "Bearer G6hw9LRbs72mChWP81jqPZzx39mF" (not a valid token).

The dailyprices command will retrieve and persist daily price summaries for the given date range.  The quote command retrieves and displays realtime stock quotes.  The list of stocks can be found in resources/data/symbols.dat.

```
[realtime-securities] (master)$ cd cmd/dailyprices/

[dailyprices] (master)$ go build dailyprices.go

[dailyprices] (master)$ ./dailyprices 02/13/2020
Loading daily stock prices...
   1: BAC
      02/13/2020
      02/14/2020
      02/18/2020
   2: F
      02/13/2020
      02/14/2020
      02/18/2020
   3: GOOG
      02/13/2020
      02/14/2020
      02/18/2020
   4: MSFT
      02/13/2020
      02/14/2020
      02/18/2020
   5: NFLX
      02/13/2020
      02/14/2020
      02/18/2020

[dailyprices] (master)$ cd ../quotes/

[quotes] (master)$ go build quotes.go

[quotes] (master)$ ./quotes

Stock     Last   Change   %Change
=====     ====   ======   =======
    F     8.05    -0.05    -0.56%
  BAC    34.27    -0.59    -1.67%
 MSFT   187.41     2.06     1.11%
 NFLX   389.15     8.75     2.30%
 GOOG  1524.00     3.26     0.22%
```
