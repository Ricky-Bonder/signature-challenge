# My Application to the Coding Challenge

This Signature Service is a program that allows users to manage 
signature devices and sign transactions using various algorithms.

## Prerequisites
Before running the Signature Service, 
ensure you have the Go (v1.16+) installed.

## Usage

Navigate to the project directory:

`cd fiskaly-coding-challenge/signing-service-challenge-go`

Run the code:

`go run main.go`

## API Endpoints
The Signature Service provides the following API endpoints:

- **GET** `/api/v0/health`: Check the health of the service.
- **POST** `/api/v0/create-signature-device` : Create a new signature device.
        
    Request Body Example:
    ```json
    {
      "algorithm": "RSA", // Accepted algorithms: RSA, ECC
      "label": "An Interesting Label"
    }
    ```
- **POST** `/api/v0/sign-transaction`: Sign a transaction using a signature device ID.
        
    Request Body Example:
    ```json
    {
      "id": "ef219680-af8a-4d7c-8949-5cf947a76c23", //previously created device's UUID
      "label": "Signed Transaction Label"
    }
    ```
- **GET** `/api/v0/get-signature-device?id=<device-UUID>`: Get information about a specific signature device.
    
- **GET** `/api/v0/get-all-devices`: Get information about all signature devices.

Everything was user tested on http://localhost:8080 through the Postman Agent.
## Testing
To run tests, use the following command:

`go test ./...`

## Further Implementation

While the current implementation provides basic functionality for managing signature devices 
and signing transactions, there are several areas where the code could be extended or improved:

1. #### Database Integration

   Currently, the Signature Service stores data in memory, but it could be easily substituted
with a relational Database through the Storage interface.

2. #### Signature Serialization

    Improvement of the serialization of the signature
in the response body of /api/v0/sign-transaction.

3. #### Unit and Integration Testing
    
    Expansion of the test suite to cover more edge cases and ensure comprehensive test coverage.
Implementing both unit tests and integration tests will help identify and prevent regressions 
as the codebase evolves.

4. #### Logging and Monitoring

    Implement proper logging and monitoring to track system behavior and performance. 
Log important events and errors to facilitate troubleshooting and debugging.