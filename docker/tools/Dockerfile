# Use a Python base image
FROM python:slim

# Set the working directory
WORKDIR /app

# Copy the script into the container
COPY gen-minio-key.py .

# Run the script
CMD ["python", "gen-minio-key.py"]
