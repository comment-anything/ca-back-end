


# Additional Features

2/7/2023 
 - Included a server CLI. This is necessary and useful as we go forwards testing the server. We can use the CLI to get output for various statuses, such as the # of users logged in, the # of pages cached.
 - Instead of 1 logger function, there are two logger functions. An initial log occurs when a request is received, and it prints the API endpoint and the HTTP Method. The other logger is called last in the lifecycle, and it prints the User ID.


# Package communication


2/10/2023

- Added Token Struct 

    - While dealing with logout, realized that we would not be able to use simple cookies with our responses. The reason is because of CORS security policies in the browser. If we allow CORS requests from anywhere, then we can't let our fetch requests pass cookies because it would be trivial for a 3rd party app to use stored cookies to access our users accounts. From MDN:
    > Note: Access-Control-Allow-Origin is prohibited from using a wildcard for requests with credentials: 'include'. In such cases, the exact origin must be provided; even if you are using a CORS unblocker extension, the requests will still fail.
    [See MDN Docs](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch)

- Added LogoutResponse struct

2/3/2023

- added methods for, e.g., wrapping responses


## DB Structs

2/2/2023

- got rid of "NewConnect" function in database/store.go; only using New now, with a boolean if you don't want to connect to the database


## Server struct

2/10/2023

- changed server method postLogout to putLogout to reflect the HTTP method actually used

1/28/2023

- made http.Server a member, removed router member to httpServer member only (enables shutdown of server programatically)


2/4/2023

Ideas:

- Candidate for refactor: Change controller.SetCookie to something more descriptive since it is no longer setting a cookie.

- Candidate for refactor: Move server.ReadsAuth to UserManager.ReadsAuth
- Candidate for refactor: Move server.EnsureController to UserManager.EnsureController
- Candidate for refactor: Move the logic GuestController.HandleRegister to UserManager.
- RegisterGuest and move all UserManager lifecycle commands to their own file "users.go" with methods such as UserManager.- LoginUser or UserManager.LogoutUser
- Cadidate for refactor: Server should have an instance of http.Server and cause that to run when appropriate, and stop it when appropriate. I wound up calling http.ListenAndServer just to get it working, but that should be fixed up. The Server should be able to stop. It would also be cool if there were a CLI for the server in another thread as well. 
- Candidate for Feature: Server CLI