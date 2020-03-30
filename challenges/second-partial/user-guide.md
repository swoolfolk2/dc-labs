User Guide
==========

Start the Server
----------------
First you need to start the server. To do this, run de main.go prograam.
To run it in the terminal, go to the directory and run it with the commnand
$ go run main.go

This will run the server in localhost with the port 8080.

Use API
---------------
First, you need to open another terminal. Once done this, there are 4 Post 
Petitions you can use: Login, Logout, Status and Upload.

-Login-

To login, you need to enter the following command:

$ curl -u <username>:<password> http://localhost:8080/login

If you entered correctly, the following message will appear:

{
	"message": "Hi <username>, welcome to the DPIP System",
	"token": <ACCESS_TOKEN>
}

This will save you username and password in an internal database, 
remembering your information.

-Status-

If you wanna check the server status, you have to provide 
your given access-token and enter the following command:

$ curl -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/status

If the token exists and there are no errors with the server, 
then the following message will appear:

{
	"message": "Hi username, the DPIP System is Up and Running"
	"time": "2015-03-07 11:06:39"
}

-Upload-

If you are logged in and want to upload an image, the 
following command should work:

$ curl -F ‘data=@path’ -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/upload

If the provided image and access-token are correct, and there are no errors,
the following message should appear:

{
	"message": "An image has been successfully uploaded",
	"filename": "image.png",
	"size": "500kb"
}

-Logout-

If there are no other actions you want to do, you can logout providing your access-token:

$ curl -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/logout

This will delete your ACCESS_TOKEN and info from the internal database and display the
message:

{
	"message": "Bye username, your token has been revoked"
}


