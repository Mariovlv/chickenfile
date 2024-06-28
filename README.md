# Secure File Upload and Sharing App

This project is a Go-based backend application that allows users to upload files without authentication. Each file is associated with a secret word or hash, which can be shared with other users for secure access to the file. Files are stored in AWS S3.

## Features

- Upload files without authentication
- Secure file sharing through secret words or hashes
- Files stored securely in AWS S3

## Prerequisites

Before you begin, ensure you have met the following requirements:

- Go (version 1.16 or later)
- AWS account with S3 bucket
- `go get` the following packages:
  - `github.com/aws/aws-sdk-go-v2/aws`
  - `github.com/aws/aws-sdk-go-v2/config`
  - `github.com/aws/aws-sdk-go-v2/service/s3`
  - `github.com/aws/aws-sdk-go-v2/feature/s3/manager`
  - `github.com/joho/godotenv`
  - `github.com/labstack/echo/v4`
  - `github.com/labstack/echo/v4/middleware`

## Installation

1. Clone the repository

   ```bash
   git clone https://github.com/mariovlv/chickenfile
   cd secure-file-upload
   ```

2. Install the dependencies

   ```bash
   go mod tidy
   ```

3. Create a `.env` file in the root of your project and add your AWS credentials and bucket name:
   ```
   AWS_ACCESS_KEY_ID=your-access-key-id
   AWS_SECRET_ACCESS_KEY=your-secret-access-key
   AWS_REGION=your-aws-region
   AWS_BUCKET_NAME=your-s3-bucket-name
   ```

## Usage

1. Run the application

   ```bash
   go run main.go
   ```

2. Open your browser and navigate to `http://localhost:1323`

3. Use the `/upload` endpoint to upload files. The response will include a secret word/hash that you can share with others to access the file.

## API Endpoints

### Upload File

- **URL:** `/upload`
- **Method:** `POST`
- **Form Data:**
  - `name` - Name of the uploader
  - `email` - Email of the uploader
  - `file` - File to be uploaded
- **Response:**
  - `200 OK` on success with a JSON object containing the hash
  - `400 Bad Request` if the file is too big or missing

## Example

### Uploading a File

Use `curl` to upload a file:

```bash
curl -F "name=John Doe" -F "email=johndoe@example.com" -F "file=@/path/to/your/file" http://localhost:1323/upload
```

### Notes

- **Security Considerations:** Ensure your AWS credentials are not hardcoded in the code. Use environment variables for better security practices.
- **Hash Generation:** You can use libraries such as `crypto/sha256` to generate unique hashes for the files.
- **S3 Permissions:** Make sure your S3 bucket permissions are set correctly to allow the app to upload and download files securely.

This README provides a basic structure to get started with your file upload and sharing app using Go and AWS S3. Feel free to customize it based on your project requirements.

## To dos

### Frontend

- [ ] Create a simple frontend using HTML, CSS, and JavaScript
  - [ ] Form for file upload with fields for name, email, and file
  - [ ] Display the generated hash/secret word after upload
  - [ ] Form for entering the hash/secret word to download the file

### Improve Error Handlers

- [ ] Enhance error handling in the backend
  - [ ] Return descriptive error messages for various error conditions (e.g., file too large, missing fields, S3 upload failure)
  - [ ] Implement proper HTTP status codes for different error scenarios

### Handling Download Routes

- [ ] Implement the `/download/:hash` endpoint
  - [ ] Retrieve the file from S3 using the provided hash
  - [ ] Return appropriate responses if the hash is invalid or the file is not found
  - [ ] Ensure secure access to files based on the hash/secret word

### Hashing and Server-Side Access to Files

- [ ] Generate a unique hash/secret word for each uploaded file
  - [ ] Use a secure hashing algorithm (e.g., SHA-256) to create the hash
  - [ ] Store the hash and associated file metadata in a secure database or in-memory store
- [ ] Implement server-side access controls to ensure files are only accessible via the generated hash/secret word
  - [ ] Validate the hash/secret word when accessing the file for download
  - [ ] Ensure the files are securely transmitted to the requesting user
