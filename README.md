# Secure File Upload and Sharing App

This project is a Go-based backend application that allows users to upload files without authentication. Each file is associated with a secret word or hash, which can be shared with other users for secure access to the file. Files are stored in AWS S3.

![Screenshot from 2024-07-13 22-04-03](https://github.com/user-attachments/assets/e58d80a2-7581-4e54-a80e-b703484454f2)
## Features

- Upload files without authentication
- Secure file sharing through secret words or hashes (currently overwritter files with same hash)
- Files stored securely in AWS S3

## Prerequisites

- Go (version 1.16 or later)
- AWS account with S3 bucket

## Installation

1. Clone the repository

   ```bash
   git clone https://github.com/mariovlv/chickenfile
   cd chickenfile
   ```

2. Install the dependencies

   ```bash
   go mod tidy
   ```

3. Create a `.env` file in the root of your project and add your AWS credentials and bucket name:

   ```bash
   PORT=PORT_NUMBER
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

2. Open your browser and navigate to `http://localhost:8080`

3. Feel ready to deploy in Heroku.
