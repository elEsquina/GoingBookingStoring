import requests
import random
import time

BASE_URL = "http://localhost:8080"
AUTH_TOKEN = "a9395c9d-d089-11ef-b45e-98e743997220"

def perform_get_request(endpoint):
    url = f"{BASE_URL}{endpoint}"
    headers = {"Authorization": f"Bearer {AUTH_TOKEN}"}
    print(f"Requesting: {url}")
    response = requests.get(url, headers=headers)
    print(f"Status Code: {response.status_code}")

# Test books/{id} from 1 to 100
for book_id in range(1, 101):
    perform_get_request(f"/books/{book_id}")
    # Occasionally request /books
    if random.randint(0, 9) == 0:
        perform_get_request("/books")
    time.sleep(0.1) 

# Test authors/{id} from 1 to 100
for author_id in range(1, 101):
    perform_get_request(f"/authors/{author_id}")
    # Occasionally request /authors
    if random.randint(0, 9) == 0:
        perform_get_request("/authors")
    time.sleep(0.1) 

print("Testing complete.")
