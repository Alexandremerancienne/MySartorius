# MySartorius: REST API to track workout progression
MySartorius is a Golang RESTful API designed for a gym center. The API enables its clients and coaches to keep track of their workout progression.

# Authentication

* Access is granted to authenticated users via JSON Web Tokens (JWTs).

# Models
* Users are divided into 4 categories: Manager, Coach, Client, Anonymous User (potential future member of the gym).
* Coaches and clients have 2 forms of interaction:
  - Sessions: live sessions at the gym center;
  - Tasks: exercises to complete at home.
* Clients can create for themselves reminders.
* Reminders are comments added by a client to a task to provide additional detail.

# Permissions
* Users are divided into three categories: Manager, Coaches, Clients;
* The Manager can read and edit any information available except the reminders posted by the clients for themselves;
* Other users can access any information regarding their activity: Coaches and Clients can access all the sessions and tasks they are involved in;
* Clients can also access and edit all the reminders they have created for themselves;
* Only Coaches and the Manager can edit sessions and tasks;
* A coach cannot edit a session or task involving another coach;

# Installation  

Database Setup: MySQL.

This locally-executable API can be installed and executed from using the following steps.
1.	Clone this repository: `git clone https://github.com/Alexandremerancienne/MySartorius.git` or download the code [as a zip file](https://github.com/Alexandremerancienne/MySartorius/archive/refs/heads/main.zip).
2.	Once locally downloaded, install required modules: `go mod tidy`.
3.	Run the server: `go run main.go`.  

# Usage and detailed endpoint documentation

One you have launched the server, you can consume the API at the following endpoints: 
* Login: [http://localhost:4000/login](http://localhost:4000/login)
* Access and edit users: [http://localhost:4000/users](http://localhost:4000/users)
* Access and edit sessions: [http://localhost:4000/sessions](http://localhost:4000/sessions)
* Access and edit tasks: [http://localhost:4000/tasks](http://localhost:4000/tasks)
* Access and edit reminders: [http://localhost:4000/tasks/{task_id}/reminders](http://localhost:4000/tasks/{task_id}/reminders)

All these endpoints support HTTP requests using GET, POST, PUT and DELETE methods.

Examples:

* Delete user with ID 1 (DELETE method): [http://localhost:4000/users/1](http://localhost:4000/users/1)
* Update reminder 1 of task 4 (PUT method): [http://localhost:4000/tasks/4/reminders/1](http://localhost:4000/tasks/4/reminders/1)

# Filters
You can apply filters to search a specific piece of information.
## Search and filter users
You can search and filter users with the following endpoint: http://localhost:4000/users/. The filters available are:
* `role=<string>` to get users by role (coach or client). Example: [http://localhost:4000/users?role=coach](http://localhost:4000/users?role=coach)
* `first_name=<string>` to search users by first name. The search does an exact match of the first name and is independent of character case. Example: [http://localhost:4000/users?first_name=john](http://localhost:4000/users?first_name=john)
* `last_name=<string>` to search users by last name. The search does an exact match of the last name and is independent of character case. Example: [http://localhost:4000/users?last_name=doe](http://localhost:4000/users?last_name=doe)
* `first_name_contains=<string>` to search users whose first name contains the search term. The search is independent of character case. Example: [http://localhost:4000/users?first_name_contains=jo](http://localhost:4000/users?first_name_contains=jo)
* `last_name_contains=<string>` to search users whose last name contains the search term. The search is independent of character case. Example: [http://localhost:4000/users?last_name_contains=do](http://localhost:4000/users?last_name_contains=do)

These filters can be combined, for example:

* Filter users by role (coach) with first name equal to "John" and last name containing "do": [http://localhost:4000/users?role=coach&first_name=john&last_name_contains=do)](http://localhost:4000/users?role=coach&first_name=john&last_name_contains=do)

## Search and filter sessions
You can search and filter sessions with the following endpoint: http://localhost:4000/sessions. The filters available are:
* `client_id=<integer>` to search sessions involving client with specific ID. Example: [http://localhost:4000/sessions?client_id=2](http://localhost:4000/sessions?client_id=2)
* `coach_id=<integer>` to search sessions involving coach with specific ID. Example: [http://localhost:4000/sessions?coach_id=5](http://localhost:4000/sessions?client_id=5)

These filters can be combined, for example:

* Filter sessions involving coach with ID 5 and client with ID 2: [http://localhost:4000/sessions?client_id=2&coach_id=5](http://localhost:4000/sessions?client_id=2&client_id=5)

## Search and filter tasks
You can search and filter tasks with the following endpoint: http://localhost:4000/tasks. The filters available are:
* `client_id=<integer>` to search tasks involving client with specific ID. Example: [http://localhost:4000/tasks?client_id=2](http://localhost:4000/tasks?client_id=2)
* `coach_id=<integer>` to search tasks involving coach with specific ID. Example: [http://localhost:4000/tasks?coach_id=5](http://localhost:4000/tasks?client_id=5)

These filters can be combined, for example:

* Filter tasks involving coach with ID 5 and client with ID 2: [http://localhost:4000/tasks?client_id=2&coach_id=5](http://localhost:4000/sessions?client_id=2&client_id=5)
