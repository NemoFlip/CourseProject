import random
import pytest
import requests

BASE_URL = "http://host.docker.internal:8080"

@pytest.fixture(scope="module")
def user_data():
    unique_id = random.randint(1, 10000)
    username = f'user_{unique_id}'
    password = f'password_{unique_id}'
    email = f"user{unique_id}@gmail.com"
    phone = "89861113245"
    return {'username': username, 'email': email, 'phone': phone, 'password': password}

def test_register_user(user_data):
    register_url = f"{BASE_URL}/registration"

    response = requests.post(register_url, json=user_data)

    assert response.status_code == 201

    return