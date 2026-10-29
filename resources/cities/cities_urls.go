package cities

// CitiesPath is the path template; the client fills {api_version}. It sits under /app, not /admin.
const CitiesPath = "/api/{api_version}/app/cities"

// CitiesKey is the key the listing's rows arrive under.
const CitiesKey = "cities"

// CitiesPageSize is the page size when limit is absent, and the largest limit accepted: a larger one is 400.
const CitiesPageSize = 100
