# **Doc-Classification**

## **Setup**

### **Secret Files**
Ensure the following files are in the appropriate directory:

1. **`firebase_service.json`**  
   Contains credentials for authenticating server-side applications with Firebase services.

2. **`google_client_secret.json`**  
   A configuration file provided by Google that contains the credentials needed to authenticate with Google APIs.

3. **`.env`**  
   Contains environment variables required by the API. Include the following variables:
   - **`OPENAI_API_KEY`**: API key for OpenAI (e.g., ChatGPT).
   - **`GOOGLE_AUTH_FILE`**: Path to the Google API configuration file  
     (e.g., `google_client_secret.json`).
m 
---

## **How to Run**

### **1. Build the Docker Image**
Run the following command to build the Docker image:

```bash
docker build . --no-cache -t doc-classification
```

### **2. Run the service**
Use this command to run the containerized service:

```bash
docker run \
    -v $(pwd)/firebase_service.json:/firebase_service.json \
    -v $(pwd)/google_client_secret.json:/google_client_secret.json \
    -v $(pwd)/.env:/.env \
    -e FIREBASE_CONFIG=/app/firebase_service.json \
    -e CLIENT_SECRET=/app/client_secret_973692223612-28ae9a7njdsfh7gv89l0fih5q36jt52m \
    -p 8000:8000 doc-classification
```