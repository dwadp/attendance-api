# Attendance API
Attendance API is a REST API to provide a way to manage employee's attendance, shifts and more. It is intended to be a coding challenge test. This project is built using:
* Go (v1.22.2)
* PostgreSQL

### How to run
#### Build the binary
First of all, make sure you have a PostgreSQL database server running. And make sure that you have installed Go (v1.22.2) or higher on your local machine to build this project. Once you're ready, clone this repository and then run:

```sh
go build -o bin/attendance
```

Then lookup into the folder `bin/` of your current directory, there will be an executable file called `attendance`. Depending on your operating system but if you're on windows it will output an `attendance.exe` file which you can run directly.
The binary provides you some ability which is:
* To provision your database & run your migrations
* To start the API server

#### Configuration
There is a `config.example.yml` file in this repository which you can copy and paste on your local machine as the main configuration of the API. Before you move further, please adjust the configuration based on your local machine ***especially*** the database (`PostgreSQl`).

**IMPORTANT**
> The binary default configuration will be at your `$HOME/.attendance-api/confg.yml`. Although you can customize it by providing a `--config` flag every time you run the binary
> which point to an absolute path to where you put the configuration file.


#### Run migrations
Before you can run the migrations, please setup your [configuration](https://github.com/dwadp/attendance-api/new/main?filename=README.md#configuration) first. To run the migration use this command:
```sh
./bin/attendance migration up
```

#### Run the server
After all of the above steps successfully finished, you can run the API server to start using the API by running this command on your terminal:
```sh
./bin/attendance
```

You can test the API by using the postman collection that i've provided [here](https://github.com/dwadp/attendance-api/blob/main/postman/Attendances.postman_collection.json).

### Notes
If you have a problem running this project on your local machine, feel free to contact me at: [dewaadiperdana@gmail.com](mailto:dewaadiperdana@gmail.com).
