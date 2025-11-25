# Docker Deployment Steps for Assignment 03

Follow these steps in order to complete the assignment.

## Step 1: Build the Docker Image

Open PowerShell in the project directory and run:

```powershell
docker build -t chatroom-server .
```

This will:
- Use the Dockerfile to build the image
- Compile the Go server application
- Create a minimal Alpine-based image

## Step 2: Test the Docker Image Locally

1. **Run the container:**
   ```powershell
   docker run -d -p 8080:8080 --name chatroom-test chatroom-server
   ```

2. **Check if it's running:**
   ```powershell
   docker ps
   ```
   
   You should see the `chatroom-test` container listed.

3. **View the logs to confirm the server started:**
   ```powershell
   docker logs chatroom-test
   ```
   
   You should see: "Chatroom server started on port 8080"

4. **Test with the client:**
   Open a new PowerShell window and run:
   ```powershell
   go run client.go
   ```
   
   Try sending a few messages to verify the server works.

5. **Stop the test container:**
   ```powershell
   docker stop chatroom-test
   docker rm chatroom-test
   ```

## Step 3: Push to Docker Hub

1. **Create a Docker Hub account** (if you don't have one):
   - Go to https://hub.docker.com
   - Sign up for a free account

2. **Login to Docker Hub from PowerShell:**
   ```powershell
   docker login
   ```
   
   Enter your Docker Hub username and password.

3. **Tag your image with your Docker Hub username:**
   ```powershell
   docker tag chatroom-server YOUR_DOCKERHUB_USERNAME/chatroom-server:latest
   ```
   
   Replace `YOUR_DOCKERHUB_USERNAME` with your actual Docker Hub username.

4. **Push the image to Docker Hub:**
   ```powershell
   docker push YOUR_DOCKERHUB_USERNAME/chatroom-server:latest
   ```
   
   This will upload your image to Docker Hub.

5. **Verify the image is public:**
   - Go to https://hub.docker.com
   - Login and check your repositories
   - Find `chatroom-server` and ensure it's set to **Public**
   - If it's private, click on the repository → Settings → Make Public

## Step 4: Update README.md

Replace all instances of `YOUR_DOCKERHUB_USERNAME` in the README.md file with your actual Docker Hub username.

For example, if your username is `johndoe`, change:
- `YOUR_DOCKERHUB_USERNAME/chatroom-server` → `johndoe/chatroom-server`

## Step 5: Push to GitHub

1. **Initialize Git repository** (if not already done):
   ```powershell
   git init
   git add .
   git commit -m "Add Docker deployment for chatroom server"
   ```

2. **Create a new repository on GitHub:**
   - Go to https://github.com
   - Click "New repository"
   - Name it (e.g., "go-chatroom-docker")
   - Do NOT initialize with README (you already have one)
   - Click "Create repository"

3. **Push to GitHub:**
   ```powershell
   git remote add origin https://github.com/YOUR_GITHUB_USERNAME/YOUR_REPO_NAME.git
   git branch -M main
   git push -u origin main
   ```

## Step 6: Verify Your Submission

Before submitting, verify:

✅ GitHub repository contains:
   - server.go
   - client.go
   - shared/types.go
   - go.mod
   - Dockerfile
   - .dockerignore
   - README.md (with your Docker Hub link)

✅ Docker Hub image:
   - Is publicly accessible
   - Can be pulled: `docker pull YOUR_DOCKERHUB_USERNAME/chatroom-server`

✅ README.md includes:
   - Direct link to Docker Hub image
   - Instructions for building, running, and testing
   - Project structure with Docker files

## Step 7: Submit

Submit the GitHub repository link to your assignment portal.

---

## Troubleshooting

### Docker build fails
- Ensure Docker Desktop is running
- Check that go.mod is present
- Verify all source files are in the directory

### Cannot connect to server in Docker
- Ensure port 8080 is mapped: `-p 8080:8080`
- Check if container is running: `docker ps`
- View logs: `docker logs chatroom-test`

### Docker push fails
- Make sure you're logged in: `docker login`
- Verify the image is tagged correctly
- Check your internet connection

### Client cannot connect
- Ensure Docker container is running
- Check port mapping is correct
- Verify firewall isn't blocking port 8080
