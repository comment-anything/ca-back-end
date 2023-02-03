


## DB Structs

2/2/2023

- got rid of "NewConnect" function in database/store.go; only using New now, with a boolean if you don't want to connect to the database


## Server struct

1/28/2023

- made http.Server a member, removed router member to httpServer member only (enables shutdown of server programatically)