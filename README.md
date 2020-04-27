# Dear colleagues,

#### This is the documentation for the API developed for the technical challenge.

##### Please note that due to not using a database for this exercise, ID of transaction is not an autoincrement PK, but a random generated number.
##### I decided not to validate against creation of this random number due to short ammount of time of completion.
##### In reality, i would choose a different approach (using a database that already solves PK generation and read/write Locks)
# Important 
## Database: 
#####Since persistence was not required, i took the liberty to create a .json file (named db.json) inside the path /src/api/database. 
#####This file will act as a DB to store transactions.

##API documentation:

###application folder:

#####url_mapping.go maps all of the endpoints for the possible requests.
- #####application.go builds app.

###controllers folder:
- #####rest_hanlder.go handles all of the HTTP requests (MVC arquitecture: controller)

###database folder: 
- ##### db.json is used for auxiliary storage.
- #####handler.go is used to R/W the auxiliary storage.

###models folder:
- #####models.go stores the structures used inside the api.

###services folder:
- #####service_provider.go handles all the logic of each process ensuring that all requirements are met.

###test folder:
- #####test.go is used to store unit test. Please keep in mind that due to low time of completion, not all functions in the api are tested.

###vendor folder:
#####This folder is used to store dependencies.

##Basic API usage:

####Compile and execute main.go. Afterwards, localhost:8080 should start listening for requests. Possible endpoints to be used are:

- #####GET http://localhost:8080/account : This endpoint is used to return current account balance.
- #####GET http://localhost:8080/transactions : This endpoint is used to return transaction history.
- #####GET http://localhost:8080/transactions/:id : This endpoint is used to return transaction searched for with id. Transaction ID must be provided in url instead of ":id"
- #####POST http://localhost:8080/transactions : This endpoint is used to generate a new transaction. Keep in mind that this requires a request body with following format

**Example:**

      {
        "type": "credit",
        "amount": 100
      }
  
  ###### Ammount value can be any float number
