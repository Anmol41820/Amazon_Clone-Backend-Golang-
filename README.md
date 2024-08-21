#API Documentation for Amazon Clone Backend (GoLang)

General Information
Base URL: ‘/’
Authentication: Some routes require authentication via a middleware.
Content-Type: All requests and responses use application/json unless otherwise specified.

1. Home
Check Prime Membership
Endpoint: /customerId={id}
Method: GET
Description: Checks if a customer is a Prime member.
If the Customer’s Prime membership expires it make “isPrime” = false
Else “isPrime” remains true
Authentication: None
Parameters:
id (path): The customer's ID.
Response:
200 OK: Returns the customer's Prime membership status.
Example:
Prime Membership : False
{
   "MatchedCount": 1,
   "ModifiedCount": 0,
   "UpsertedCount": 0,
   "UpsertedID": null
}




2. User Management
Get All Users
Endpoint: /users
Method: GET
Description: Retrieves a list of all users.
Authentication: None
Response:
200 OK: Returns a list of users.
Example:


[
   {
       "_id": "66a8be438b66d10a2315cf2e",
       "firstName": "Dhyey",
       "lastName": "Mistri",
       "email": "dhyey@gmail.com",
       "password": "8a5d92205be0c2b95df533290bc82237b452dc93f81f7b6c87ce9d332ee768c2daa0db5477",
       "mobileNumber": "0987654321",
       "dateOfBirth": "18/12/2001",
       "role": "seller",
       "isPrime": false
   },
   {
       "_id": "66ab331d97d120b6a67da9f3",
       "firstName": "Shiva",
       "lastName": "Gupta",
       "email": "shiva@gmail.com",
       "password": "76eedce0787f23a4924e897cdf7e6e20254b04f61a5f2a71875a5448477fc4586188f2394a",
       "mobileNumber": "8982511209",
       "dateOfBirth": "18/12/2000",
       "role": "customer",
       "isPrime": false
   },
   {
       "_id": "66b4537d0208bbd82ee589aa",
       "firstName": "Ram",
       "lastName": "Gupta",
       "email": "ram@gmail.com",
       "password": "12a51cb9090c1c11dd35d8dc78a856cf6672df13dbbef16bd788ba0a6750df8cf62dc5da",
       "mobileNumber": "8982511208",
       "dateOfBirth": "18/12/2000",
       "role": "customer",
       "isPrime": false
   }
]


Get Single User
Endpoint: /user/userId={id}
Method: GET
Description: Retrieves details of a specific user by their ID.
Authentication: None
Parameters:
id (path): The user's ID.
Response:
200 OK: Returns the user's details.
Example:
{
   "_id": "66b4537d0208bbd82ee589aa",
   "firstName": "Ram",
   "lastName": "Gupta",
   "email": "ram@gmail.com",
   "password": "12a51cb9090c1c11dd35d8dc78a856cf6672df13dbbef16bd788ba0a6750df8cf62dc5da",
   "mobileNumber": "8982511208",
   "dateOfBirth": "18/12/2000",
   "role": "customer",
   "isPrime": false
}


Create a User
Endpoint: /addUser
Method: POST
Description: Creates a new user.
Get the user detail
Find all the exiting users and verify weather the current user’s email or mobile number is used before or not
Validate the entered password as length ≥8 and should contain at least one super character
Validate the entered mobile number as length == 10
Validate the date of birth as it should be less than the present date
Validate the email id
Encrypt the password
Insert the user in MongoDB
If role == customer
Created a new Cart for new customer
Created a new Wishlist for new customer
Created a new Search history for new customer
Created a new Product recommendation for new customer
Created a new Recently view product for new customer
Else if role ==seller
Creating a new report for the new seller
Authentication: None
Request Body:
{
        "firstName": "Shiva",
        "lastName": "Gupta",
        "email": "shiva@gmail.com",
        "password": "Shiva@123",
        "mobileNumber": "999999999",
        "dateOfBirth": "18/12/2000",
        "role": "customer",
        "isPrime": false
}

Response:
201 Created: Returns the created user's ID.
Example:
{
   "InsertedID": "66c222ed56fe9801831cd73c"
}
{
   "InsertedID": "66c222ed56fe9801831cd73e"
}
{
   "InsertedID": "66c222ed56fe9801831cd740"
}
{
   "InsertedID": "66c222f856fe9801831cd742"
}
{
   "InsertedID": "66c222f956fe9801831cd744"
}
{
   "InsertedID": "66c222f956fe9801831cd746"
}
{
   "InsertedID": "66c222ed56fe9801831cd73c"
}



Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004",
       "role": "customer"
}
"Invalid Email!!"
409 Conflict
Pass
2
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "somgmail.com",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004",
       "role": "customer"
}


"Invalid Email!!"
409 Conflict
Pass
3
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004",
       "role": "customer"
}
"Invalid Email!!"
409 Conflict
Pass
4
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004",
       "role": "customer"
}
"Password should be more than 7 characters!!"


409 Conflict
Pass
5
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som12345",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004",
       "role": "customer"
}
"Password should contain atleast one super characters!!"


409 Conflict
Pass
6
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som@1234",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004",
       "role": "customer"
}
"Your Mobile number must be of 10 digits!!"


409 Conflict
Pass
7
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som@1234",
       "mobileNumber": "6666666666",
       "dateOfBirth": "18/12/2025",
       "role": "customer"
}
"Invalid DOB: Date is in the future!"


409 Conflict
Pass


Update a User
Endpoint: /updateUser/userId={id}
Method: PUT
Authentication: Required
Description: Updates details of a specific user by their ID.
Role can not be changed while updating, role should be fixed by frontend
Parameters:
id (path): The user's ID.
Request Body:
{
        "firstName": "Yash",
        "lastName": "Yadav",
        "email": "yash@gmail.com",
        "password": "Yash@123",
        "mobileNumber": "9876543210",
        "dateOfBirth": "18/12/2003",
        "isPrime": false
}

Response:
200 OK: Returns the updated user's details.
Example:
{
   "_id": "66c222ed56fe9801831cd73c",
   "firstName": "Yash",
   "lastName": "Yadav",
   "email": "yash@gmail.com",
   "password": "2157c1264a8ca7280fac269bc7d3fff2b9a4793d206d4507e8ec84c723ddd55999338f94",
   "mobileNumber": "9876543210",
   "dateOfBirth": "08/12/2003",
   "role": "customer",
   "isPrime": false
}

Unit Testing


S.No.
Request
Response
Status
Pass/Fail
1
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004",
       "role": "customer"
}
"Invalid Email!!"
409 Conflict
Pass
2
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "somgmail.com",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004",
       "role": "customer"
}


"Invalid Email!!"
409 Conflict
Pass
3
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004",
       "role": "customer"
}
"Invalid Email!!"
409 Conflict
Pass
4
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004",
       "role": "customer"
}
"Password should be more than 7 characters!!"


409 Conflict
Pass
5
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som12345",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004",
       "role": "customer"
}
"Password should contain atleast one super characters!!"


409 Conflict
Pass
6
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som@1234",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004",
       "role": "customer"
}
"Your Mobile number must be of 10 digits!!"


409 Conflict
Pass
7
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som@1234",
       "mobileNumber": "6666666666",
       "dateOfBirth": "18/12/2025",
       "role": "customer"
}
"Invalid DOB: Date is in the future!"


409 Conflict
Pass


Delete a User
Endpoint: /deleteUser/userId={id}
Method: DELETE
Authentication: Required
Description: Deletes a specific user by their ID.
Delete the user from user API
If role is customer delete the user from customer API as well
Same for seller
Parameters:
id (path): The user's ID.
Response:
200 OK: Confirms deletion.
Example:
{
   "DeletedID": "66c222ed56fe9801831cd73c"
}



3. Registration and Authentication
Register a Customer
Endpoint: /registerCustomer
Method: POST
Description: Registers a new customer.
Initise wallet with zero money.
Role as a customer
Empty Address, add addresses later.
Authentication: None
Request Body:
{
        "firstName": "Shiva",
        "lastName": "Gupta",
        "email": "shiva@gmail.com",
        "password": "Shiva@123",
        "mobileNumber": "999999999",
        "dateOfBirth": "18/12/2000"
}


Response:
201 Created: Returns the registered customer's ID.
Example:


{
   "InsertedID": "66c22829818ea8e8db645d55"
}
{
   "InsertedID": "66c22829818ea8e8db645d55"
}
{
   "InsertedID": "66c2282e818ea8e8db645d58"
}
{
   "InsertedID": "66c2282f818ea8e8db645d5a"
}
{
   "InsertedID": "66c22835818ea8e8db645d5c"
}
{
   "InsertedID": "66c22836818ea8e8db645d5e"
}
{
   "InsertedID": "66c22836818ea8e8db645d60"
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Invalid Email!!"
409 Conflict
Pass
2
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "somgmail.com",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}


"Invalid Email!!"
409 Conflict
Pass
3
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Invalid Email!!"
409 Conflict
Pass
4
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Password should be more than 7 characters!!"


409 Conflict
Pass
5
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som12345",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Password should contain atleast one super characters!!"


409 Conflict
Pass
6
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som@1234",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Your Mobile number must be of 10 digits!!"


409 Conflict
Pass
7
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som@1234",
       "mobileNumber": "6666666666",
       "dateOfBirth": "18/12/2025"
}
"Invalid DOB: Date is in the future!"


409 Conflict
Pass


Register a Seller
Endpoint: /registerSeller
Method: POST
Description: Registers a new seller.
Role as a seller
Also create a new report for the seller’s products
Authentication: None
Request Body:


{

        "firstName": "Shiva",
        "lastName": "Gupta",
        "email": "shiva@gmail.com",
        "password": "Shiva@123",
        "mobileNumber": "999999999",
        "dateOfBirth": "18/12/2000"
}


Response:
201 Created: Returns the registered seller's ID.
Example:
{
   "InsertedID": "66c22878818ea8e8db645d62"
}
{
   "InsertedID": "66c22878818ea8e8db645d62"
}
{
   "InsertedID": "66c22879818ea8e8db645d65"
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Invalid Email!!"
409 Conflict
Pass
2
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "somgmail.com",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}


"Invalid Email!!"
409 Conflict
Pass
3
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Invalid Email!!"
409 Conflict
Pass
4
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Password should be more than 7 characters!!"


409 Conflict
Pass
5
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som12345",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Password should contain atleast one super characters!!"


409 Conflict
Pass
6
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som@1234",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Your Mobile number must be of 10 digits!!"


409 Conflict
Pass
7
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som@1234",
       "mobileNumber": "6666666666",
       "dateOfBirth": "18/12/2025"
}
"Invalid DOB: Date is in the future!"


409 Conflict
Pass


Login
Endpoint: /login
Method: POST
Description: Authenticates a user and provides a session token.
Get the entered email, password
Decrypt the stored password for that entered email
Check the stored password from the entered password
Set the cookie for the token generated by the user email and id
Authentication: None
Request Body:
{
        "email": "yash@gmail.com",
        "password": "Yash@123"
}


Response:
200 OK: Returns a session token.
Example:
"Welcome!!"
Yash, You have login as customer!!


Name
Value
Domain
Path
Expires
HttpOnly
Secure
token
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Inlhc2hAZ21haWwuY29tIiwiZXhwIjoxNzI0MDg1NjI1LCJpZCI6IjY2YzIyMmVkNTZmZTk4MDE4MzFjZDczYyJ9.kprzhrbMjLHXrjS3D0e8ys84O8GPzj_CdYUICM4D254
localhost
/
Mon, 19 Aug 2024 16:40:25 GMT
true
false

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
       "email": "som",
       "password": "Som@1234"
}
"Invalid Email!!"


409 Conflict
Pass
2
{
       "email": "somgmail.com",
       "password": "Som@1234"
}
"Invalid Email!!"
409 Conflict
Pass
3
{
       "email": "som@gmail",
       "password": "Som@1234"
}
"Invalid Email!!"
409 Conflict
Pass
4
{
       "email": "som@gmail",
       "password": "Som@123"
}
"Invalid User or Wrong password!!"


401 Unauthorized
Pass


Logout
Endpoint: /logout
Method: POST
Description: Logs out the user, invalidating the session token.
Delete the stored token from the cookie
Authentication: None
Response:
200 OK: Confirms logout.
Example:
Logged out successfully

Forgot Password
Endpoint: /forgotPassword
Method: POST
Description: Initiates the password recovery process for a user.
Check weather the email is available in the database or not
Validate the new password
Update the user with new password
Authentication: None
Request Body:

{
        "email": "anmol@gmail.com",
        "newPassword": "Anmol@123"
}


Response:
200 OK: Confirms the initiation of the password recovery process.
Example:
Password changed successfully!!

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
       "email": "som",
       "newPassword": "Som@1234"
}
"Invalid Email!!"


409 Conflict
Pass
2
{
       "email": "somgmail.com",
       "newPassword": "Som@1234"
}
"Invalid Email!!"
409 Conflict
Pass
3
{
       "email": "som@gmail",
       "newPassword": "Som@1234"
}
"Invalid Email!!"
409 Conflict
Pass
4
{
       "email": "sommm@gmail",
       "newPassword": "Som@1234"
}
"User not registered, Please Register first!!"


502 Bad Gateway
Pass
5
{
       "email": "som@gmail",
       "newPassword": "Som@123"
}
"Invalid User or Wrong password!!"


401 Unauthorized
Pass


Change Password
Endpoint: /changePassword
Method: POST
Description: Changes the user's password.
check weather the email is available in the database or not
check the present password is correct or not
Validate the new password
Update the user with new password
Authentication: None
Request Body:
{
        "email": "anmol@gmail.com",
        "password": "Anmol@123",
        "newPassword": "Anmol@1234"
}


Response:
200 OK: Confirms the password change.
Example:
Password reset successfully!!

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
       "email": "som",
       "password": "Som@1234",
       "newPassword": "Som@12345"
}
"Invalid Email!!"


409 Conflict
Pass
2
{
       "email": "somgmail.com",
       "password": "Som@1234",
       "newPassword": "Som@12345"
}
"Invalid Email!!"
409 Conflict
Pass
3
{
       "email": "som@gmail",
       "password": "Som@1234",
       "newPassword": "Som@12345"
}
"Invalid Email!!"
409 Conflict
Pass
4
{
       "email": "sommm@gmail",
       "password": "Som@1234",
       "newPassword": "Som@12345"
}
"User not registered, Please Register first!!"
502 Bad Gateway
Pass
5
{
       "email": "som@gmail",
       "password": "Som@123",
       "newPassword": "Som@12345"
}
"Invalid User or Wrong password!!"


401 Unauthorized
Pass
6
{
       "email": "som@gmail.com",
       "password": "Som@123",
       "newPassword": "Som@12345"
}


"Incorrect Password!!"


502 Bad Gateway
Pass




4.Customer Management
Get All Customers
Endpoint: /customers
Method: GET
Description: Retrieves a list of all customers.
Authentication: None
Response:
200 OK: Returns a list of customers.
Example:
[
   {
       "_id": "66ab331d97d120b6a67da9f3",
       "firstName": "Shiva",
       "lastName": "Gupta",
       "email": "shiva@gmail.com",
       "password": "76eedce0787f23a4924e897cdf7e6e20254b04f61a5f2a71875a5448477fc4586188f2394a",
       "mobileNumber": "8982511209",
       "dateOfBirth": "18/12/2000",
       "role": "customer",
       "isPrime": false,
       "primeExpireDate": "0001-01-01T00:00:00Z",
       "productsPurchased": [],
       "addresses": [
           {
               "_id": "66b06aeb19b253de5b363873",
               "fullName": "Shiva Gupta",
               "mobileNumber": "8982511209",
               "pincode": "471515",
               "line1": "near old bus stand,",
               "line2": "In frotn f government hospital",
               "landmark": "RC brother",
               "city": "LavKush nagar",
               "state": "MP",
               "country": "India",
               "isDefault": true
           },
           {
               "_id": "66b0748bf4c5a5864f10e7b5",
               "fullName": "Shiva Gupta",
               "mobileNumber": "9893792209",
               "pincode": "560000",
               "line1": "Stanza living lobito house,",
               "line2": "Mathikere",
               "landmark": "Devasandra bus stop",
               "city": "Bangalore",
               "state": "Karnataka",
               "country": "India",
               "isDefault": false
           }
       ],
       "wallet": 6105
   }
]


