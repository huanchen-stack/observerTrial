# observerTrial

FIRST, CHANGE THE LIMIT OF #FD: `ulimit -n 1024000`

`make clean`

`make speedGO`

`./speedGO`

Notes:
 - TLS connections and GET requests are implemneted in the common package
 - Redirects are not enabled for now due to an incomplete version of http.Client cache
 - input.csv is already sanitized using ./sanitize/sanitize.go, meaning IPs that do not respond are filtered
