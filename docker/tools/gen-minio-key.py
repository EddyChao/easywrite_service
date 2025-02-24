import random
import string

def generate_random_key(length=20):
    characters = string.ascii_letters + string.digits
    return ''.join(random.choice(characters) for _ in range(length))

def generate_minio_keys():
    minio_access_key = generate_random_key(20)  # Ensure access key is at least 3 characters
    minio_secret_key = generate_random_key(40)  # Ensure secret key is at least 8 characters
    return f"- MINIO_ACCESS_KEY_FILE={minio_access_key}", f"- MINIO_SECRET_KEY_FILE={minio_secret_key}"

if __name__ == "__main__":
    access_key, secret_key = generate_minio_keys()
    print(access_key)
    print(secret_key)

# Windows: python .\gen_minio_key.py
# Other: python3 .\gen_minio_key.py
# Docker: docker build -t minio-key-generator . ; docker run --rm minio-key-generator