Get Single Customer
Endpoint: /customer/customerId={id}
Method: GET
Description: Retrieves details of a specific customer by their ID.
Authentication: None
Parameters:
id (path): The customer's ID.
Response:
200 OK: Returns the customer's details.
Example:
{
   "_id": "66c22829818ea8e8db645d55",
   "firstName": "Suryansh",
   "lastName": "Gupta",
   "email": "suryansh@gmail.com",
   "password": "f95202d1f0893f5cea42e4a17a63006672a91b364ed74908f1a57bdc13aaecc17e43485d9918555b",
   "mobileNumber": "7658934536",
   "dateOfBirth": "25/07/2000",
   "role": "customer",
   "isPrime": false,
   "primeExpireDate": "0001-01-01T00:00:00Z",
   "productsPurchased": [],
   "addresses": [],
   "wallet": 0
}


Create a Customer
Endpoint: /addCustomer
Method: POST
Description: Creates a new customer.
Created a new Cart for new customer
Created a new Wishlist for new customer
Created a new Search history for new customer
Created a new Product recommendation for new customer
Created a new Recently view product for new customer
Authentication: None
Request Body:


{
        "firstName": "Ravi",
        "lastName": "Kiran",
        "email": "ravi@gmail.com",
        "password": "Ravi@123",
        "mobileNumber": "3489765678",
        "dateOfBirth": "18/12/2002"
}


Response:
201 Created: Returns the created customer's ID.
Example:
{
   "InsertedID": "66c22faf818ea8e8db645d6f"
}
{
   "InsertedID": "66c22faf818ea8e8db645d6f"
}
{
   "InsertedID": "66c22faf818ea8e8db645d71"
}
{
   "InsertedID": "66c22fb0818ea8e8db645d73"
}
{
   "InsertedID": "66c22fb0818ea8e8db645d75"
}
{
   "InsertedID": "66c22fb1818ea8e8db645d77"
}
{
   "InsertedID": "66c22fb1818ea8e8db645d79"
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Invalid Email!!"
409 Conflict
Pass
2
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "somgmail.com",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}


"Invalid Email!!"
409 Conflict
Pass
3
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Invalid Email!!"
409 Conflict
Pass
4
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Password should be more than 7 characters!!"


409 Conflict
Pass
5
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som12345",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Password should contain atleast one super characters!!"


409 Conflict
Pass
6
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som@1234",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Your Mobile number must be of 10 digits!!"


409 Conflict
Pass
7
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som@1234",
       "mobileNumber": "6666666666",
       "dateOfBirth": "18/12/2025"
}
"Invalid DOB: Date is in the future!"


409 Conflict
Pass


Update a Customer
Endpoint: /updateCustomer/customerId={id}
Method: PUT
Authentication: Required
Description: Updates details of a specific customer by their ID.
Parameters:
id (path): The customer's ID.
Request Body:


{
        "firstName": "Ravi",
        "lastName": "Kiran",
        "email": "ravi@gmail.com",
        "password": "Ravi@123",
        "mobileNumber": "6789768978",
        "dateOfBirth": "18/12/2002"
}


Response:
200 OK: Returns the updated customer's details.
405 Method not found: Invalid url or method
Example:
{
   "_id": "66c22faf818ea8e8db645d6f",
   "firstName": "Ravi",
   "lastName": "Kiran",
   "email": "ravi@gmail.com",
   "password": "8dcd48c9a9fad927c8358a7e7e3b9ebe3b639d9827be73d4129794d716d2f2cb01c21452",
   "mobileNumber": "6789768978",
   "dateOfBirth": "18/12/2002",
   "role": "customer",
   "isPrime": false,
   "primeExpireDate": "0001-01-01T00:00:00Z",
   "productsPurchased": [],
   "addresses": [],
   "wallet": 0
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Invalid Email!!"
409 Conflict
Pass
2
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "somgmail.com",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}


"Invalid Email!!"
409 Conflict
Pass
3
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Invalid Email!!"
409 Conflict
Pass
4
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Password should be more than 7 characters!!"


409 Conflict
Pass
5
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som12345",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Password should contain atleast one super characters!!"


409 Conflict
Pass
6
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som@1234",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Your Mobile number must be of 10 digits!!"


409 Conflict
Pass
7
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som@1234",
       "mobileNumber": "6666666666",
       "dateOfBirth": "18/12/2025"
}
"Invalid DOB: Date is in the future!"


409 Conflict
Pass



Delete a Customer
Endpoint: /deleteCustomer/customerId={id}
Method: DELETE
Authentication: Required
Description: Deletes a specific customer by their ID.
Delete the customer from customer API
Delete the customer from user API as well
Parameters:
id (path): The customer's ID.
Response:
200 OK: Confirms deletion.
Example:
{
   "DeletedID": "66c22faf818ea8e8db645d6f"
}


5.Wallet Management
Add Money to Wallet
Endpoint: /addMoneyInWallet/customerId={id}
Method: POST
Authentication: Required
Description: Adds money to the customer's wallet.
Parameters:
id (path): The customer's ID.
Request Body:


{
    "money": 12000
}


Response:
200 OK: Confirms the addition of money to the wallet.
Example:
{
   "MatchedCount": 1,
   "ModifiedCount": 1,
   "UpsertedCount": 0,
   "UpsertedID": null
}
{
   "_id": "66c22faf818ea8e8db645d6f",
   "firstName": "Ravi",
   "lastName": "Kiran",
   "email": "ravi@gmail.com",
   "password": "1227391d5995e97a0766e55226d3777f29b1d26238c3a74be38e2fef58c7d633d5de6840",
   "mobileNumber": "6789768978",
   "dateOfBirth": "18/12/2002",
   "role": "customer",
   "isPrime": false,
   "primeExpireDate": "0001-01-01T00:00:00Z",
   "productsPurchased": null,
   "addresses": null,
   "wallet": 12000
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
    "money": 12000
}
"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass
2
{
    "moneys": 12000
}
Not give response
–
Pass
3
{
    "money": -12000
}
Can't be add negative money in the wallet!!


400 Bad Request
Pass


Withdraw Money from Wallet
Endpoint: /withdrawMoneyFromWallet/customerId={id}
Method: POST
Authentication: Required
Description: Withdraws money from the customer's wallet.
If withdrawAll is true, then money in wallet becomes zero
Parameters:
id (path): The customer's ID.
Request Body:
{
    "money": 1000,
    "withdrawAll": false
}


Response:
200 OK: Confirms the withdrawal of money from the wallet.
Example:
{
   "MatchedCount": 1,
   "ModifiedCount": 1,
   "UpsertedCount": 0,
   "UpsertedID": null
}
{
   "_id": "66c22faf818ea8e8db645d6f",
   "firstName": "Ravi",
   "lastName": "Kiran",
   "email": "ravi@gmail.com",
   "password": "1227391d5995e97a0766e55226d3777f29b1d26238c3a74be38e2fef58c7d633d5de6840",
   "mobileNumber": "6789768978",
   "dateOfBirth": "18/12/2002",
   "role": "customer",
   "isPrime": false,
   "primeExpireDate": "0001-01-01T00:00:00Z",
   "productsPurchased": null,
   "addresses": null,
   "wallet": 11000
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
    "money": 12000,
    "withdrawAll": false
}
"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass
2
{
    "moneys": 12000,
    "withdrawAlls": false
}
Not give response
–
Pass


{
    "money": 12000000,
    "withdrawAll": false
}
Not enough Money in Wallet!!


400 Bad Request
Pass


Buy Prime Membership
Endpoint: /buyPrimeMembership/customerId={id}
Method: POST
Authentication: Required
Description: Buys a Prime membership for the customer.
It has 4 subscriptions, 1 month(299), 3 months(599), 6 months(899) and              1 year(1099)
Features of Prime memberships are
Free delivery charge
Customer change change the delivery date according to availability
Only one day delivery

Parameters:
id (path): The customer's ID.
Request Body:
{
    "months":12
}


Response:
200 OK: Confirms the purchase of the Prime membership.
Example:
{
   "MatchedCount": 1,
   "ModifiedCount": 1,
   "UpsertedCount": 0,
   "UpsertedID": null
}
{
   "_id": "66c22faf818ea8e8db645d6f",
   "firstName": "Ravi",
   "lastName": "Kiran",
   "email": "ravi@gmail.com",
   "password": "1227391d5995e97a0766e55226d3777f29b1d26238c3a74be38e2fef58c7d633d5de6840",
   "mobileNumber": "6789768978",
   "dateOfBirth": "18/12/2002",
   "role": "customer",
   "isPrime": true,
   "primeExpireDate": "2025-08-18T23:19:35.150990413+05:30",
   "productsPurchased": null,
   "addresses": null,
   "wallet": 10901
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
    "month": 6
}
"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass
2
{
    "months": 12000
}
Not give response
–
Pass
3
{
    "month": 12
}
Not enough Money in Wallet!!
400 Bad Request
Pass




6.Seller Management
Get All Sellers
Endpoint: /sellers
Method: GET
Description: Retrieves a list of all sellers.
Authentication: None
Response:
200 OK: Returns a list of sellers.
Example:
[
   {
       "_id": "66a8be438b66d10a2315cf2e",
       "firstName": "Dhyey",
       "lastName": "Mistri",
       "email": "dhyey@gmail.com",
       "password": "8a5d92205be0c2b95df533290bc82237b452dc93f81f7b6c87ce9d332ee768c2daa0db5477",
       "mobileNumber": "0987654321",
       "dateOfBirth": "18/12/2001",
       "role": "seller",
       "isPrime": false,
       "productsListedIds": [
           "66b499c97dd2d413f5399ee9",
           "66a8d9c2697de17b2285d258",
           "66a9e09b63831c9fba7f84a2",
           "66a9e0ec63831c9fba7f84a5",
           "66aa0d90c9440c136263d090",
           "66aa320d5270e6700c8c94d1",
           "66ab273b03aa2deac83093e2",
           "66bca49f58a66a6eb4953126"
       ],
       "productsSoldIds": null
   },
   {
       "_id": "66c22878818ea8e8db645d62",
       "firstName": "Anmol",
       "lastName": "Gupta",
       "email": "anmol@gmail.com",
       "password": "c8d907bf70b2813a1616cf82f6f7ef4861afb5a2599ad98eacf1e90cae47daafc4c2ca128c9b",
       "mobileNumber": "9897678945",
       "dateOfBirth": "18/07/2000",
       "role": "seller",
       "isPrime": false,
       "productsListedIds": null,
       "productsSoldIds": null
   }
]


Get Single Seller
Endpoint: /seller/sellerId={id}
Method: GET
Description: Retrieves details of a specific seller by their ID.
Authentication: None
Parameters:
id (path): The seller's ID.
Response:
200 OK: Returns the seller's details.
Example:
{
   "_id": "66c22878818ea8e8db645d62",
   "firstName": "Anmol",
   "lastName": "Gupta",
   "email": "anmol@gmail.com",
   "password": "c8d907bf70b2813a1616cf82f6f7ef4861afb5a2599ad98eacf1e90cae47daafc4c2ca128c9b",
   "mobileNumber": "9897678945",
   "dateOfBirth": "18/07/2000",
   "role": "seller",
   "isPrime": false,
   "productsListedIds": null,
   "productsSoldIds": null
}


Create a Seller
Endpoint: /addSeller
Method: POST
Description: Creates a new seller.
New seller is added in seller API
New user is added in user API as a role seller
Creating a new report for the new seller
Authentication: None
Request Body:
{
       "firstName": "Dev",
       "lastName": "Gupta",
       "email": "dev@gmail.com",
       "password": "Dev@1234",
       "mobileNumber": "7654321234",
       "dateOfBirth": "26/03/2003"
}


Response:
201 Created: Returns the created seller's ID.
Example:
{
   "InsertedID": "66c2355b818ea8e8db645d8d"
}
{
   "InsertedID": "66c2355b818ea8e8db645d8d"
}
{
   "InsertedID": "66c2355c818ea8e8db645d8f"
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Invalid Email!!"
409 Conflict
Pass
2
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "somgmail.com",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}


"Invalid Email!!"
409 Conflict
Pass
3
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Invalid Email!!"
409 Conflict
Pass
4
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Password should be more than 7 characters!!"


409 Conflict
Pass
5
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som12345",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Password should contain atleast one super characters!!"


409 Conflict
Pass
6
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som@1234",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Your Mobile number must be of 10 digits!!"


409 Conflict
Pass
7
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som@1234",
       "mobileNumber": "6666666666",
       "dateOfBirth": "18/12/2025"
}
"Invalid DOB: Date is in the future!"


409 Conflict
Pass


Update a Seller
Endpoint: /updateSeller/sellerId={id}
Method: PUT
Authentication: Required
Description: Updates details of a specific seller by their ID.
Parameters:
id (path): The seller's ID.
Request Body:
{
       "firstName": "Dev",
       "lastName": "Gupta",
       "email": "dev@gmail.com",
       "password": "Dev@1234",
       "mobileNumber": "7654321234",
       "dateOfBirth": "26/03/2003"
}


Response:
200 OK: Returns the updated seller's details.
Example:
{
   "_id": "66c2355b818ea8e8db645d8d",
   "firstName": "Dev",
   "lastName": "Gupta",
   "email": "dev@gmail.com",
   "password": "c8a8501682059cc9ca65ee47b501294268f257ca79da2cf654ecea32edb769c161c6802e",
   "mobileNumber": "7654321234",
   "dateOfBirth": "26/03/2003",
   "role": "seller",
   "isPrime": false,
   "productsListedIds": null,
   "productsSoldIds": null
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Invalid Email!!"
409 Conflict
Pass
2
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "somgmail.com",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}


"Invalid Email!!"
409 Conflict
Pass
3
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Invalid Email!!"
409 Conflict
Pass
4
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som@123",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Password should be more than 7 characters!!"


409 Conflict
Pass
5
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som12345",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Password should contain atleast one super characters!!"


409 Conflict
Pass
6
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som@1234",
       "mobileNumber": "666666",
       "dateOfBirth": "18/12/2004"
}
"Your Mobile number must be of 10 digits!!"


409 Conflict
Pass
7
{
       "firstName": "Som",
       "lastName": "Gupta",
       "email": "som@gmail.com",
       "password": "Som@1234",
       "mobileNumber": "6666666666",
       "dateOfBirth": "18/12/2025"
}
"Invalid DOB: Date is in the future!"


409 Conflict
Pass



Delete a Seller
Endpoint: /deleteSeller/sellerId={id}
Method: DELETE
Authentication: Required
Description: Deletes a specific seller by their ID.
Also deleting the user from user API
Parameters:
id (path): The seller's ID.
Response:
200 OK: Confirms deletion.
Example:


{
   "DeletedId": "66c2355b818ea8e8db645d8d"
}


