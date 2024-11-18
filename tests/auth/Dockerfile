FROM python:3.9-slim

WORKDIR /auth_tests

RUN pip install requests pytest

COPY tests/auth/auth_tests.py .

CMD ["pytest", "auth_tests.py"]