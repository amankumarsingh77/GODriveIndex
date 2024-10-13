<div align="center">
  <!-- <br />
    <a href="https://youtu.be/lEflo_sc82g?feature=shared" target="_blank">
      <img src="https://github.com/JavaScript-Mastery-Pro/medicare-dev/assets/151519281/160a9367-29e8-4e63-ae78-29476b04bff3" alt="Project Banner">
    </a>
  <br /> -->

  <div>
    <img src="https://img.shields.io/badge/Go-00ADD8?logo=Go&logoColor=white&style=for-the-badge" alt="golang" />
    <img src="https://img.shields.io/badge/Google%20drive-grey?style=for-the-badge&logo=googledrive" alt="google drive" />
  </div>

  <h3 align="center">GODriveIndex</h3>

   <div align="center">
     A google drive go API to index, search and download files from google drive without download quota limits.
    </div>

    ⚠️ Warning: The quota free download feature is currently still under development. I have not tested it for large number of downloads. Feel free to create an issue if you face any bugs.
</div>

## 📋 <a name="table">Table of Contents</a>

1.  [Introduction](#introduction)
2.  [Tech Stack](#tech-stack)
3.  [Features](#features)
4.  [Quick Start](#quick-start)
5.  [Connect With Me](#connect-with-me)

## <a name="introduction">🤖 Introduction</a>

GoDriveIndex is a powerful open-source Google Drive indexing solution built with Go. This project aims to provide a seamless way to index, search, and download files from Google Drive without being constrained by download quota limits. By leveraging the Google Drive API and Go's robust concurrency features, GoDriveIndex offers a fast and reliable method to interact with your Google Drive contents.
Whether you're managing a large collection of files or need a reliable way to access and share Google Drive contents, GoDriveIndex provides a robust solution that combines the speed of Go with the versatility of Google Drive.


## <a name="tech-stack">⚙️ Tech Stack</a>

- Golang
- Google Drive API

1. Efficient Indexing: Quickly catalog your entire Google Drive structure.
2. Advanced Search: Easily find files and folders using various search parameters.
3. Quota-Free Downloads: Bypass Google Drive's download quotas for smoother file retrieval.
4. Secure Authentication: Utilize OAuth 2.0 for safe and authorized access to Google Drive.
5. User-Friendly Interface: A clean and intuitive web interface for easy navigation and file management.

## <a name="features">🔋 Features</a>

👉 **Efficient Indexing**: Quickly catalog your entire Google Drive structure.

👉 **Advanced Search**: Easily find files and folders using various search parameters.

👉 **Quota-Free Downloads**: Bypass Google Drive's download quotas for smoother file retrieval.

👉 **Secure Authentication**: Utilize OAuth 2.0 for safe and authorized access to Google Drive.


NOTE: MANY FEATURES ARE STILL UNDER DEVELOPMENT.

## <a name="quick-start">🤸 Quick Start</a>

Follow these steps to set up the project locally on your machine.

**Prerequisites**
## Quick Start

### Prerequisites
- **Go**: Ensure you have Go installed on your machine. You can download it from [here](https://golang.org/dl/).
- **Git**: Make sure Git is installed. You can download it from [here](https://git-scm.com/).
- **Environment Variables**: Create a `.env` file in the root directory with necessary configurations.

### Setup

1. **Clone the Repository**
   ```bash
   git clone https://github.com/amankumarsingh77/GODriveIndex.git
   cd GODriveIndex
   ```

2. **Install Dependencies**
   - Navigate to the project directory and run:
   ```bash
   go mod tidy
   ```

3. **Configure Environment Variables**
   - Create a `.env` file in the root directory and add your configurations. Example:
   ```plaintext
   GOOGLE_CLIENT_ID=YOUR_GOOGLE_CLIENT_ID
   GOOGLE_CLIENT_SECRET=YOUR_GOOGLE_CLIENT_SECRET
   OAUTH_REDIRECT_URL=YOUR_OAUTH_REDIRECT_URL
   ```

***Additional Configuration***
- **Setup Service Account**: Ensure `serviceAccount.json` is properly configured is if you are using service account for authentication. Follow [this](https://cloud.google.com/iam/docs/service-account-overview) guide to create a service account  and use the `serviceAccount.demo.json` file as reference.
- **Authentication**: Modify `auth.json` as needed for your authentication setup. Follow  `auth.demo.json` for reference.

**Run the Server**
   - Start the server using:
   ```bash
   go run cmd/server/main.go
   ```

**Access the Application**
   - Open your browser and go to `http://localhost:8080` to access the application.

## Connect With Me

Q) Want to chat or need assistance with setting up a project?

A) You can connect with me on [X](https://x.com/amankumar404) and [Gmail](mailto:amankumarsingh7702@gmail.com)