7.Product Management
Get All Products
Endpoint: /products/page={pageNumber}
Method: GET
Description: Retrieves a paginated list of all products.
Limit of product is 2 per page
Authentication: None
Parameters:
pageNumber (path): The page number for pagination.
Response:
200 OK: Returns a list of products.
Example:
[
   {
       "_id": "66a8d9c2697de17b2285d258",
       "productName": "IPhone 15",
       "aboutProduct": "DYNAMIC ISLAND COMES TO IPHONE 15 -- Dynamic Island bubbles up alerts and Live Activities -- so you don't miss them while you're doing something else. You can see who's calling, track your next ride, check your flight status, and so much more.",
       "sellerId": "66a8be438b66d10a2315cf2e",
       "brand": "Apple",
       "color": "",
       "releaseDate": "2024-07-30T00:00:00Z",
       "bestSeller": false,
       "newRelease": false,
       "replacePolicy": false,
       "returnPolicy": false,
       "maxRetailPrice": 79900,
       "sellingPrice": 70900,
       "discount": 11.26,
       "quantity": 3,
       "unitsSold": 0,
       "productCategories": [
           "Mobile",
           "Electronic"
       ],
       "averageRating": 0,
       "productImages": [
           "/assets/photos/iphone15.png"
       ],
       "productProperties": {
           "Memory Storage": "128 GB",
           "Model Name": "iPhone 15",
           "Operating System": "IOS",
           "Screen Size": "6.1 inches"
       },
       "numberOfReviews": 0
   },
   {
       "_id": "66a9e09b63831c9fba7f84a2",
       "productName": "iPhone 16",
       "aboutProduct": "DYNAMIC ISLAND COMES TO IPHONE 15 -- Dynamic Island bubbles up alerts and Live Activities -- so you don't miss them while you're doing something else. You can see who's calling, track your next ride, check your flight status, and so much more.",
       "sellerId": "66a8be438b66d10a2315cf2e",
       "brand": "Apple",
       "color": "",
       "releaseDate": "2024-07-31T00:00:00Z",
       "bestSeller": true,
       "newRelease": true,
       "replacePolicy": false,
       "returnPolicy": false,
       "maxRetailPrice": 99900,
       "sellingPrice": 80900,
       "discount": 19.02,
       "quantity": 4,
       "unitsSold": 6,
       "productCategories": [
           "Mobile",
           "Electronic"
       ],
       "averageRating": 0,
       "productImages": [
           "/assets/photos/iphone16.png"
       ],
       "productProperties": {
           "Memory Storage": "528 GB",
           "Model Name": "iPhone 16",
           "Operating System": "IOS",
           "Screen Size": "6.1 inches"
       },
       "numberOfReviews": 0
   }
]


Get Single Product for Authenticated Customer
Endpoint: /product/productId={productId}&customerId={id}
Method: GET
Authentication: Required
Description: Retrieves details of a specific product by its ID for a specific authenticated customer.
Update the customer’s recommended product API
Also add that product in the top 10 recently view product API for that customer
Parameters:
productId (path): The product's ID.
id (path): The customer's ID.
Response:
200 OK: Returns the product's details.
Example:
{
   "MatchedCount": 1,
   "ModifiedCount": 1,
   "UpsertedCount": 0,
   "UpsertedID": null
}
{
   "MatchedCount": 1,
   "ModifiedCount": 1,
   "UpsertedCount": 0,
   "UpsertedID": null
}
{
   "_id": "66a8d9c2697de17b2285d258",
   "productName": "IPhone 15",
   "aboutProduct": "DYNAMIC ISLAND COMES TO IPHONE 15 -- Dynamic Island bubbles up alerts and Live Activities -- so you don't miss them while you're doing something else. You can see who's calling, track your next ride, check your flight status, and so much more.",
   "sellerId": "66a8be438b66d10a2315cf2e",
   "brand": "Apple",
   "color": "",
   "releaseDate": "2024-07-30T00:00:00Z",
   "bestSeller": false,
   "newRelease": false,
   "replacePolicy": false,
   "returnPolicy": false,
   "maxRetailPrice": 79900,
   "sellingPrice": 70900,
   "discount": 11.26,
   "quantity": 3,
   "unitsSold": 0,
   "productCategories": [
       "Mobile",
       "Electronic"
   ],
   "averageRating": 0,
   "productImages": [
       "/assets/photos/iphone15.png"
   ],
   "productProperties": {
       "Memory Storage": "128 GB",
       "Model Name": "iPhone 15",
       "Operating System": "IOS",
       "Screen Size": "6.1 inches"
   },
   "numberOfReviews": 0
}



Create a Product
Endpoint: /addProduct/sellerId={id}
Method: POST
Authentication: Required
Description: Creates a new product under a specific seller.
Parameters:
id (path): The seller's ID.
Request Body:
{
        "productName": "Sony i7 5G",
        "aboutProduct": "DYNAMIC ISLAND COMES TO sony pixel 9A -- Dynamic Island bubbles up alerts and Live Activities -- so you don't miss them while you're doing something else. You can see who's calling, track your next ride, check your flight status, and so much more.",
        "brand": "Sony",
        "color": "Black",
        "bestSeller" : true,
        "maxRetailPrice": 15900,
        "sellingPrice": 14990,
        "quantity": 2,
        "productCategories": ["Mobile","Electronic"],
        "productImages": ["assets/photos/sony.png"],
        "productProperties": {
            "Operating System" : "Android",
            "Memory Storage" : "256 GB",
            "Screen Size" : "6.1 inches",
            "Model Name" : "Sony i7 5G"
        },
        "replacePolicy": true,
        "returnPolicy": true
}


Response:
201 Created: Returns the created product's ID.
Example:
{
   "InsertedID": "66c2499f4d1bd49184234263"
}
{
   "MatchedCount": 1,
   "ModifiedCount": 1,
   "UpsertedCount": 0,
   "UpsertedID": null
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
        "productName": "Sony i7 5G",
        "aboutProduct": "DYNAMIC ISLAND COMES TO sony pixel 9A -- Dynamic Island bubbles up alerts and Live Activities -- so you don't miss them while you're doing something else. You can see who's calling, track your next ride, check your flight status, and so much more.",
        "brand": "Sony",
        "color": "Black",
        "bestSeller" : true,
        "maxRetailPrice": 15900,
        "sellingPrice": 14990,
        "quantity": 2,
        "productCategories": ["Mobile","Electronic"],
        "productImages": ["assets/photos/sony.png"],
        "productProperties": {
            "Operating System" : "Android",
            "Memory Storage" : "256 GB",
            "Screen Size" : "6.1 inches",
            "Model Name" : "Sony i7 5G"
        },
        "replacePolicy": true,
        "returnPolicy": true
}
"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass


Update a Product
Endpoint: /updateProduct/productId={productId}&sellerId={id}
Method: PUT
Authentication: Required
Description: Updates details of a specific product by its ID under a specific seller.
Update the quantity of that product
Parameters:
productId (path): The product's ID.
id (path): The seller's ID.
Request Body:


{
        "productName": "Sony i7 5G",
        "aboutProduct": "DYNAMIC ISLAND COMES TO sony pixel 9A -- Dynamic Island bubbles up alerts and Live Activities -- so you don't miss them while you're doing something else. You can see who's calling, track your next ride, check your flight status, and so much more.",
        "brand": "Sony",
        "color": "Black",
        "bestSeller" : true,
        "maxRetailPrice": 15900,
        "sellingPrice": 14990,
        "quantity": 2,
        "productCategories": ["Mobile","Electronic"],
        "productImages": ["assets/photos/sony.png"],
        "productProperties": {
            "Operating System" : "Android",
            "Memory Storage" : "256 GB",
            "Screen Size" : "6.1 inches",
            "Model Name" : "Sony i7 5G"
        },
        "replacePolicy": true,
        "returnPolicy": true
}


Response:
200 OK: Returns the updated product's details.
Example:


{
   "_id": "66c2499f4d1bd49184234263",
   "productName": "BMW GT Model",
   "aboutProduct": "A toy for plaing, real model of BMW.",
   "sellerId": "66c2355b818ea8e8db645d8d",
   "brand": "HotWheels",
   "color": "White",
   "releaseDate": "2024-08-18T19:21:03.267Z",
   "bestSeller": true,
   "newRelease": true,
   "replacePolicy": true,
   "returnPolicy": true,
   "maxRetailPrice": 1299,
   "sellingPrice": 1199,
   "discount": 7.7,
   "quantity": 4,
   "unitsSold": 0,
   "productCategories": [
       "Car",
       "Toy"
   ],
   "averageRating": 0,
   "productImages": [
       "assets/photos/farrari.png"
   ],
   "productProperties": {
       "Body": "Metal",
       "Color": "Red",
       "Model Name": "BMW GT",
       "Screen Size": "4 inches"
   },
   "numberOfReviews": 0
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
        "productName": "Sony i7 5G",
        "aboutProduct": "DYNAMIC ISLAND COMES TO sony pixel 9A -- Dynamic Island bubbles up alerts and Live Activities -- so you don't miss them while you're doing something else. You can see who's calling, track your next ride, check your flight status, and so much more.",
        "brand": "Sony",
        "color": "Black",
        "bestSeller" : true,
        "maxRetailPrice": 15900,
        "sellingPrice": 14990,
        "quantity": 2,
        "productCategories": ["Mobile","Electronic"],
        "productImages": ["assets/photos/sony.png"],
        "productProperties": {
            "Operating System" : "Android",
            "Memory Storage" : "256 GB",
            "Screen Size" : "6.1 inches",
            "Model Name" : "Sony i7 5G"
        },
        "replacePolicy": true,
        "returnPolicy": true
}
"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass


Delete a Product
Endpoint: /deleteProduct/productId={productId}&sellerId={id}
Method: DELETE
Authentication: Required
Description: Deletes a specific product by its ID under a specific seller.
Parameters:
productId (path): The product's ID.
id (path): The seller's ID.
Response:
200 OK: Confirms deletion.
Example:


