```
# Get basic whoami info
echo "" | nc localhost 8080

# Get bench response
echo "/bench" | nc localhost 8080
# Returns: "1"

# Get data of specific size
echo "/data?size=10" | nc localhost 8080
# Returns: "|ABCDEFGHI|"

echo "/data?size=5" | nc localhost 8080
# Returns: "|ABCD|"

echo "/data?size=0" | nc localhost 8080
# Returns: ""
```
