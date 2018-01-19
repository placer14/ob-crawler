# OpenBazaar Crawler

## Installation

`go get github.com/placer14/ob-crawler` or `docker pull placer14/ob-crawler`

## Usage

```
./ob-crawler -api-host 0.0.0.0 -api-port 4002 -api-timeout 30
-max-visits 100 -n 50
```

#### `-auth-cookie` (__string__) (__required__)

The crawler must be pointed at an OpenBazaar API using the
Authentication cookie it generated upon start. (More information at
https://github.com/OpenBazaar/openbazaar-go/blob/master/docs/security.md#authentication-cookie)

This option allows the crawler to include it on requests made to the API.

#### `-api-timeout` (__int__) (__optional__, __default: 60__)

This option tells the crawler how many seconds to wait before disconnecting and trying the next operation.

#### `-max-visits` (__int__) (__optional__, __default: 0__)

This option tells the crawler how many nodes to retrieve information about before stopping.

#### `-n` (__int__) (__optional__, __default: 10__)

This options tells the crawler how many workers to start.

#### `-api-host` (__string__) (__optional__, __default: api__)
#### `-api-port` (__int__) (__optional__, __default: 4002__)

## Purpose

This is a small toy golang program that I used to explore OpenBazaar's internal API.

## The Challenge

Using Go write an OpenBazaar network crawler that logs the total number of listings on the network. The crawler may connect to a running openbazaar-go instance and use the API.

You can find the API code here: https://github.com/OpenBazaar/openbazaar-go/tree/master/api

Some endpoints of interest may be:
- /ob/peers
- /ob/closestpeers
- /ob/profile
- /ob/listings

## Released under the MIT License

Copyright 2018 Mike Greenberg

Permission is hereby granted, free of charge, to any person obtaining a
copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be included
in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS
OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