{
   "DeletedId": "66c2355b818ea8e8db645d8d"
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1


"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass


Search Products
Endpoint: /products/search/page={pageNumber}&customerId={id}
Method: GET
Description: Searches for products based on query parameters and returns paginated results for a specific customer.
Update the search history (store 10 recent search for the customer) 
Authentication: None
Parameters:
pageNumber (path): The page number for pagination.
id (path): The customer's ID.
Request Body:
{
	"Search_str" : “iphone”
}

Response:
200 OK: Returns a list of matching products.
Example:
[
   {
       "_id": "66a8d9c2697de17b2285d258",
       "productName": "IPhone 15",
       "aboutProduct": "DYNAMIC ISLAND COMES TO IPHONE 15 -- Dynamic Island bubbles up alerts and Live Activities -- so you don't miss them while you're doing something else. You can see who's calling, track your next ride, check your flight status, and so much more.",
       "sellerId": "66a8be438b66d10a2315cf2e",
       "brand": "Apple",
       "color": "",
       "releaseDate": "2024-07-30T00:00:00Z",
       "bestSeller": false,
       "newRelease": false,
       "replacePolicy": false,
       "returnPolicy": false,
       "maxRetailPrice": 79900,
       "sellingPrice": 70900,
       "discount": 11.26,
       "quantity": 3,
       "unitsSold": 0,
       "productCategories": [
           "Mobile",
           "Electronic"
       ],
       "averageRating": 0,
       "productImages": [
           "/assets/photos/iphone15.png"
       ],
       "productProperties": {
           "Memory Storage": "128 GB",
           "Model Name": "iPhone 15",
           "Operating System": "IOS",
           "Screen Size": "6.1 inches"
       },
       "numberOfReviews": 0
   },
   {
       "_id": "66a9e09b63831c9fba7f84a2",
       "productName": "iPhone 16",
       "aboutProduct": "DYNAMIC ISLAND COMES TO IPHONE 15 -- Dynamic Island bubbles up alerts and Live Activities -- so you don't miss them while you're doing something else. You can see who's calling, track your next ride, check your flight status, and so much more.",
       "sellerId": "66a8be438b66d10a2315cf2e",
       "brand": "Apple",
       "color": "",
       "releaseDate": "2024-07-31T00:00:00Z",
       "bestSeller": true,
       "newRelease": true,
       "replacePolicy": false,
       "returnPolicy": false,
       "maxRetailPrice": 99900,
       "sellingPrice": 80900,
       "discount": 19.02,
       "quantity": 4,
       "unitsSold": 6,
       "productCategories": [
           "Mobile",
           "Electronic"
       ],
       "averageRating": 0,
       "productImages": [
           "/assets/photos/iphone16.png"
       ],
       "productProperties": {
           "Memory Storage": "528 GB",
           "Model Name": "iPhone 16",
           "Operating System": "IOS",
           "Screen Size": "6.1 inches"
       },
       "numberOfReviews": 0
   }
]
{
   "MatchedCount": 1,
   "ModifiedCount": 1,
   "UpsertedCount": 0,
   "UpsertedID": null
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
	"Search_str" : “iphone”
}
"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass


Filter Products by Category
Endpoint: /products/categoryFilter/page={pageNumber}
Method: GET
Description: Filters products by category and returns paginated results.
Authentication: None
Parameters:
pageNumber (path): The page number for pagination.
Request Body:
{
   "categoryName" : "Car",
   "sortPriceAcending" : false,
   "sortPriceDecending" : true,
   "brand" : "HotWheels",
   "minPrice" : 0,
   "maxPrice" : 10000,
   "color" : "Red"
}


Response:
200 OK: Returns a list of products filtered by category.
Example:
[
   {
       "_id": "66bca49f58a66a6eb4953126",
       "productName": "Farrari 296 GTB Model",
       "aboutProduct": "A toy for plaing, real model of farrari.",
       "sellerId": "66a8be438b66d10a2315cf2e",
       "brand": "HotWheels",
       "color": "Red",
       "releaseDate": "2024-08-14T12:35:43.235Z",
       "bestSeller": true,
       "newRelease": true,
       "replacePolicy": true,
       "returnPolicy": true,
       "maxRetailPrice": 1299,
       "sellingPrice": 1199,
       "discount": 7.7,
       "quantity": 4,
       "unitsSold": 0,
       "productCategories": [
           "Car",
           "Toy"
       ],
       "averageRating": 0,
       "productImages": [
           "assets/photos/farrari.png"
       ],
       "productProperties": {
           "Body": "Metal",
           "Color": "Red",
           "Model Name": "Farrari 296 GTB",
           "Screen Size": "4 inches"
       },
       "numberOfReviews": 0
   }
]


Filter Products by Best Seller
Endpoint: /products/categoryFilter/bestSeller/page={pageNumber}
Method: GET
Description: Filters products by best seller status and returns paginated results.
Authentication: None
Parameters:
pageNumber (path): The page number for pagination.
Response:
200 OK: Returns a list of best-selling products.
Example:
[
   {
       "_id": "66a9e09b63831c9fba7f84a2",
       "productName": "iPhone 16",
       "aboutProduct": "DYNAMIC ISLAND COMES TO IPHONE 15 -- Dynamic Island bubbles up alerts and Live Activities -- so you don't miss them while you're doing something else. You can see who's calling, track your next ride, check your flight status, and so much more.",
       "sellerId": "66a8be438b66d10a2315cf2e",
       "brand": "Apple",
       "color": "",
       "releaseDate": "2024-07-31T00:00:00Z",
       "bestSeller": true,
       "newRelease": true,
       "replacePolicy": false,
       "returnPolicy": false,
       "maxRetailPrice": 99900,
       "sellingPrice": 80900,
       "discount": 19.02,
       "quantity": 4,
       "unitsSold": 6,
       "productCategories": [
           "Mobile",
           "Electronic"
       ],
       "averageRating": 0,
       "productImages": [
           "/assets/photos/iphone16.png"
       ],
       "productProperties": {
           "Memory Storage": "528 GB",
           "Model Name": "iPhone 16",
           "Operating System": "IOS",
           "Screen Size": "6.1 inches"
       },
       "numberOfReviews": 0
   },
   {
       "_id": "66aa320d5270e6700c8c94d1",
       "productName": "Google Pixel",
       "aboutProduct": "DYNAMIC ISLAND COMES TO google pixel -- Dynamic Island bubbles up alerts and Live Activities -- so you don't miss them while you're doing something else. You can see who's calling, track your next ride, check your flight status, and so much more.",
       "sellerId": "66a8be438b66d10a2315cf2e",
       "brand": "Google",
       "color": "White",
       "releaseDate": "2024-07-31T06:12:11.391Z",
       "bestSeller": true,
       "newRelease": true,
       "replacePolicy": true,
       "returnPolicy": true,
       "maxRetailPrice": 34900,
       "sellingPrice": 32990,
       "discount": 5.47,
       "quantity": 3,
       "unitsSold": 0,
       "productCategories": [
           "Mobile",
           "Electronic"
       ],
       "averageRating": 0,
       "productImages": [
           "/assets/photos/google.png"
       ],
       "productProperties": {
           "Memory Storage": "64 GB",
           "Model Name": "Google Pixel 7A 5G",
           "Operating System": "Android",
           "Screen Size": "6.1 inches"
       },
       "numberOfReviews": 0
   }
]


Filter Products by New Release
Endpoint: /products/categoryFilter/newRelease/page={pageNumber}
Method: GET
Description: Filters products by new release status and returns paginated results.
Authentication: None
Parameters:
pageNumber (path): The page number for pagination.
Response:
200 OK: Returns a list of newly released products.
Example:
[
   {
       "_id": "66c2499f4d1bd49184234263",
       "productName": "BMW GT Model",
       "aboutProduct": "A toy for plaing, real model of BMW.",
       "sellerId": "66c2355b818ea8e8db645d8d",
       "brand": "HotWheels",
       "color": "White",
       "releaseDate": "2024-08-18T19:21:03.267Z",
       "bestSeller": true,
       "newRelease": true,
       "replacePolicy": true,
       "returnPolicy": true,
       "maxRetailPrice": 1299,
       "sellingPrice": 1199,
       "discount": 7.7,
       "quantity": 4,
       "unitsSold": 0,
       "productCategories": [
           "Car",
           "Toy"
       ],
       "averageRating": 0,
       "productImages": [
           "assets/photos/farrari.png"
       ],
       "productProperties": {
           "Body": "Metal",
           "Color": "Red",
           "Model Name": "BMW GT",
           "Screen Size": "4 inches"
       },
       "numberOfReviews": 0
   },
   {
       "_id": "66bca49f58a66a6eb4953126",
       "productName": "Farrari 296 GTB Model",
       "aboutProduct": "A toy for plaing, real model of farrari.",
       "sellerId": "66a8be438b66d10a2315cf2e",
       "brand": "HotWheels",
       "color": "Red",
       "releaseDate": "2024-08-14T12:35:43.235Z",
       "bestSeller": true,
       "newRelease": true,
       "replacePolicy": true,
       "returnPolicy": true,
       "maxRetailPrice": 1299,
       "sellingPrice": 1199,
       "discount": 7.7,
       "quantity": 4,
       "unitsSold": 0,
       "productCategories": [
           "Car",
           "Toy"
       ],
       "averageRating": 0,
       "productImages": [
           "assets/photos/farrari.png"
       ],
       "productProperties": {
           "Body": "Metal",
           "Color": "Red",
           "Model Name": "Farrari 296 GTB",
           "Screen Size": "4 inches"
       },
       "numberOfReviews": 0
   }
]




8.Category Management
Get All Categories
Endpoint: /categories
Method: GET
Description: Retrieves a list of all categories.
Only for the frontend to display all the categories available
Authentication: None
Response:
200 OK: Returns a list of categories.
Example:
[
   {
       "_id": "66a9ca3133ee621888eb1acd",
       "categoryName": "Mobile",
       "brands": [
           "iPhone",
           "One Plus",
           "Red me",
           "Realme",
           "Google Pixel",
           "Poco",
           "Samsung",
           "Sony"
       ],
       "colors": [
           "Black",
           "White",
           "Gray",
           "Green"
       ],
       "priceRanges": {
           "0": "9999",
           "10000": "19999",
           "20000": "29999",
           "30000": "59999",
           "60000": "200000"
       }
   }
]



Get Single Category
Endpoint: /category/categoryId={categoryId}
Method: GET
Description: Retrieves details of a specific category by its ID.
Authentication: None
Parameters:
categoryId (path): The category's ID.
Response:
200 OK: Returns the category's details.
Example:
{
   "_id": "66a9ca3133ee621888eb1acd",
   "categoryName": "Mobile",
   "brands": [
       "iPhone",
       "One Plus",
       "Red me",
       "Realme",
       "Google Pixel",
       "Poco",
       "Samsung",
       "Sony"
   ],
   "colors": [
       "Black",
       "White",
       "Gray",
       "Green"
   ],
   "priceRanges": {
       "0": "9999",
       "10000": "19999",
       "20000": "29999",
       "30000": "59999",
       "60000": "200000"
   }
}


Create a Category
Endpoint: /addCategory
Method: POST
Authentication: None
Description: Creates a new category.
Request Body:


{
         "colors":["Black","White","Gray","Green"],  "priceRanges":{"0":"9999","10000":"19999","20000":"29999","30000":"59999","60000":"200000"},
     "categoryName":"Mobile",
     "brands":["iPhone","One Plus","Red me","Realme","Google Pixel","Poco","Samsung","Sony"]
}


Response:
201 Created: Returns the created category's ID.
Example:
{
   "InsertedID": "66c24e274d1bd4918423427a"
}

Delete a Category
Endpoint: /deleteCategory/categoryId={categoryId}
Method: DELETE
Authentication: None
Description: Deletes a specific category by its ID.
Parameters:
categoryId (path): The category's ID.
Response:
200 OK: Confirms deletion.
Example:
{
   "DeletedId": "66c24e274d1bd4918423427a"
}



9.Cart Management
Get All Carts
Endpoint: /carts
Method: GET
Description: Retrieves a list of all carts.
Authentication: None
Response:
200 OK: Returns a list of carts.
Example:
[
   {
       "_id": "66ab331e97d120b6a67da9f5",
       "customerId": "66ab331d97d120b6a67da9f3",
       "cartItems": [
           {
               "productId": "66aa0d90c9440c136263d090",
               "productName": "Redmi",
               "productPrice": 12999,
               "quantityInCart": 4,
               "selectedForBuying": false,
               "productImage": "/assets/photos/redmi.png",
               "productInStock": true
           },
           {
               "productId": "66a8d9c2697de17b2285d258",
               "productName": "IPhone 15",
               "productPrice": 70900,
               "quantityInCart": 1,
               "selectedForBuying": false,
               "productImage": "/assets/photos/iphone15.png",
               "productInStock": true
           }
       ],
       "numberOfProduct": 0,
       "totalAmount": 0
   },
   {
       "_id": "66b4537d0208bbd82ee589ac",
       "customerId": "66b4537d0208bbd82ee589aa",
       "cartItems": [],
       "numberOfProduct": 0,
       "totalAmount": 0
   },
   {
       "_id": "66c222ed56fe9801831cd73e",
       "customerId": "66c222ed56fe9801831cd73c",
       "cartItems": null,
       "numberOfProduct": 0,
       "totalAmount": 0
   },
   {
       "_id": "66c2282e818ea8e8db645d58",
       "customerId": "66c22829818ea8e8db645d55",
       "cartItems": [],
       "numberOfProduct": 0,
       "totalAmount": 0
   },
   {
       "_id": "66c22faf818ea8e8db645d71",
       "customerId": "66c22faf818ea8e8db645d6f",
       "cartItems": [],
       "numberOfProduct": 0,
       "totalAmount": 0
   }
]


Get Cart by Customer ID
Endpoint: /cart/customerId={id}
Method: GET
Authentication: Required
Description: Retrieves the cart for a specific customer by their ID.
Initially that product is checked for buying
Customer can toggle check in their cart weather that product is added in the buying list or not
Also check that each product is out of stock or not
Parameters:
id (path): The customer's ID.
Response:
200 OK: Returns the customer's cart.
Example:
{
   "_id": "66c22faf818ea8e8db645d71",
   "customerId": "66c22faf818ea8e8db645d6f",
   "cartItems": [],
   "numberOfProduct": 0,
   "totalAmount": 0
}


Add Item to Cart
Endpoint: /addToCart/customerId={id}
Method: POST
Authentication: Required
Description: Adds an item to the customer's cart.
Customer can not add any out of stock product in their cart
Parameters:
id (path): The customer's ID.
Request Body:
{
    "customerId" : "66c22faf818ea8e8db645d6f",
    "productId" : "66c2499f4d1bd49184234263"
    "quantity" : 2
}


Response:
201 Created: Returns the updated cart.
Example:
{
   "_id": "66c22faf818ea8e8db645d71",
   "customerId": "66c22faf818ea8e8db645d6f",
   "cartItems": [
       {
           "productId": "66c2499f4d1bd49184234263",
           "productName": "BMW GT Model",
           "productPrice": 1199,
           "quantityInCart": 2,
           "selectedForBuying": true,
           "productImage": "assets/photos/farrari.png",
           "productInStock": true
       }
   ],
   "numberOfProduct": 2,
   "totalAmount": 2398
}


Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
    "customerId" : "66ab331d97d120b6a67da9f3",
    "productId" : "66aa0d90c9440c136263d090",
    "quantity" : 2
}
"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass
2
{
    "customerId" : "66ab331d97d120b6a67da9f3",
    "productId" : "66aa0d90c9440c136263d090",
    "quantity" : 2
}
"Out of Stock or Less Quantity Available!!"
(If that product is out of stock)
400 Bad Request
Pass


Remove All Items from Cart
Endpoint: /removeAllItemsFromCart/customerId={id}
Method: POST
Authentication: Required
Description: Removes all items from the customer's cart.
Parameters:
id (path): The customer's ID.
Response:
200 OK: Confirms removal of all items.
Example:


