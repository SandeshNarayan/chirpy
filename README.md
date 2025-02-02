Chirpy

Chirpy is a lightweight microblogging platform where users can post short messages called "chirps." It features authentication, user profile management, and chirp creation, retrieval, and deletion.

Features

User authentication (JWT-based login, refresh, and revoke tokens)

Create, retrieve, and delete chirps

User profile updates (email, password, and is_chirpy_red setting)

Admin metrics tracking

Secure API with role-based access control

Installation

Prerequisites

Go 1.18+

PostgreSQL

Git

Setup

# Clone the repository
git clone https://github.com/SandeshNarayan/chirpy.git
cd chirpy

# Create a .env file and configure environment variables
cp .env.example .env

# Install dependencies
go mod tidy

Database Setup

# Start PostgreSQL and create a database
psql -U postgres -c "CREATE DATABASE chirpy;"

# Apply migrations
psql -U postgres -d chirpy -f migrations.sql

Running the Server

go run main.go

Server runs on http://localhost:8080.

API Endpoints

User Authentication

POST /api/users – Create a new user

POST /api/login – Authenticate and get an access token

POST /api/refresh – Refresh access token

POST /api/revoke – Revoke token

PUT /api/users – Update user email and password

Chirps

POST /api/chirps – Create a new chirp

GET /api/chirps – Retrieve all chirps

GET /api/chirps/{chirpID} – Get a chirp by ID

DELETE /api/chirps/{chirpID} – Delete a chirp (only by the author)

Admin

GET /admin/metrics – View server metrics

POST /admin/reset – Reset the database

Error Handling

401 Unauthorized: Invalid or missing token

403 Forbidden: Unauthorized action

404 Not Found: Resource does not exist

500 Internal Server Error: Unexpected server failure

