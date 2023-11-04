local http = require("http")

response, error_message = http.request("GET", "http://example.com", {
    query="page=1",
    timeout="30s",
    headers={
        Accept="*/*"
    }
})

write(response.body)
write(error_message)
