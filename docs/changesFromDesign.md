


# Additional Features

4/16/2023: 
 - Created a message pipeline for ViewUser...
    - apiEndpoint: "/viewUser"
    - communication.ViewUser = { Username: string }
    - sendsBack: communication.PublicUserProfile = UserProfile

2/7/2023 
 - Included a server CLI. This is necessary and useful as we go forwards testing the server. We can use the CLI to get output for various statuses, such as the # of users logged in, the # of pages cached.
 - Instead of 1 logger function, there are two logger functions. An initial log occurs when a request is received, and it prints the API endpoint and the HTTP Method. The other logger is called last in the lifecycle, and it prints the User ID.


# Package communication

3/28/2023:
- Communication.Ban now bans based on username rather than ID. Username is sent, not ID (and user is subsequently looked up in DB to get their ID). Also has Ban.ban to determine whether an action is supposed to be an unban.

3/22/2023:
- Communication.Moderate, ReportID changed to *int64 to allow to nil values. (Comments can be moderated without associated report)
- Moderate.ViewModRecords was given field 'From', 'To', 'ByUser', and 'ForDomain'

3/21/2023 
- ViewLogs.ForUserID to ViewLogs.ForUser
- ViewLogs.ForDomain to ViewLogs.ForEndpoint
- AdminAccessLog.AtTime was added.

2/16/2023
- Change Comment.UserId from type 'string' to type 'int64'. (Typo in design doc. )
- Added FullPage server-client communication entity for when a client first requests comments for a page.  


2/12/2023

- Got rid of PasswordResetCode client-server communication entity and combined into one SetNewPass struct. Got rid of associated API endpoint. 
- Added ProfileUpdateResponse which is dispatched to the client when a change to their profile has been realized on the server.


2/10/2023

- Added Token Struct 

    - While dealing with logout, realized that we would not be able to use simple cookies with our responses. The reason is because of CORS security policies in the browser. If we allow CORS requests from anywhere, then we can't let our fetch requests pass cookies because it would be trivial for a 3rd party app to use stored cookies to access our users accounts. From MDN:
    > Note: Access-Control-Allow-Origin is prohibited from using a wildcard for requests with credentials: 'include'. In such cases, the exact origin must be provided; even if you are using a CORS unblocker extension, the requests will still fail.
    [See MDN Docs](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch)

- Added LogoutResponse struct

2/3/2023

- added methods for, e.g., wrapping responses


## DB Structs

2/12/2023

- The password reset code is now equal to the unique primary key ID of the reset code in the reset code database. 


2/2/2023

- got rid of "NewConnect" function in database/store.go; only using New now, with a boolean if you don't want to connect to the database


## Server

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