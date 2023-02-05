


# Package communication

2/3/2023

- added methods for, e.g., wrapping responses


## DB Structs

2/2/2023

- got rid of "NewConnect" function in database/store.go; only using New now, with a boolean if you don't want to connect to the database


## Server struct

1/28/2023

- made http.Server a member, removed router member to httpServer member only (enables shutdown of server programatically)


2/4/2023

Ideas:

- Candidate for refactor: Move server.ReadsAuth to UserManager.ReadsAuth
- Candidate for refactor: Move server.EnsureController to UserManager.EnsureController
- Candidate for refactor: Move the logic GuestController.HandleRegister to UserManager.
- RegisterGuest and move all UserManager lifecycle commands to their own file "users.go" with methods such as UserManager.- LoginUser or UserManager.LogoutUser
- Cadidate for refactor: Server should have an instance of http.Server and cause that to run when appropriate, and stop it when appropriate. I wound up calling http.ListenAndServer just to get it working, but that should be fixed up. The Server should be able to stop. It would also be cool if there were a CLI for the server in another thread as well. 
- Candidate for Feature: Server CLI