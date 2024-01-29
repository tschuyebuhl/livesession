# Database cache library 

## Running the tests

```bash
go test ./service -race 
```
## Running the application 
```bash
docker compose up -d
```
Add `--build` flag to rebuild the image after making changes to the code.

# Code documentation
## Package service

### Overview
The `UserService` package provides a service for managing user data in a Go application. It leverages caching to improve performance and reduce unnecessary data fetching. This service is particularly useful for applications that require efficient user data retrieval and management.

The package consists of the `UserService` struct, one exported method, and one helper method.

### `UserService` Struct
This struct is the core of the package and includes the following fields:

- `repo`: An instance of `data.Repository` for reading data from the database.
- `cache`: An instance of `cache.Cache` for caching user data.
- `pending`: A map tracking pending requests for user data.
- `mu`: A mutex to ensure thread-safety.

#### Method: `GetUser`
Retrieves a user by ID.

##### Parameters:
- `id`: The unique identifier of the user.

##### Returns:
- `*data.User`: The user object if found.
- `bool`: A flag indicating whether the user was found in the cache.
- `error`: Error object, if any.

##### Description:
1. Checks if the user exists in the cache. If found, returns the user.
2. Implements double-checked locking to handle concurrent requests:
    - If another goroutine has updated the cache, it fetches the user from the cache.
    - Otherwise, it marks the request as pending.
3. Fetches the user from the repository if not in cache.
4. Updates the cache with the fetched user data.
5. Notifies all waiters (goroutines waiting for this data).

#### Method: `notifyWaiters`
Notifies all goroutines waiting for a specific user data.

##### Parameters:
- `user`: The user data to be sent to the waiting goroutines.
- `id`: The unique identifier of the user.

##### Description:
1. Locks the mutex to ensure thread safety.
2. Sends the user data to all goroutines in the pending list for the specified user ID.
3. Removes the user ID from the pending list.

## Usage example
```go
// Initialize repository and cache instances 
repo := MyRepoImplementation{}
cache := MyCacheImplementation{}

// Create a UserService instance
userService := service.NewUserService(repo, cache)

// Get a user by ID
user, found, err := userService.GetUser(myUserID)
if err != nil {
    log.Fatal(err)
}
if found {
    fmt.Println("User found:", user)
} else {
    fmt.Println("User not found in cache, fetched from repository")
}
```
The `Repository` and `Cache` interfaces are defined as follows:
```go
type Repository interface {
	Get(ID) (*User, error)
}
```
```go
type Cache interface {
	Get(key data.ID) (*data.User, bool, bool)
	Put(value *data.User)
	Delete(key data.ID)
	Nuke(sure bool)
}

```