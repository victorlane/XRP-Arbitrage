## Arbitraging project
Get data from `n` amount of sources and see if there is an arbitrage possibility

## Sources
- Exchanges
    - Bitrue
    - Binance
    - etc
- Ledgers
    - XRPL
    - Stellar Dex
    - etc

## Todo
- [ ] Getting order book data from XRPL returns nothing but does not throw an error
- [ ] Figure out how to get actual XRPL price data<sup>1</sup>

<br>

1. Could potentially get the last actual trade for a pair and then compare to order book, could alternatively also just accept that there is spread and get the average of the order book

## Maybe

- Add more pairs from other Ledgers + DEX's<sup>2</sup>
- Create a website that updates the prices every `n` seconds or gives alerts if the spread is `> n%` 

<br>

2. DEX's usually have good arbitrage opportunities if the currency is listed on an exchange