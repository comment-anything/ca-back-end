-- bans 3/29


// TODO 

- ~~unban user ~~
- ~~add ban record for domain ban~~
- ~~add ban record for global ban~~
- ~~global bans ~~
- ~~prevent global banned user from logging in~~
- ~~log global banned user out immediately~~
- ~~prevent domain banned user from commenting~~
- ~~push ban lists to user profile~~

--- Misc

- limit max size of communication entity
- make usernames not case sensitive?
- encrypt passwords

--- Comments

2/19/2023
- should disallow empty comments
- should disallow comments over some max limit
- should regex validate comments to filter bad words and such


2/18/23 
Will know more after testing, but....

should Server-Client communication comments have info regarding the url they are for? 

should Page.NewComment push a communication.Comment to all user's on the page? Will that cause the other user's on that page to possibly have wrong responses when they subsequently get comments for a different page? 
 Solutions: A map for each user with pending comments...
            Keeping domain and path data with Server.Comment structures

The latter is probably the right solution. Then the front end can choose to display or not display the comments in the server response, if they match the page the user is actually wanting to view. 

Actually, may not be a problem. A NewPage struct will always be later in the ServerResponse array and will cause a clearing of the current page on the front end....
Maybe I'm overthinking it!



--- UserProfile

- in both register and login, the UserProfile portion of a LoginResponse is generated line by line, but is incomplete. Additional queries need to be performed to determine, for example, which domains a user can moderate. A func which converts a User to a UserProfile response could be used in both situations. 

2/18/23 : Do this with Store 


--- Validators

- these will be fun to implement

password entropy check : https://golangexample.com/validate-the-strength-of-a-password-in-go/

should also check a rainbow table of common passwords

Can move to util module


---- auth 

- gettoken function in server/auth ; we should make the expiration time changable in the .env file perhaps



-- decent tuts on digitalocean, always: https://www.digitalocean.com/community/tutorials/how-to-make-an-http-server-in-go