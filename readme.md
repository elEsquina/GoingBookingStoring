# Online Bookstore Project

This project implements a RESTful API for an online bookstore that allows users to manage books, authors, and sales reports. The system includes a backend for handling database operations, middleware for authentication and logging, and endpoints for interacting with the data.

---

## Features

### Core Functionality:

- **Authentication**:

  - Authentication with the use of UUID Bearer tokens.

- **Book Management**:

  - Add, update, retrieve, and delete books.
  - Support for filtering books by title, author, or genre.

- **Author Management**:

  - Add, update, retrieve, and delete authors.

- **Sales Reporting**:

  - Generate daily sales reports, including total revenue and top-selling books.
  - Save reports as JSON files in `output-reports`.

### Middlewares and Security

- **Authentication**:
    - Token-based authentication for securing endpoints.
    - Make a `GET` request to `http://baseurl:8080/login` to obtain a Bearer token, which can be then attached to the `Authorization` header in any future request.
- **Request Logging**:
    - Logs all requests, including timestamps, methods, and response statuses, to `requests.log`.
- **Context Middleware**:
    - The `ContextGeneration` middleware adds a database connection (`DBTemplate`) to the request context, ensuring that each request has access to a shared database template.
    - Allows handlers to access the database without directly passing it through function arguments.
    - Ensures a timeout of 5 seconds for each request, preventing long-running queries from blocking resources.


### Scalability:

- **Database-Driven Architecture**:
  - Uses PostgreSQL for efficient data handling and storage.
- **Layered Design**:
  - Separates API, data, and middleware layers for maintainability and scalability.


### Graceful Shutdown:

The server shuts down cleanly upon receiving forced termination (ctrl+c).
Waits for ongoing requests to complete within a defined timeout period (5 seconds) before closing connections.

---

## Folder Structure

```
.
├── api                     # API handlers and routing
│   ├── authorHandler.go    # Handlers for author-related operations
│   ├── bookHandler.go      # Handlers for book-related operations
│   ├── middleWares.go      # Middleware for logging and authentication
├── configs                 # Configuration files
├── data                    # Database and data access logic
│   ├── dbTemplate.go       # Database interaction template
│   ├── reportGeneration.go # Logic for generating sales reports
│   ├── DAOFactory.go       # Ensures the existence of one single instance of each repository
│   ├── IDAO.go             # Abstract generic DAO interface
│   ├── structs.go          # Entities layer
│   ├── *DAO.go             # Concrete repositories
├── docs                    # Documentation
├── output-reports          # Directory for saved sales reports
├── sql                     # SQL scripts for schema and migrations
├── tests                   # Tests
├── main.go                 # Entry point of the application
├── requests.log            # Log file for HTTP requests
```

---

## API Endpoints

### Base URL

```
http://localhost:8080
```

### Login 

| Endpoint      | Method | Description                                                                  |
| ------------- | ------ | ---------------------------------------------------------------------------- |
| `/login`      | Any    | Returns an authorization token to use as header to access server ressources  |

### Books

| Endpoint      | Method | Description                          |
| ------------- | ------ | ------------------------------------ |
| `/books`      | GET    | List all books or filter by criteria |
| `/books`      | POST   | Add a new book                       |
| `/books/{id}` | GET    | Retrieve book details by ID          |
| `/books/{id}` | PUT    | Update a book by ID                  |
| `/books/{id}` | DELETE | Delete a book by ID                  |

### Authors

| Endpoint        | Method | Description                   |
| --------------- | ------ | ----------------------------- |
| `/authors`      | GET    | List all authors              |
| `/authors`      | POST   | Add a new author              |
| `/authors/{id}` | GET    | Retrieve author details by ID |
| `/authors/{id}` | PUT    | Update an author by ID        |
| `/authors/{id}` | DELETE | Delete an author by ID        |

---

## Design Reasoning

### Why Separate Layers?

- **API Layer**:
  - Handles HTTP requests and routes them to the appropriate business logic.
  - Provides decoupling between user interactions and core logic.

- **Middleware**:
  - Improves security and maintainability by centralizing common tasks like authentication and logging.

- **Data Layer**:
  - Abstracts database interactions, enabling easier testing and flexibility for future database migrations.

### Why Use `DBTemplate`?

- Ensures consistency in database queries by centralizing logic for `SELECT`, `INSERT`, `UPDATE`, and `DELETE`.
- Simplifies error handling and improves maintainability.

### Why Token-Based Authentication?

- Stateless.
- Easily integrated with middleware to secure endpoints.

### Why Save Reports to JSON?

- JSON is a widely accepted format that can be processed by various tools and languages.
- Allows for easy integration with external systems for analytics or backups.

---

## How to Run

1. **Set up the database**:

   - Ensure PostgreSQL is installed and running.
   - Use the SQL scripts in the `sql` folder to set up the schema.

2. **Start the application**:

   ```bash
   go run main.go
   ```

3. **Access the API**:

   - Base URL: `http://localhost:8080`

4. **Generate reports**:

   - Reports are automatically generated every 24 hours (24 seconds for testing) and saved in the `output-reports` directory.

---

## Logging

- **Request Logs**:
  - Stored in `requests.log`.
  - Includes timestamps, HTTP methods, and response statuses.
- **Error Logs**:
  - Errors during report generation or database operations are logged to the console.


## Testing

- **Data testing**:
  - In the tests/dummydata folder, wrote a python script that generates a custom number of rows (faker library for dummy data).
  - This is to test wether the backend works as expected and test for robustness test when dealing with large amount of data.

- **Stress testing**:
  - Wrote a python script (tests/stresstest) that requests all the GET endpoints to test the responsivity of the application under load.