{
   "MatchedCount": 1,
   "ModifiedCount": 1,
   "UpsertedCount": 0,
   "UpsertedID": null
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1


"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass
2


"Cart is already empty!!"
(If no product is in the cart)
400 Bad Request
Pass


Toggle Product to Buy
Endpoint: /cart/toggleProductToBuy/customerId={id}
Method: POST
Authentication: Required
Description: Toggles the purchase status of a product in the cart.
Parameters:
id (path): The customer's ID.
Request Body:
{
    "customerId" : "66c22faf818ea8e8db645d6f",
    "productId" : "66c2499f4d1bd49184234263"
}


Response:
200 OK: Returns the updated cart.
Example:

{
   "_id": "66c22faf818ea8e8db645d71",
   "customerId": "66c22faf818ea8e8db645d6f",
   "cartItems": [
       {
           "productId": "66c2499f4d1bd49184234263",
           "productName": "BMW GT Model",
           "productPrice": 1199,
           "quantityInCart": 2,
           "selectedForBuying": false,
           "productImage": "assets/photos/farrari.png",
           "productInStock": true
       }
   ],
   "numberOfProduct": 0,
   "totalAmount": 0
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
    "customerId" : "66c22faf818ea8e8db645d6f",
    "productId" : "66c2499f4d1bd49184234263"
}
"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass
2


"Product Out of Stock, Unable to toggle to Buy!!"
400 Bad Request
Pass


Increase Product Quantity in Cart
Endpoint: /cart/increasingProductQuantity/customerId={id}
Method: POST
Authentication: Required
Description: Increases the quantity of a product in the customer's cart.
Only if that product’s quantity is available in the inventory
Not decreasing the quantity from inventory, it will when order placed / payment initiated
Parameters:
id (path): The customer's ID.
Request Body:
{
    "customerId" : "66c22faf818ea8e8db645d6f",
    "productId" : "66c2499f4d1bd49184234263"
}


Response:
200 OK: Returns the updated cart.
Example:



{
   "_id": "66c22faf818ea8e8db645d71",
   "customerId": "66c22faf818ea8e8db645d6f",
   "cartItems": [
       {
           "productId": "66c2499f4d1bd49184234263",
           "productName": "BMW GT Model",
           "productPrice": 1199,
           "quantityInCart": 3,
           "selectedForBuying": true,
           "productImage": "assets/photos/farrari.png",
           "productInStock": true
       }
   ],
   "numberOfProduct": 3,
   "totalAmount": 3597
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
    "customerId" : "66c22faf818ea8e8db645d6f",
    "productId" : "66c2499f4d1bd49184234263"
}
"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass
2
{
    "customerId" : "66c22faf818ea8e8db645d6f",
    "productId" : "66c2499f4d1bd49184234263"
}
"Out of Stock or Less Quantity Available in Inventory!!"
400 Bad Request
Pass


Decrease Product Quantity in Cart
Endpoint: /cart/decreasingProductQuantity/customerId={id}
Method: POST
Authentication: Required
Description: Decreases the quantity of a product in the customer's cart.
If it decreases to zero product will be removed from the customer’s cart
Parameters:
id (path): The customer's ID.
Request Body:
{
    "customerId" : "66ab331d97d120b6a67da9f3",
    "productId" : "66aa0d90c9440c136263d090"
}


Response:
200 OK: Returns the updated cart.
Example:


{
   "_id": "66c22faf818ea8e8db645d71",
   "customerId": "66c22faf818ea8e8db645d6f",
   "cartItems": [
       {
           "productId": "66c2499f4d1bd49184234263",
           "productName": "BMW GT Model",
           "productPrice": 1199,
           "quantityInCart": 3,
           "selectedForBuying": true,
           "productImage": "assets/photos/farrari.png",
           "productInStock": true
       }
   ],
   "numberOfProduct": 3,
   "totalAmount": 3597
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
    "customerId" : "66c22faf818ea8e8db645d6f",
    "productId" : "66c2499f4d1bd49184234263"
}
"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass



10.Wishlist Management
Get All Wishlists
Endpoint: /wishlists
Method: GET
Description: Retrieves all wishlists in the system.
Authentication: None
Response:
200 OK: A list of wishlists.
500 Internal Server Error: If there is a server error.
Example Response:
[
   {
       "_id": "66acb02c3fb43d5d8661a83c",
       "customerId": "66ab331d97d120b6a67da9f3",
       "wishlistItems": [
           {
               "productId": "66aa0d90c9440c136263d090",
               "productName": "Redmi",
               "productImage": "/assets/photos/redmi.png",
               "productPrice": 12999,
               "productInStock": true
           },
           {
               "productId": "66a8d9c2697de17b2285d258",
               "productName": "IPhone 15",
               "productImage": "/assets/photos/iphone15.png",
               "productPrice": 70900,
               "productInStock": true
           },
           {
               "productId": "66ab273b03aa2deac83093e2",
               "productName": "Google Pixel",
               "productImage": "/assets/photos/google.png",
               "productPrice": 42990,
               "productInStock": false
           }
       ],
       "numberOfProduct": 3
   }
]


Get Wishlist by Customer ID
Endpoint: /wishlists/customerId={id}
Method: GET
Description: Retrieves the wishlist for a specific customer.
Authentication: Required.
Parameters:
id (path): The customer's ID.
Response:
200 OK: The wishlist for the given customer ID.
404 Not Found: If no wishlist is found for the customer.
500 Internal Server Error: If there is a server error.
Example
{
   "_id": "66c22fb0818ea8e8db645d73",
   "customerId": "66c22faf818ea8e8db645d6f",
   "wishlistItems": [],
   "numberOfProduct": 0
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1


"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass



Add Product to Wishlist
Endpoint: /addToWishlist/customerId={id}
Method: POST
Description: Adds a product to the wishlist of a specific customer.
If the product is out of stock, then also it will add in wishlist
Authentication: Required.
Parameters:
id (path): The customer's ID.
Request Body:
{
    "customerId" : "66c22faf818ea8e8db645d6f",
    "productId" : "66c2499f4d1bd49184234263"
}


Response:
201 Created: Product added to wishlist successfully.
400 Bad Request: If the request body is invalid.
500 Internal Server Error: If there is a server error.
Example
{
   "_id": "66c22fb0818ea8e8db645d73",
   "customerId": "66c22faf818ea8e8db645d6f",
   "wishlistItems": [
       {
           "productId": "66c2499f4d1bd49184234263",
           "productName": "BMW GT Model",
           "productImage": "assets/photos/farrari.png",
           "productPrice": 1199,
           "productInStock": true
       }
   ],
   "numberOfProduct": 1
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
    "customerId" : "66c22faf818ea8e8db645d6f",
    "productId" : "66c2499f4d1bd49184234263"
}
"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass
2
{
    "customerId" : "66c22faf818ea8e8db645d6f",
    "productId" : "66c2499f4d1bd49184234263"
}
"Product is already in your Wishlist!!"
409 Conflict
Pass


Remove All Items from Wishlist
Endpoint: /RemoveAllItemsFromWishlist/customerId={id}
Method: POST
Description: Removes all products from the wishlist of a specific customer.
Authentication: Required.
Parameters:
id (path): The customer's ID.
Response:
200 OK: All items removed from wishlist successfully.
404 Not Found: If no wishlist is found for the customer.
500 Internal Server Error: If there is a server error.
Example
{
   "MatchedCount": 1,
   "ModifiedCount": 1,
   "UpsertedCount": 0,
   "UpsertedID": null
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1


"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass
2


"Wishlist is already empty!!"
400 Bad Request
Pass


Remove a Product from Wishlist
Endpoint: /wishlists/removeProductFromWishlist/customerId={id}
Method: POST
Description: Removes a specific product from the wishlist of a specific customer.
Authentication: Required.
Parameters:
id (path): The customer's ID.
Request Body:
{
    "customerId" : "66c22faf818ea8e8db645d6f",
    "productId" : "66c2499f4d1bd49184234263"
}


Response:
200 OK: Product removed from wishlist successfully.
404 Not Found: If the product is not found in the wishlist.
500 Internal Server Error: If there is a server error.
Example
{
   "_id": "66c22fb0818ea8e8db645d73",
   "customerId": "66c22faf818ea8e8db645d6f",
   "wishlistItems": [
       {
           "productId": "66c2499f4d1bd49184234263",
           "productName": "BMW GT Model",
           "productImage": "assets/photos/farrari.png",
           "productPrice": 1199,
           "productInStock": true
       }
   ],
   "numberOfProduct": 1
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
    "customerId" : "66c22faf818ea8e8db645d6f",
    "productId" : "66c2499f4d1bd49184234263"
}
"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass



11. Address Management
Get All Addresses
Endpoint: /addresses
Method: GET
Description: Retrieves all addresses stored in the system.
Authentication: None
Response:
200 OK: A list of addresses.
500 Internal Server Error: If there is a server error.
Example
[
   {
       "_id": "66b06aeb19b253de5b363873",
       "fullName": "Shiva Gupta",
       "mobileNumber": "8982511209",
       "pincode": "471515",
       "line1": "near old bus stand,",
       "line2": "In frotn f government hospital",
       "landmark": "RC brother",
       "city": "LavKush nagar",
       "state": "MP",
       "country": "India",
       "isDefault": true
   },
   {
       "_id": "66b0748bf4c5a5864f10e7b5",
       "fullName": "Shiva Gupta",
       "mobileNumber": "9893792209",
       "pincode": "560000",
       "line1": "Stanza living lobito house,",
       "line2": "Mathikere",
       "landmark": "Devasandra bus stop",
       "city": "Bangalore",
       "state": "Karnataka",
       "country": "India",
       "isDefault": false
   },
   {
       "_id": "66b9a11ec046307c7abca648",
       "fullName": "Ram Gupta",
       "mobileNumber": "1234567890",
       "pincode": "560000",
       "line1": "....,",
       "line2": "...",
       "landmark": "...",
       "city": "Bangalore",
       "state": "Karnataka",
       "country": "India",
       "isDefault": false
   }
]


Add Address
Endpoint: /addAddress/customerId={id}
Method: POST
Description: Adds a new address for a specified customer.
Address will be added in the customer API
Initially the most recent added address will be set is default and rest will be uncheck from the default automatically
All the validation for the required field is applied
Customer can set the new address as default, then all previous default address will be false
Authentication: Required
Parameters:
id (path): The customer's ID.
Request Body:
{
    "fullName" : "Shiva Gupta",
    "mobileNumber" : "1234567890",
    "pincode" : "560000",
    "line1" : "....,",
    "line2" : "...",
    "landmark" : "...",
    "city" : "Bangalore",
    "state" : "Karnataka",
    "country" : "India",
   “isDefault” : true
}


Response:
201 Created: Address successfully added.
400 Bad Request: If the request is invalid.
500 Internal Server Error: If there is a server error.
Example:
{
   "InsertedID": "66c253c8522685bb6b7f75c6"
}
{
   "MatchedCount": 1,
   "ModifiedCount": 1,
   "UpsertedCount": 0,
   "UpsertedID": null
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
    "fullName" : "Shiva Gupta",
    "mobileNumber" : "1234567890",
    "pincode" : "560000",
    "line1" : "....,",
    "line2" : "...",
    "landmark" : "...",
    "city" : "Bangalore",
    "state" : "Karnataka",
    "country" : "India",
   “isDefault” : true
}
"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass
2
{
    "fullName" : "",
    "mobileNumber" : "1234567890",
    "pincode" : "560000",
    "line1" : "....,",
    "line2" : "...",
    "landmark" : "...",
    "city" : "Bangalore",
    "state" : "Karnataka",
    "country" : "India",
   “isDefault” : true
}
“Full name is empty!!”
409 Conflict
Pass
3
{
    "fullName" : "Shiva Gupta",
    "mobileNumber" : "1234567890",
    "pincode" : "",
    "line1" : "....,",
    "line2" : "...",
    "landmark" : "...",
    "city" : "Bangalore",
    "state" : "Karnataka",
    "country" : "India",
   “isDefault” : true
}


“Pincode is empty!!”
409 Conflict
Pass
4
{
    "fullName" : "Shiva Gupta",
    "mobileNumber" : "1234567890",
    "pincode" : "560000",
    "line1" : "....,",
    "line2" : "...",
    "landmark" : "...",
    "city" : "Bangalore",
    "state" : "Karnataka",
    "country" : "India",
   “isDefault” : true
}
“Line 1 is empty!!”
409 Conflict
Pass
5
{
    "fullName" : "Shiva Gupta",
    "mobileNumber" : "1234567890",
    "pincode" : "560000",
    "line1" : "stanza living",
    "line2" : "...",
    "landmark" : "...",
    "city" : "Bangalore",
    "state" : "Karnataka",
    "country" : "India",
   “isDefault” : true
}
“Landmark is empty!!”


409 Conflict
Pass
6
{
    "fullName" : "Shiva Gupta",
    "mobileNumber" : "1234567890",
    "pincode" : "560000",
    "line1" : "stanza living",
    "line2" : "...",
    "landmark" : "rostok house",
    "city" : "",
    "state" : "Karnataka",
    "country" : "India",
   “isDefault” : true
}
"City is empty!!"


409 Conflict
Pass
7
{
    "fullName" : "Shiva Gupta",
    "mobileNumber" : "1234567890",
    "pincode" : "560000",
    "line1" : "stanza living",
    "line2" : "...",
    "landmark" : "rostok house",
    "city" : "Bangalore",
    "state" : "",
    "country" : "India",
   “isDefault” : true
}
"State is empty!!"
409 Conflict
Pass
8
{
    "fullName" : "Shiva Gupta",
    "mobileNumber" : "1234567890",
    "pincode" : "560000",
    "line1" : "stanza living",
    "line2" : "...",
    "landmark" : "rostok house",
    "city" : "Bangalore",
    "state" : "Karnataka",
    "country" : "",
   “isDefault” : true
}
"Country is empty!!"
409 Conflict
Pass
9
{
    "fullName" : "Shiva Gupta",
    "mobileNumber" : "12345",
    "pincode" : "560000",
    "line1" : "stanza living",
    "line2" : "...",
    "landmark" : "rostok house",
    "city" : "Bangalore",
    "state" : "Karnataka",
    "country" : "India",
   “isDefault” : true
}
"Your Mobile number must be of 10 digits!!"


409 Conflict
Pass


Update Address
Endpoint: /updateAddress/addressId={addressId}&customerId={id}
Method: PUT
Description: Updates an existing address for a specified customer.
Any field can’t be empty
All the validation for the required field is applied
Address will also be updated in the customer’s API
Authentication: Required
Parameters:
addressId (path): The address ID.
id (path): The customer's ID.
Request Body
{
    "fullName" : "Shiva Gupta",
    "mobileNumber" : "1234567890",
    "pincode" : "560000",
    "line1" : "....,",
    "line2" : "...",
    "landmark" : "...",
    "city" : "Bangalore",
    "state" : "Karnataka",
    "country" : "India"
}


Response:
200 OK: Address successfully updated.
400 Bad Request: If the request is invalid.
404 Not Found: If the address or customer is not found.
500 Internal Server Error: If there is a server error.
Example:
{
   "MatchedCount": 1,
   "ModifiedCount": 1,
   "UpsertedCount": 0,
   "UpsertedID": null
}
{
   "MatchedCount": 1,
   "ModifiedCount": 1,
   "UpsertedCount": 0,
   "UpsertedID": null
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
    "fullName" : "Shiva Gupta",
    "mobileNumber" : "1234567890",
    "pincode" : "560000",
    "line1" : "....,",
    "line2" : "...",
    "landmark" : "...",
    "city" : "Bangalore",
    "state" : "Karnataka",
    "country" : "India",
   “isDefault” : true
}
"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass
2
{
    "fullName" : "",
    "mobileNumber" : "1234567890",
    "pincode" : "560000",
    "line1" : "....,",
    "line2" : "...",
    "landmark" : "...",
    "city" : "Bangalore",
    "state" : "Karnataka",
    "country" : "India",
   “isDefault” : true
}
“Full name is empty!!”
409 Conflict
Pass
3
{
    "fullName" : "Shiva Gupta",
    "mobileNumber" : "1234567890",
    "pincode" : "",
    "line1" : "....,",
    "line2" : "...",
    "landmark" : "...",
    "city" : "Bangalore",
    "state" : "Karnataka",
    "country" : "India",
   “isDefault” : true
}


“Pincode is empty!!”
409 Conflict
Pass
4
{
    "fullName" : "Shiva Gupta",
    "mobileNumber" : "1234567890",
    "pincode" : "560000",
    "line1" : "....,",
    "line2" : "...",
    "landmark" : "...",
    "city" : "Bangalore",
    "state" : "Karnataka",
    "country" : "India",
   “isDefault” : true
}
“Line 1 is empty!!”
409 Conflict
Pass
5
{
    "fullName" : "Shiva Gupta",
    "mobileNumber" : "1234567890",
    "pincode" : "560000",
    "line1" : "stanza living",
    "line2" : "...",
    "landmark" : "...",
    "city" : "Bangalore",
    "state" : "Karnataka",
    "country" : "India",
   “isDefault” : true
}
“Landmark is empty!!”


409 Conflict
Pass
6
{
    "fullName" : "Shiva Gupta",
    "mobileNumber" : "1234567890",
    "pincode" : "560000",
    "line1" : "stanza living",
    "line2" : "...",
    "landmark" : "rostok house",
    "city" : "",
    "state" : "Karnataka",
    "country" : "India",
   “isDefault” : true
}
"City is empty!!"


409 Conflict
Pass
7
{
    "fullName" : "Shiva Gupta",
    "mobileNumber" : "1234567890",
    "pincode" : "560000",
    "line1" : "stanza living",
    "line2" : "...",
    "landmark" : "rostok house",
    "city" : "Bangalore",
    "state" : "",
    "country" : "India",
   “isDefault” : true
}
"State is empty!!"
409 Conflict
Pass
8
{
    "fullName" : "Shiva Gupta",
    "mobileNumber" : "1234567890",
    "pincode" : "560000",
    "line1" : "stanza living",
    "line2" : "...",
    "landmark" : "rostok house",
    "city" : "Bangalore",
    "state" : "Karnataka",
    "country" : "",
   “isDefault” : true
}
"Country is empty!!"
409 Conflict
Pass
9
{
    "fullName" : "Shiva Gupta",
    "mobileNumber" : "12345",
    "pincode" : "560000",
    "line1" : "stanza living",
    "line2" : "...",
    "landmark" : "rostok house",
    "city" : "Bangalore",
    "state" : "Karnataka",
    "country" : "India",
   “isDefault” : true
}
"Your Mobile number must be of 10 digits!!"


409 Conflict
Pass


Delete Address
Endpoint: /deleteAddress/addressId={addressId}&customerId={id}
Method: DELETE
Description: Deletes an existing address for a specified customer and addressId.
Authentication: Required
Parameters:
addressId (path): The address ID.
id (path): The customer's ID.
Response:
200 OK: Address successfully deleted.
404 Not Found: If the address or customer is not found.
500 Internal Server Error: If there is a server error.
Example :
{
   "DeletedCount": 1
}
{
   "MatchedCount": 1,
   "ModifiedCount": 1,
   "UpsertedCount": 0,
   "UpsertedID": null
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
    "fullName" : "Shiva Gupta",
    "mobileNumber" : "1234567890",
    "pincode" : "560000",
    "line1" : "....,",
    "line2" : "...",
    "landmark" : "...",
    "city" : "Bangalore",
    "state" : "Karnataka",
    "country" : "India",
   “isDefault” : true
}
"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass



12. Order Management
Get All Orders
Endpoint: /orders
Method: GET
Description: Retrieves all orders in the system.
Authentication: None
Response:
200 OK: A list of orders.
500 Internal Server Error: If there is a server error.
Example
[
   {
       "_id": "66b9b7ffe2932a902210149f",
       "customerId": "66b4537d0208bbd82ee589aa",
       "productIds": [
           "66aa0d90c9440c136263d090",
           "66a9e09b63831c9fba7f84a2"
       ],
       "productNames": [
           "Redmi",
           "iPhone 16"
       ],
       "orderQuantitys": [
           2,
           1
       ],
       "orderedDate": "2024-08-12T07:22:41.445Z",
       "deliveredDates": [
           [
               "2024-08-12T07:22:41.847Z",
               "2024-08-12T07:23:35.686Z",
               "2024-08-12T07:28:12.397Z",
               "2024-08-12T08:28:53.161Z",
               "2024-08-12T08:31:17.71Z",
               "2024-08-12T08:31:55.966Z",
               "2024-08-12T08:33:56.911Z",
               "2024-08-12T08:34:47.456Z"
           ],
           [
               "2024-08-12T07:22:42.679Z",
               "2024-08-16T07:22:42.679Z"
           ]
       ],
       "status": [
           [
               "Order Confirmed",
               "Order Delivered",
               "Replace Confirmed",
               "Cancelled",
               "Return Confirmed",
               "Cancelled",
               "Return Confirmed",
               "Refunded"
           ],
           [
               "Order Confirmed"
           ]
       ],
       "shippingAddress": {
           "_id": "66b9a11ec046307c7abca648",
           "fullName": "Ram Gupta",
           "mobileNumber": "1234567890",
           "pincode": "560000",
           "line1": "....,",
           "line2": "...",
           "landmark": "...",
           "city": "Bangalore",
           "state": "Karnataka",
           "country": "India",
           "isDefault": false
       },
       "paymentDetails": {
           "_id": "66b9b841e2932a90221014a3",
           "paymentMethod": "Wallet",
           "totalAmount": 106898,
           "userId": "66b4537d0208bbd82ee589aa",
           "cardNumber": "",
           "cardExpiryDate": "",
           "nameOnCard": "",
           "upiId": ""
       },
       "priceDetail": {
           "listPrice": 129900,
           "sellingPrice": 106898,
           "deliveryCharge": 0,
           "totalAmount": 106898
       },
       "totalAmount": 106898
   }
]


Get Orders by Customer ID
Endpoint: /orders/customerId={id}
Method: GET
Description: Retrieves all orders for a specified customer.
Used for the order history for the customer’s profile
Authentication: Required
Parameters:
id (path): The customer's ID.
Response:
200 OK: A list of orders for the specified customer.
404 Not Found: If no orders are found for the customer.
500 Internal Server Error: If there is a server error.
Example


[
   {
       "_id": "66b9b7ffe2932a902210149f",
       "customerId": "66b4537d0208bbd82ee589aa",
       "productIds": [
           "66aa0d90c9440c136263d090",
           "66a9e09b63831c9fba7f84a2"
       ],
       "productNames": [
           "Redmi",
           "iPhone 16"
       ],
       "orderQuantitys": [
           2,
           1
       ],
       "orderedDate": "2024-08-12T07:22:41.445Z",
       "deliveredDates": [
           [
               "2024-08-12T07:22:41.847Z",
               "2024-08-12T07:23:35.686Z",
               "2024-08-12T07:28:12.397Z",
               "2024-08-12T08:28:53.161Z",
               "2024-08-12T08:31:17.71Z",
               "2024-08-12T08:31:55.966Z",
               "2024-08-12T08:33:56.911Z",
               "2024-08-12T08:34:47.456Z"
           ],
           [
               "2024-08-12T07:22:42.679Z",
               "2024-08-16T07:22:42.679Z"
           ]
       ],
       "status": [
           [
               "Order Confirmed",
               "Order Delivered",
               "Replace Confirmed",
               "Cancelled",
               "Return Confirmed",
               "Cancelled",
               "Return Confirmed",
               "Refunded"
           ],
           [
               "Order Confirmed"
           ]
       ],
       "shippingAddress": {
           "_id": "66b9a11ec046307c7abca648",
           "fullName": "Ram Gupta",
           "mobileNumber": "1234567890",
           "pincode": "560000",
           "line1": "....,",
           "line2": "...",
           "landmark": "...",
           "city": "Bangalore",
           "state": "Karnataka",
           "country": "India",
           "isDefault": false
       },
       "paymentDetails": {
           "_id": "66b9b841e2932a90221014a3",
           "paymentMethod": "Wallet",
           "totalAmount": 106898,
           "userId": "66b4537d0208bbd82ee589aa",
           "cardNumber": "",
           "cardExpiryDate": "",
           "nameOnCard": "",
           "upiId": ""
       },
       "priceDetail": {
           "listPrice": 129900,
           "sellingPrice": 106898,
           "deliveryCharge": 0,
           "totalAmount": 106898
       },
       "totalAmount": 106898
   }
]

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1


"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass


Buy Now
Endpoint: /buyNow/customerId={id}
Method: POST
Description: Creates a new order for the specified customer to buy a single or multiple product immediately.
It will also create the bill that is detail charges for the products
Later will be proceed for the payment
Shipping address will be the default address of that customer
Order status look like:
Order Confirmed
Canceled
Order Delivered
Replace Confirmed
Canceled
Replace Order Delivered
Return Confirmed
Canceled
Refund Initiated / Refunded
Authentication: Required
Parameters:
id (path): The customer's ID.
Request Body:
{
    "customerId": "66ab331d97d120b6a67da9f3",
    "productIds": ["66aa320d5270e6700c8c94d1"],
    "productNames" : ["Google Pixel"],
    "orderQuantitys": [1]
}


Response:
201 Created: Order successfully created.
400 Bad Request: If the request is invalid.
500 Internal Server Error: If there is a server error.
Example
{
   "InsertedID": "66c2590a512bf19c6acef70b"
}
"66c2590a512bf19c6acef70b"
{
   "listPrice": 34900,
   "sellingPrice": 32990,
   "deliveryCharge": 0,
   "totalAmount": 32990
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
    "customerId": "66ab331d97d120b6a67da9f3",
    "productIds": ["66aa320d5270e6700c8c94d1"],
    "productNames" : ["Google Pixel"],
    "orderQuantitys": [1]
}
"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass
2
{
    "customerId": "66ab331d97d120b6a67da9f3",
    "productIds": ["66aa320d5270e6700c8c94d1"],
    "productNames" : ["Google Pixel"],
    "orderQuantitys": [1]
}
“Google Pixel is Out of Stock, Please Reduce the quantity of the product!!"
409 Conflict
Pass


Continue to Payment
Endpoint: /continueToPayment/orderId={orderId}&customerId={id}
Method: POST
Description: Proceeds with payment for an existing order.
Order id placed successfully
Order is also updated in delivery order API for the delivery agent to see all the order to be delivery with the expected delivery date and time
Also product quantity is decreased from the inventory
Authentication: Required
Parameters:
orderId (path): The ID of the order.
id (path): The customer's ID.
Request Body:
{
    "paymentMethod" : "Wallet",
    "userId": "66ab331d97d120b6a67da9f3"
}


Response:
200 OK: Successfully moved to payment.
400 Bad Request: If the request is invalid.
404 Not Found: If the order or customer is not found.
500 Internal Server Error: If there is a server error.
Example
{
   "MatchedCount": 1,
   "ModifiedCount": 1,
   "UpsertedCount": 0,
   "UpsertedID": null
}
{
   "_id": "66c25a30512bf19c6acef710",
   "customerId": "66c22faf818ea8e8db645d6f",
   "productIds": [
       "66aa320d5270e6700c8c94d1"
   ],
   "productNames": [
       "Google Pixel"
   ],
   "orderQuantitys": [
       1
   ],
   "orderedDate": "2024-08-19T02:06:14.164837563+05:30",
   "deliveredDates": [
       [
           "2024-08-19T02:06:14.525692564+05:30",
           "2024-08-20T02:06:14.525693313+05:30"
       ]
   ],
   "status": [
       [
           "Order Confirmed"
       ]
   ],
   "shippingAddress": {
       "_id": "66c253c8522685bb6b7f75c6",
       "fullName": "Ravi Kiran",
       "mobileNumber": "1234567890",
       "pincode": "18000",
       "line1": "....,",
       "line2": "...",
       "landmark": "...",
       "city": "Bangalore",
       "state": "Karnataka",
       "country": "India",
       "isDefault": true
   },
   "paymentDetails": {
       "_id": "66c25b3e25767654f5c32739",
       "paymentMethod": "Wallet",
       "totalAmount": 32990,
       "userId": "66c22faf818ea8e8db645d6f",
       "cardNumber": "",
       "cardExpiryDate": "",
       "nameOnCard": "",
       "upiId": ""
   },
   "priceDetail": {
       "listPrice": 34900,
       "sellingPrice": 32990,
       "deliveryCharge": 0,
       "totalAmount": 32990
   },
   "totalAmount": 32990
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
    "paymentMethod" : "Wallet",
    "userId": "66ab331d97d120b6a67da9f3"
}
"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass
2
{
    "paymentMethod" : "Wallet",
    "userId": "66ab331d97d120b6a67da9f3"
}
"Not Enough Amount in your Wallet!!"
400 Bad Request
Pass
3
{
    "paymentMethod" : "UPI",
    “upiId”: “abc”, 
    "userId": "66ab331d97d120b6a67da9f3"
}
"Only Wallet Payment Available!!"
400 Bad Request
Pass


Replace Order Request
Endpoint: /replaceOrderRequest/orderId={orderId}&productId={productId}&customerId={id}
Method: PUT
Description: Requests a replacement for a product in an existing order.
This request is from the customer side
Product can only be replace if the product has replacement policy and product can be replaced within 7 days of delivery date
Customer can replace a product once in the lifetime
Status will be updated and the expected replacement date will also be updated
Replace order API will also be updated for the delivery agent to see all the replacement orders and track it with the expected date of replacement
If that product is not available in the inventory, so customer can’t be replace it, customer have to place a return request or talk to customer care
If the product is not damage, product quantity will also be increased in the inventory
Authentication: Required
Parameters:
orderId (path): The ID of the order.
productId (path): The ID of the product.
id (path): The customer's ID.
Request Body:
{
    "dontLikeDueToColorOrSize": true,
    "isDamage": false
}



Response:
200 OK: Replacement request successfully created.
400 Bad Request: If the request is invalid.
404 Not Found: If the order, product, or customer is not found.
500 Internal Server Error: If there is a server error.
Example
{
   "_id": "66c25cda25767654f5c32744",
   "customerId": "66c22faf818ea8e8db645d6f",
   "productIds": [
       "66aa320d5270e6700c8c94d1"
   ],
   "productNames": [
       "Google Pixel"
   ],
   "orderQuantitys": [
       1
   ],
   "orderedDate": "2024-08-18T20:43:31.621Z",
   "deliveredDates": [
       [
           "2024-08-18T20:43:32.099Z",
           "2024-08-18T20:45:30.97Z",
           "2024-08-19T02:16:52.783220891+05:30",
           "2024-08-22T02:16:52.783221904+05:30"
       ]
   ],
   "status": [
       [
           "Order Confirmed",
           "Order Delivered",
           "Replace Confirmed"
       ]
   ],
   "shippingAddress": {
       "_id": "66c253c8522685bb6b7f75c6",
       "fullName": "Ravi Kiran",
       "mobileNumber": "1234567890",
       "pincode": "18000",
       "line1": "....,",
       "line2": "...",
       "landmark": "...",
       "city": "Bangalore",
       "state": "Karnataka",
       "country": "India",
       "isDefault": true
   },
   "paymentDetails": {
       "_id": "66c25cf325767654f5c3274a",
       "paymentMethod": "Wallet",
       "totalAmount": 32990,
       "userId": "66c22faf818ea8e8db645d6f",
       "cardNumber": "",
       "cardExpiryDate": "",
       "nameOnCard": "",
       "upiId": ""
   },
   "priceDetail": {
       "listPrice": 34900,
       "sellingPrice": 32990,
       "deliveryCharge": 0,
       "totalAmount": 32990
   },
   "totalAmount": 32990
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
    "dontLikeDueToColorOrSize": true,
    "isDamage": false
}
"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass
2
{
    "dontLikeDueToColorOrSize": true,
    "isDamage": false
}
"No Replace Policy for this product!!"
400 Bad Request
Pass
3
{
    "dontLikeDueToColorOrSize": true,
    "isDamage": false
}
"Please Replace when Order Delivered or Only One time Replace Policy!!"
400 Bad Request
Pass
4
{
    "dontLikeDueToColorOrSize": true,
    "isDamage": false
}
"Product Unavailable, Please try for Return!!"
400 Bad Request
Pass
5
{
    "dontLikeDueToColorOrSize": true,
    "isDamage": false
}
"Replace Expires!!"
400 Bad Request
Pass


Return Order Request
Endpoint: /returnOrderRequest/orderId={orderId}&productId={productId}&customerId={id}
Method: PUT
Description: Requests a return for a product in an existing order.
This request is from the customer side
Product can only be returned if the product has return policy and product can be return within 7 days of delivery date or replacement delivery date
Status will be updated and the expected return date will also be updated
Return order API will also be updated for the delivery agent to see all the return orders and track it with the expected date of return
If the product is not damage, product quantity will also be increased in the inventory
Authentication: Required
Parameters:
orderId (path): The ID of the order.
productId (path): The ID of the product.
id (path): The customer's ID.
Request Body:
{
    "dontLikeDueToColorOrSize": true,
    "isDamage": false
}

Response:
200 OK: Return request successfully created.
400 Bad Request: If the request is invalid.
404 Not Found: If the order, product, or customer is not found.
500 Internal Server Error: If there is a server error.
Example
{
   "_id": "66c25cda25767654f5c32744",
   "customerId": "66c22faf818ea8e8db645d6f",
   "productIds": [
       "66aa320d5270e6700c8c94d1"
   ],
   "productNames": [
       "Google Pixel"
   ],
   "orderQuantitys": [
       1
   ],
   "orderedDate": "2024-08-18T20:43:31.621Z",
   "deliveredDates": [
       [
           "2024-08-18T20:43:32.099Z",
           "2024-08-18T20:45:30.97Z",
           "2024-08-18T20:46:52.783Z",
           "2024-08-18T20:48:11.403Z",
           "2024-08-19T02:22:23.565272434+05:30",
           "2024-08-20T02:22:23.565273173+05:30"
       ]
   ],
   "status": [
       [
           "Order Confirmed",
           "Order Delivered",
           "Replace Confirmed",
           "Replace Order Delivered",
           "Return Confirmed"
       ]
   ],
   "shippingAddress": {
       "_id": "66c253c8522685bb6b7f75c6",
       "fullName": "Ravi Kiran",
       "mobileNumber": "1234567890",
       "pincode": "18000",
       "line1": "....,",
       "line2": "...",
       "landmark": "...",
       "city": "Bangalore",
       "state": "Karnataka",
       "country": "India",
       "isDefault": true
   },
   "paymentDetails": {
       "_id": "66c25cf325767654f5c3274a",
       "paymentMethod": "Wallet",
       "totalAmount": 32990,
       "userId": "66c22faf818ea8e8db645d6f",
       "cardNumber": "",
       "cardExpiryDate": "",
       "nameOnCard": "",
       "upiId": ""
   },
   "priceDetail": {
       "listPrice": 34900,
       "sellingPrice": 32990,
       "deliveryCharge": 0,
       "totalAmount": 32990
   },
   "totalAmount": 32990
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
    "dontLikeDueToColorOrSize": true,
    "isDamage": false
}
"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass
2
{
    "dontLikeDueToColorOrSize": true,
    "isDamage": false
}
"No Return Policy for this product!!"
400 Bad Request
Pass
3
{
    "dontLikeDueToColorOrSize": true,
    "isDamage": false
}
"Return Expires!!"
400 Bad Request
Pass
4
{
    "dontLikeDueToColorOrSize": true,
    "isDamage": false
}
"Please Refund when Order Delivered!!"
400 Bad Request
Pass


Cancel Order Request
Endpoint: /cancelOrderRequest/orderId={orderId}&productId={productId}&customerId={id}
Method: PUT
Description: Cancels an order or a specific product in an existing order.
This request is from the customer side
Product can be canceled within 1 days of ordered date, replacement date or return date
Status will be updated and the canceled date will also be updated
Authentication: Required
Parameters:
orderId (path): The ID of the order.
productId (path): The ID of the product.
id (path): The customer's ID.
Response:
200 OK: Order or product successfully canceled.
400 Bad Request: If the request is invalid.
404 Not Found: If the order, product, or customer is not found.
500 Internal Server Error: If there is a server error.
Example
{
   "_id": "66c25a30512bf19c6acef710",
   "customerId": "66c22faf818ea8e8db645d6f",
   "productIds": [
       "66aa320d5270e6700c8c94d1"
   ],
   "productNames": [
       "Google Pixel"
   ],
   "orderQuantitys": [
       1
   ],
   "orderedDate": "2024-08-18T20:36:14.164Z",
   "deliveredDates": [
       [
           "2024-08-18T20:36:14.525Z",
           "2024-08-19T02:10:29.325995695+05:30"
       ]
   ],
   "status": [
       [
           "Order Confirmed",
           "Cancelled"
       ]
   ],
   "shippingAddress": {
       "_id": "66c253c8522685bb6b7f75c6",
       "fullName": "Ravi Kiran",
       "mobileNumber": "1234567890",
       "pincode": "18000",
       "line1": "....,",
       "line2": "...",
       "landmark": "...",
       "city": "Bangalore",
       "state": "Karnataka",
       "country": "India",
       "isDefault": true
   },
   "paymentDetails": {
       "_id": "66c25b3e25767654f5c32739",
       "paymentMethod": "Wallet",
       "totalAmount": 32990,
       "userId": "66c22faf818ea8e8db645d6f",
       "cardNumber": "",
       "cardExpiryDate": "",
       "nameOnCard": "",
       "upiId": ""
   },
   "priceDetail": {
       "listPrice": 34900,
       "sellingPrice": 32990,
       "deliveryCharge": 0,
       "totalAmount": 32990
   },
   "totalAmount": 32990
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1


"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass
2


"Cancel Expires!!"
400 Bad Request
Pass
3


"You can't Cancel an order after Delivered!!"
400 Bad Request
Pass



Change Delivery Date
Endpoint: /changeDeliveryDate/orderId={orderId}&productId={productId}&customerId={id}
Method: PUT
Description: Changes the delivery date for an existing order.
This feature is only available for the prime membership
Chance the delivery date by increasing the respected delivery date in order API as well as in delivery order API
Authentication: Required
Parameters:
orderId (path): The ID of the order.
productId (path): The ID of the product.
id (path): The customer's ID.
Request Body:
{
    "noOfDayToIncrease":1
}

Response:
200 OK: Delivery date successfully updated.
400 Bad Request: If the request is invalid.
404 Not Found: If the order, product, or customer is not found.
500 Internal Server Error: If there is a server error.
Example


{
   "MatchedCount": 1,
   "ModifiedCount": 1,
   "UpsertedCount": 0,
   "UpsertedID": null
}
{
   "_id": "66c25a30512bf19c6acef710",
   "customerId": "66c22faf818ea8e8db645d6f",
   "productIds": [
       "66aa320d5270e6700c8c94d1"
   ],
   "productNames": [
       "Google Pixel"
   ],
   "orderQuantitys": [
       1
   ],
   "orderedDate": "2024-08-19T02:06:14.164837563+05:30",
   "deliveredDates": [
       [
           "2024-08-19T02:06:14.525692564+05:30",
           "2024-08-21T02:06:14.525693313+05:30"
       ]
   ],
   "status": [
       [
           "Order Confirmed"
       ]
   ],
   "shippingAddress": {
       "_id": "66c253c8522685bb6b7f75c6",
       "fullName": "Ravi Kiran",
       "mobileNumber": "1234567890",
       "pincode": "18000",
       "line1": "....,",
       "line2": "...",
       "landmark": "...",
       "city": "Bangalore",
       "state": "Karnataka",
       "country": "India",
       "isDefault": true
   },
   "paymentDetails": {
       "_id": "66c25b3e25767654f5c32739",
       "paymentMethod": "Wallet",
       "totalAmount": 32990,
       "userId": "66c22faf818ea8e8db645d6f",
       "cardNumber": "",
       "cardExpiryDate": "",
       "nameOnCard": "",
       "upiId": ""
   },
   "priceDetail": {
       "listPrice": 34900,
       "sellingPrice": 32990,
       "deliveryCharge": 0,
       "totalAmount": 32990
   },
   "totalAmount": 32990
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
    "noOfDayToIncrease":1
}
"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass
2
{
    "noOfDayToIncrease":1
}
"You can't change the Delivery date after the expected delivery date!!"
400 Bad Request
Pass
3
{
    "noOfDayToIncrease":1
}
"You can't change the Delivery date after the order been delivered!!"
400 Bad Request
Pass
4
{
    "noOfDayToIncrease":1
}
"You are not a prime member, kindly buy prime membership to change the delivery date!!"
400 Bad Request
Pass



13. Delivery Partner Management
Get All Deliver Orders
Endpoint: /deliverOrders
Method: GET
Description: Retrieves all orders that are ready to be delivered by delivery partners.
Authentication: None
Response:
200 OK: A list of orders to be delivered.
500 Internal Server Error: If there is a server error.
Example
[
   {
       "_id": "66b5a10dbfdff59f57914ebb",
       "customerId": "66ab331d97d120b6a67da9f3",
       "productId": "66aa0d90c9440c136263d090",
       "orderId": "66b5a10cbfdff59f57914eb7",
       "customerName": "Shiva Gupta",
       "customerAddress": {
           "_id": "66b06aeb19b253de5b363873",
           "fullName": "Shiva Gupta",
           "mobileNumber": "8982511209",
           "pincode": "471515",
           "line1": "near old bus stand,",
           "line2": "In frotn f government hospital",
           "landmark": "RC brother",
           "city": "LavKush nagar",
           "state": "MP",
           "country": "India",
           "isDefault": false
       },
       "expectedDeliveryDate": "2024-08-12T04:54:37.7Z"
   },
   {
       "_id": "66b9b843e2932a90221014a9",
       "customerId": "66b4537d0208bbd82ee589aa",
       "productId": "66a9e09b63831c9fba7f84a2",
       "orderId": "66b9b7ffe2932a902210149f",
       "customerName": "Ram Gupta",
       "customerAddress": {
           "_id": "66b9a11ec046307c7abca648",
           "fullName": "Ram Gupta",
           "mobileNumber": "1234567890",
           "pincode": "560000",
           "line1": "....,",
           "line2": "...",
           "landmark": "...",
           "city": "Bangalore",
           "state": "Karnataka",
           "country": "India",
           "isDefault": false
       },
       "expectedDeliveryDate": "2024-08-16T07:22:42.679Z"
   }
]


Get All Replace Orders
Endpoint: /replaceOrders
Method: GET
Description: Retrieves all orders that are ready to be replaced by delivery partners.
Authentication: None
Response:
200 OK: A list of orders to be replaced.
500 Internal Server Error: If there is a server error.
Example
null
Get All Return Orders
Endpoint: /returnOrders
Method: GET
Description: Retrieves all orders that are ready to be returned by delivery partners.
Authentication: None
Response:
200 OK: A list of orders to be returned.
500 Internal Server Error: If there is a server error.
Example
null
Deliver Order by Delivery Partner
Endpoint: /deliverOrderByDeliveryPartner/orderId={orderId}&productId={productId}&customerId={customerId}
Method: POST
Description: Marks an order as delivered by the delivery partner.
Updated the status of the product along with the delivery date
Delete the order from the delivery order API
Also update the report of the seller’s product
Authentication: None
Parameters:
orderId (path): The ID of the order.
productId (path): The ID of the product.
customerId (path): The ID of the customer.
Response:
200 OK: Order marked as delivered.
400 Bad Request: If the request is invalid.
404 Not Found: If the order, product, or customer is not found.
500 Internal Server Error: If there is a server error.
Example
{
   "MatchedCount": 1,
   "ModifiedCount": 1,
   "UpsertedCount": 0,
   "UpsertedID": null
}
{
   "MatchedCount": 1,
   "ModifiedCount": 1,
   "UpsertedCount": 0,
   "UpsertedID": null
}
{
   "DeletedCount": 1
}
{
   "_id": "66c25cda25767654f5c32744",
   "customerId": "66c22faf818ea8e8db645d6f",
   "productIds": [
       "66aa320d5270e6700c8c94d1"
   ],
   "productNames": [
       "Google Pixel"
   ],
   "orderQuantitys": [
       1
   ],
   "orderedDate": "2024-08-18T20:43:31.621Z",
   "deliveredDates": [
       [
           "2024-08-18T20:43:32.099Z",
           "2024-08-19T02:15:30.970447879+05:30"
       ]
   ],
   "status": [
       [
           "Order Confirmed",
           "Order Delivered"
       ]
   ],
   "shippingAddress": {
       "_id": "66c253c8522685bb6b7f75c6",
       "fullName": "Ravi Kiran",
       "mobileNumber": "1234567890",
       "pincode": "18000",
       "line1": "....,",
       "line2": "...",
       "landmark": "...",
       "city": "Bangalore",
       "state": "Karnataka",
       "country": "India",
       "isDefault": true
   },
   "paymentDetails": {
       "_id": "66c25cf325767654f5c3274a",
       "paymentMethod": "Wallet",
       "totalAmount": 32990,
       "userId": "66c22faf818ea8e8db645d6f",
       "cardNumber": "",
       "cardExpiryDate": "",
       "nameOnCard": "",
       "upiId": ""
   },
   "priceDetail": {
       "listPrice": 34900,
       "sellingPrice": 32990,
       "deliveryCharge": 0,
       "totalAmount": 32990
   },
   "totalAmount": 32990
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1


"Order is already Delivered!!"
400 Bad Request
Pass


Replace Order by Delivery Partner
Endpoint: /replaceOrderByDeliveryPartner/orderId={orderId}&productId={productId}&customerId={customerId}
Method: POST
Description: Marks an order as replaced by the delivery partner.
Updated the status of the product along with the replace date
Delete the order from the replace order API
Authentication: None
Parameters:
orderId (path): The ID of the order.
productId (path): The ID of the product.
customerId (path): The ID of the customer.
Response:
200 OK: Order marked as replaced.
400 Bad Request: If the request is invalid.
404 Not Found: If the order, product, or customer is not found.
500 Internal Server Error: If there is a server error.
Example
{
   "MatchedCount": 1,
   "ModifiedCount": 1,
   "UpsertedCount": 0,
   "UpsertedID": null
}
{
   "MatchedCount": 1,
   "ModifiedCount": 1,
   "UpsertedCount": 0,
   "UpsertedID": null
}
{
   "DeletedCount": 1
}
{
   "_id": "66c25cda25767654f5c32744",
   "customerId": "66c22faf818ea8e8db645d6f",
   "productIds": [
       "66aa320d5270e6700c8c94d1"
   ],
   "productNames": [
       "Google Pixel"
   ],
   "orderQuantitys": [
       1
   ],
   "orderedDate": "2024-08-18T20:43:31.621Z",
   "deliveredDates": [
       [
           "2024-08-18T20:43:32.099Z",
           "2024-08-18T20:45:30.97Z",
           "2024-08-18T20:46:52.783Z",
           "2024-08-19T02:18:11.403380639+05:30"
       ]
   ],
   "status": [
       [
           "Order Confirmed",
           "Order Delivered",
           "Replace Confirmed",
           "Replace Order Delivered"
       ]
   ],
   "shippingAddress": {
       "_id": "66c253c8522685bb6b7f75c6",
       "fullName": "Ravi Kiran",
       "mobileNumber": "1234567890",
       "pincode": "18000",
       "line1": "....,",
       "line2": "...",
       "landmark": "...",
       "city": "Bangalore",
       "state": "Karnataka",
       "country": "India",
       "isDefault": true
   },
   "paymentDetails": {
       "_id": "66c25cf325767654f5c3274a",
       "paymentMethod": "Wallet",
       "totalAmount": 32990,
       "userId": "66c22faf818ea8e8db645d6f",
       "cardNumber": "",
       "cardExpiryDate": "",
       "nameOnCard": "",
       "upiId": ""
   },
   "priceDetail": {
       "listPrice": 34900,
       "sellingPrice": 32990,
       "deliveryCharge": 0,
       "totalAmount": 32990
   },
   "totalAmount": 32990
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1


"Replace Order is already Delivered!!"
400 Bad Request
Pass


Return Order by Delivery Partner
Endpoint: /returnOrderByDeliveryPartner/orderId={orderId}&productId={productId}&customerId={customerId}
Method: POST
Description: Marks an order as returned by the delivery partner.
Updated the status of the product along with the return date
Delete the order from the return order API
Authentication: None
Parameters:
orderId (path): The ID of the order.
productId (path): The ID of the product.
customerId (path): The ID of the customer.
Response:
200 OK: Order marked as returned.
400 Bad Request: If the request is invalid.
404 Not Found: If the order, product, or customer is not found.
500 Internal Server Error: If there is a server error.
Example
{
   "MatchedCount": 1,
   "ModifiedCount": 1,
   "UpsertedCount": 0,
   "UpsertedID": null
}
{
   "MatchedCount": 1,
   "ModifiedCount": 1,
   "UpsertedCount": 0,
   "UpsertedID": null
}
{
   "DeletedCount": 1
}
{
   "MatchedCount": 1,
   "ModifiedCount": 1,
   "UpsertedCount": 0,
   "UpsertedID": null
}
{
   "_id": "66c25cda25767654f5c32744",
   "customerId": "66c22faf818ea8e8db645d6f",
   "productIds": [
       "66aa320d5270e6700c8c94d1"
   ],
   "productNames": [
       "Google Pixel"
   ],
   "orderQuantitys": [
       1
   ],
   "orderedDate": "2024-08-18T20:43:31.621Z",
   "deliveredDates": [
       [
           "2024-08-18T20:43:32.099Z",
           "2024-08-18T20:45:30.97Z",
           "2024-08-18T20:46:52.783Z",
           "2024-08-18T20:48:11.403Z",
           "2024-08-18T20:52:23.565Z",
           "2024-08-19T02:24:18.566789554+05:30"
       ]
   ],
   "status": [
       [
           "Order Confirmed",
           "Order Delivered",
           "Replace Confirmed",
           "Replace Order Delivered",
           "Return Confirmed",
           "Refunded"
       ]
   ],
   "shippingAddress": {
       "_id": "66c253c8522685bb6b7f75c6",
       "fullName": "Ravi Kiran",
       "mobileNumber": "1234567890",
       "pincode": "18000",
       "line1": "....,",
       "line2": "...",
       "landmark": "...",
       "city": "Bangalore",
       "state": "Karnataka",
       "country": "India",
       "isDefault": true
   },
   "paymentDetails": {
       "_id": "66c25cf325767654f5c3274a",
       "paymentMethod": "Wallet",
       "totalAmount": 32990,
       "userId": "66c22faf818ea8e8db645d6f",
       "cardNumber": "",
       "cardExpiryDate": "",
       "nameOnCard": "",
       "upiId": ""
   },
   "priceDetail": {
       "listPrice": 34900,
       "sellingPrice": 32990,
       "deliveryCharge": 0,
       "totalAmount": 32990
   },
   "totalAmount": 32990
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1


"Return Order is already Taken or Only one time you can return the product!!"
400 Bad Request
Pass



14. Review Management
Get All Reviews
Endpoint: /reviews
Method: GET
Description: Retrieves all reviews in the system.
Authentication: None
Response:
200 OK: A list of all reviews.
500 Internal Server Error: If there is a server error.
Example
[
   {
       "_id": "66b471f8b76055321d3c2697",
       "productId": "66aa0d90c9440c136263d090",
       "customerId": "66ab331d97d120b6a67da9f3",
       "customerName": "Shiva",
       "rating": 3,
       "headline": "It is not good.",
       "description": "I dont like it, but its good product.",
       "reviewDate": "2024-08-08T07:21:28.552Z",
       "reviewImages": []
   },
   {
       "_id": "66b472a9cb9832ff04db4189",
       "productId": "66aa0d90c9440c136263d090",
       "customerId": "66ab331d97d120b6a67da9f3",
       "customerName": "Shiva",
       "rating": 4,
       "headline": "It is not good.",
       "description": "I dont like it, but its good product.",
       "reviewDate": "2024-08-08T07:24:25.242Z",
       "reviewImages": []
   }
]


Get Reviews by Product ID
Endpoint: /reviews/productId={productId}
Method: GET
Description: Retrieves all reviews for a specific product.
Authentication: None
Parameters:
productId (path): The ID of the product.
Response:
200 OK: A list of reviews for the specified product.
404 Not Found: If the product is not found.
500 Internal Server Error: If there is a server error.
Example
[
   {
       "_id": "66b471f8b76055321d3c2697",
       "productId": "66aa0d90c9440c136263d090",
       "customerId": "66ab331d97d120b6a67da9f3",
       "customerName": "Shiva",
       "rating": 3,
       "headline": "It is not good.",
       "description": "I dont like it, but its good product.",
       "reviewDate": "2024-08-08T07:21:28.552Z",
       "reviewImages": []
   },
   {
       "_id": "66b472a9cb9832ff04db4189",
       "productId": "66aa0d90c9440c136263d090",
       "customerId": "66ab331d97d120b6a67da9f3",
       "customerName": "Shiva",
       "rating": 4,
       "headline": "It is not good.",
       "description": "I dont like it, but its good product.",
       "reviewDate": "2024-08-08T07:24:25.242Z",
       "reviewImages": []
   }
]


Get Most Recent Reviews by Product ID
Endpoint: /mostRecentReviews/productId={productId}
Method: GET
Description: Retrieves the most recent reviews for a specific product.
Authentication: None
Parameters:
productId (path): The ID of the product.
Response:
200 OK: A list of the most recent reviews for the specified product.
404 Not Found: If the product is not found.
500 Internal Server Error: If there is a server error.
Example
[
   {
       "_id": "66b472a9cb9832ff04db4189",
       "productId": "66aa0d90c9440c136263d090",
       "customerId": "66ab331d97d120b6a67da9f3",
       "customerName": "Shiva",
       "rating": 4,
       "headline": "It is not good.",
       "description": "I dont like it, but its good product.",
       "reviewDate": "2024-08-08T07:24:25.242Z",
       "reviewImages": []
   },
   {
       "_id": "66b471f8b76055321d3c2697",
       "productId": "66aa0d90c9440c136263d090",
       "customerId": "66ab331d97d120b6a67da9f3",
       "customerName": "Shiva",
       "rating": 3,
       "headline": "It is not good.",
       "description": "I dont like it, but its good product.",
       "reviewDate": "2024-08-08T07:21:28.552Z",
       "reviewImages": []
   }
]


Get Top Reviews by Product ID
Endpoint: /topReviews/productId={productId}
Method: GET
Description: Retrieves the top reviews (highest rated) for a specific product.
Authentication: None
Parameters:
productId (path): The ID of the product.
Response:
200 OK: A list of the top reviews for the specified product.
404 Not Found: If the product is not found.
500 Internal Server Error: If there is a server error.
Example
[
   {
       "_id": "66b472a9cb9832ff04db4189",
       "productId": "66aa0d90c9440c136263d090",
       "customerId": "66ab331d97d120b6a67da9f3",
       "customerName": "Shiva",
       "rating": 4,
       "headline": "It is not good.",
       "description": "I dont like it, but its good product.",
       "reviewDate": "2024-08-08T07:24:25.242Z",
       "reviewImages": []
   },
   {
       "_id": "66b471f8b76055321d3c2697",
       "productId": "66aa0d90c9440c136263d090",
       "customerId": "66ab331d97d120b6a67da9f3",
       "customerName": "Shiva",
       "rating": 3,
       "headline": "It is not good.",
       "description": "I dont like it, but its good product.",
       "reviewDate": "2024-08-08T07:21:28.552Z",
       "reviewImages": []
   }
]


Add Review
Endpoint: /addReview/productId={productId}&customerId={id}
Method: POST
Description: Adds a review for a specific product by a customer.
Authentication: Required
Parameters:
productId (path): The ID of the product.
id (path): The ID of the customer.
Request Body:
{
   "customerName": "Shiva",
   "rating" : 4,
   "headline" : "It is not good.",
   "description" : "I don't like it, but it's good product.",
   "reviewImages": []
}


Response:
201 Created: Review successfully added.
400 Bad Request: If the request is invalid.
404 Not Found: If the product or customer is not found.
500 Internal Server Error: If there is a server error.
Example
{
   "MatchedCount": 1,
   "ModifiedCount": 1,
   "UpsertedCount": 0,
   "UpsertedID": null
}
{
   "InsertedID": "66c261b00ed8dcd74d10e141"
}
{
   "_id": "66c261b00ed8dcd74d10e141",
   "productId": "66aa320d5270e6700c8c94d1",
   "customerId": "66c22faf818ea8e8db645d6f",
   "customerName": "Ravi",
   "rating": 3,
   "headline": "It is not good.",
   "description": "...",
   "reviewDate": "2024-08-19T02:33:44.876637671+05:30",
   "reviewImages": []
}



15.Report Management
View Reports by Seller ID
Endpoint: /report/sellerId={id}
Method: GET
Description: Retrieves all reports associated with a specific seller.
Report can be based on for last 1 year (monthly bases) for each products
Or for 10 years (yearly bases) for each products
Authentication: Required
Parameters:
id (path): The seller's ID.
Request Body:
{
   "statsBasedOn": "year"
}


Response:
200 OK: A list of reports for the specified seller.
404 Not Found: If the seller is not found.
500 Internal Server Error: If there is a server error.
Example
Stats Based On Year For the Last One Decade!!
{
   "2014": 0,
   "2015": 0,
   "2016": 0,
   "2017": 0,
   "2018": 0,
   "2019": 0,
   "2020": 0,
   "2021": 0,
   "2022": 0,
   "2023": 0,
   "2024": 2,
   "avergaePrice": 80900,
   "productImage": [
       "/assets/photos/iphone16.png"
   ],
   "productName": "iPhone 16",
   "totalUnitSold": 2
}
{
   "2014": 0,
   "2015": 0,
   "2016": 0,
   "2017": 0,
   "2018": 0,
   "2019": 0,
   "2020": 0,
   "2021": 0,
   "2022": 0,
   "2023": 0,
   "2024": 3,
   "avergaePrice": 12999,
   "productImage": [
       "/assets/photos/redmi.png"
   ],
   "productName": "Redmi",
   "totalUnitSold": 3
}
{
   "2014": 0,
   "2015": 0,
   "2016": 0,
   "2017": 0,
   "2018": 0,
   "2019": 0,
   "2020": 0,
   "2021": 0,
   "2022": 0,
   "2023": 0,
   "2024": 1,
   "avergaePrice": 32990,
   "productImage": [
       "/assets/photos/google.png"
   ],
   "productName": "Google Pixel",
   "totalUnitSold": 1
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
   "statsBasedOn": "year"
}
"Unauthorized User!!"
(if seller is not logged in)
401 Unauthorized
Pass



View Report by Product ID and Seller ID
Endpoint: /report/productId={productId}&sellerId={id}
Method: GET
Description: Retrieves a specific report for a product associated with a seller.
Report can be based on for last 1 year (monthly bases)
Or for 10 years (yearly bases)
Authentication: Required
Parameters:
productId (path): The product ID.
id (path): The seller's ID.
Request Body:
{
   "statsBasedOn": "month"
}


Response:
200 OK: The report for the specified product and seller.
404 Not Found: If the product or seller is not found.
500 Internal Server Error: If there is a server error.
Example
Stats Based On Month For the Last One Year for Google Pixel!!
{
   "productName": "Google Pixel",
   "productImage": [
       "/assets/photos/google.png"
   ],
   "totalUnitSold": 1,
   "avergaePrice": 32990,
   "january": 0,
   "february": 0,
   "march": 0,
   "april": 0,
   "may": 0,
   "june": 0,
   "july": 0,
   "august": 1,
   "september": 0,
   "october": 0,
   "november": 0,
   "december": 0
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
   "statsBasedOn": "month"
}
"Unauthorized User!!"
(if seller is not logged in)
401 Unauthorized
Pass



16. Search History Management
Get All Search Histories
Endpoint: /searchHistory
Method: GET
Description: Retrieves all search histories in the system for all the customers
Authentication: None
Response:
200 OK: A list of all search histories.
500 Internal Server Error: If there is a server error.
Example
[
   {
       "_id": "66b5e0d8c370303caccfa5aa",
       "customerId": "66ab331d97d120b6a67da9f3",
       "searchText": [
           "iphone"
       ]
   },
   {
       "_id": "66b5e11fc370303caccfa5ab",
       "customerId": "66b4537d0208bbd82ee589aa",
       "searchText": []
   },
   {
       "_id": "66c222f856fe9801831cd742",
       "customerId": "66c222ed56fe9801831cd73c",
       "searchText": []
   },
   {
       "_id": "66c22835818ea8e8db645d5c",
       "customerId": "66c22829818ea8e8db645d55",
       "searchText": []
   },
   {
       "_id": "66c22fb0818ea8e8db645d75",
       "customerId": "66c22faf818ea8e8db645d6f",
       "searchText": [
           ""
       ]
   }
]


Get Search History by Customer ID
Endpoint: /searchHistory/customerId={id}
Method: GET
Description: Retrieves the search history for a specific customer by their ID.
Authentication: Required
Parameters:
id (path): The customer's ID.
Response:
200 OK: The search history for the specified customer.
404 Not Found: If the customer is not found.
500 Internal Server Error: If there is a server error.
Example
[
   "iphone"
]

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1


"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass



17. Customer Care Management
Chat with Customer Care
Endpoint: /customerCare/message/customerId={id}
Method: POST
Description: Allows a customer to send a message to customer care.
Authentication: Required
Parameters:
id (path): The customer's ID.
Request Body:
{
   "message": ""
}


Response:
200 OK: The message was successfully sent to customer care.
400 Bad Request: If the message format is invalid.
500 Internal Server Error: If there is a server error.
Example


Hi! It's Amazon's messaging assistant again.
Product Name: Redmi
Is this what you need help with?
Choose Anyone of the following...
No, Something else
Yes, that's correct

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1
{
   "message": ""
}
"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass




18. Product Recommendation Management
Get Recommended Products by Customer ID
Endpoint: /productRecommendation/customerId={id}
Method: GET
Description: Retrieves product recommendations for a specific customer.
10 products related to last search based on category stored in product recommendation API
Authentication: Required
Parameters:
id (path): The customer's ID.
Response:
200 OK: A list of recommended products for the customer.
404 Not Found: If no recommendations are found for the customer.
500 Internal Server Error: If there is a server error.
Example
{
   "_id": "66bae957efbf2f19ef428380",
   "customerId": "66ab331d97d120b6a67da9f3",
   "productIds": [
       "66a8d9c2697de17b2285d258",
       "66a9e09b63831c9fba7f84a2",
       "66a9e0ec63831c9fba7f84a5",
       "66aa0d90c9440c136263d090",
       "66aa320d5270e6700c8c94d1",
       "66ab273b03aa2deac83093e2",
       "66b499c97dd2d413f5399ee9",
       "66a8d9c2697de17b2285d258",
       "66a9e09b63831c9fba7f84a2",
       "66a9e0ec63831c9fba7f84a5"
   ],
   "productName": [
       "IPhone 15",
       "iPhone 16",
       "Samsung",
       "Redmi",
       "Google Pixel",
       "Google Pixel",
       "Sony i7 5G",
       "IPhone 15",
       "iPhone 16",
       "Samsung"
   ],
   "aboutProduct": [
       "DYNAMIC ISLAND COMES TO IPHONE 15 -- Dynamic Island bubbles up alerts and Live Activities -- so you don't miss them while you're doing something else. You can see who's calling, track your next ride, check your flight status, and so much more.",
       "DYNAMIC ISLAND COMES TO IPHONE 15 -- Dynamic Island bubbles up alerts and Live Activities -- so you don't miss them while you're doing something else. You can see who's calling, track your next ride, check your flight status, and so much more.",
       "DYNAMIC ISLAND COMES TO samsung -- Dynamic Island bubbles up alerts and Live Activities -- so you don't miss them while you're doing something else. You can see who's calling, track your next ride, check your flight status, and so much more.",
       "DYNAMIC ISLAND COMES TO redmi -- Dynamic Island bubbles up alerts and Live Activities -- so you don't miss them while you're doing something else. You can see who's calling, track your next ride, check your flight status, and so much more.",
       "DYNAMIC ISLAND COMES TO google pixel -- Dynamic Island bubbles up alerts and Live Activities -- so you don't miss them while you're doing something else. You can see who's calling, track your next ride, check your flight status, and so much more.",
       "DYNAMIC ISLAND COMES TO google pixel 9A -- Dynamic Island bubbles up alerts and Live Activities -- so you don't miss them while you're doing something else. You can see who's calling, track your next ride, check your flight status, and so much more.",
       "DYNAMIC ISLAND COMES TO sony pixel 9A -- Dynamic Island bubbles up alerts and Live Activities -- so you don't miss them while you're doing something else. You can see who's calling, track your next ride, check your flight status, and so much more.",
       "DYNAMIC ISLAND COMES TO IPHONE 15 -- Dynamic Island bubbles up alerts and Live Activities -- so you don't miss them while you're doing something else. You can see who's calling, track your next ride, check your flight status, and so much more.",
       "DYNAMIC ISLAND COMES TO IPHONE 15 -- Dynamic Island bubbles up alerts and Live Activities -- so you don't miss them while you're doing something else. You can see who's calling, track your next ride, check your flight status, and so much more.",
       "DYNAMIC ISLAND COMES TO samsung -- Dynamic Island bubbles up alerts and Live Activities -- so you don't miss them while you're doing something else. You can see who's calling, track your next ride, check your flight status, and so much more."
   ],
   "brand": [
       "Apple",
       "Apple",
       "Samsung",
       "Redmi",
       "Google",
       "Google",
       "Sony",
       "Apple",
       "Apple",
       "Samsung"
   ],
   "bestSeller": [
       false,
       true,
       false,
       false,
       true,
       true,
       true,
       false,
       true,
       false
   ],
   "newRelease": [
       false,
       true,
       true,
       true,
       true,
       true,
       true,
       false,
       true,
       true
   ],
   "maxRetailPrice": [
       79900,
       99900,
       150000,
       15000,
       34900,
       44900,
       15900,
       79900,
       99900,
       150000
   ],
   "sellingPrice": [
       70900,
       80900,
       139000,
       12999,
       32990,
       42990,
       14990,
       70900,
       80900,
       139000
   ],
   "discount": [
       11.26,
       19.02,
       7.33,
       13.34,
       5.47,
       4.25,
       5.72,
       11.26,
       19.02,
       7.33
   ],
   "averageRating": [
       0,
       0,
       0,
       3.5,
       0,
       0,
       0,
       0,
       0,
       0
   ],
   "productImage": [
       "/assets/photos/iphone15.png",
       "/assets/photos/iphone16.png",
       "/assets/photos/samsumg.png",
       "/assets/photos/redmi.png",
       "/assets/photos/google.png",
       "/assets/photos/google.png",
       "assets/photos/sony.png",
       "/assets/photos/iphone15.png",
       "/assets/photos/iphone16.png",
       "/assets/photos/samsumg.png"
   ],
   "numberOfReviews": [
       0,
       0,
       0,
       2,
       0,
       0,
       0,
       0,
       0,
       0
   ]
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1


"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass




19. Recently Viewed Products Management
Get Recently Viewed Products by Customer ID
Endpoint: /recentlyViewedProduct/customerId={id}
Method: GET
Description: Retrieves the list of products recently viewed by a specific customer.
Last search product stored in recently viewed product API
Recently viewed product API stores 10 recent product details 
Authentication: Required
Parameters:
id (path): The customer's ID.
Response:
200 OK: A list of recently viewed products for the customer.
404 Not Found: If no recently viewed products are found for the customer.
500 Internal Server Error: If there is a server error.
Example


{
   "_id": "66bafa18efbf2f19ef428382",
   "customerId": "66ab331d97d120b6a67da9f3",
   "productIds": [
       "66b499c97dd2d413f5399ee9"
   ],
   "productName": [
       "Sony i7 5G"
   ],
   "aboutProduct": [
       "DYNAMIC ISLAND COMES TO sony pixel 9A -- Dynamic Island bubbles up alerts and Live Activities -- so you don't miss them while you're doing something else. You can see who's calling, track your next ride, check your flight status, and so much more."
   ],
   "brand": [
       "Sony"
   ],
   "bestSeller": [
       true
   ],
   "newRelease": [
       true
   ],
   "maxRetailPrice": [
       15900
   ],
   "sellingPrice": [
       14990
   ],
   "discount": [
       5.72
   ],
   "averageRating": [
       0
   ],
   "productImage": [
       "assets/photos/sony.png"
   ],
   "numberOfReviews": [
       0
   ]
}

Unit Testing
S.No.
Request
Response
Status
Pass/Fail
1


"Unauthorized User!!"
(if customer is not logged in)
401 Unauthorized
Pass




Thank you
