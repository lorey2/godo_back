This project is a go todolist backend with 
    Gin framework for handeling http requests and
    GORM for Postgres database

First it calls database.Connect() located in database/db.go to connect to pg db
the . operator is like :: c in this case

gin.Default is a router that listen for incoming request

The code configures CORS to allow browsers to make requests from a different origin
    (like your frontend running on a different port).

There is two type of route
    /register and /login that are public and visible to anyone
    /[everything tast related] that you need to be connected to have access
            handlers.AuthMiddleware() to check that




models.go store the blueprint of a task and a user with tags usefull later

db.go handle the connection to the db
    DB.AutoMigrate is cool: it create the database table automagically or update them

auth.go
    Register
        BindJSON fill the model.User with the info and check if json is correct
            If there is an error we send a json with 404 and the error
        Then we hash the password and create the database user
            If not created error
        Default task for cool onboarding created too
        FINALLY we send to gin.context json : Httpcreated, Registration successful
    Login
        on remplit un model.user avec la request
            si pas bon erreur
        on remplit un autre model.user avec la db depuis le mail
        on compare les passwd
            si cest bon on cree un token jwt (json web token)
        we send statusOk+token
    authmiddleware
        we get the authorization header and parse it
        with the token we can get the userID
            there is some bulshit notation but isok like Anonymous Function
