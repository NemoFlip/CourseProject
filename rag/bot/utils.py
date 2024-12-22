from collections import defaultdict
from time import time


class SpamDetector:
    def __init__(self):
        self.suspicious_users = {}
        self.message_count = defaultdict(list)
        self.blocked_users = {}

    def is_blocked(self, user_id, current_time):
        return user_id in self.blocked_users and current_time < self.blocked_users[user_id]

    def detect_spam(self, user_id, current_time):
        self.message_count[user_id].append(current_time)
        self.message_count[user_id] = [
            t for t in self.message_count[user_id] if current_time - t <= 5
        ]

        if len(self.message_count[user_id]) > 5:
            self.blocked_users[user_id] = current_time + 10
            return True
        return False

    def is_suspicious(self, user_id):
        return user_id in self.suspicious_users

    def verify_captcha(self, user_id, answer):
        return str(self.suspicious_users.get(user_id)) == answer.strip()

    def update_message_count(self, user_id, current_time):
        self.message_count[user_id].append(current_time)

    def add_captcha(self, user_id, captcha_answer):
        self.suspicious_users[user_id] = captcha_answer

    def remove_captcha(self, user_id):
        if user_id in self.suspicious_users:
            del self.suspicious_users[user_id]

    def detect_suspicious_activity(self, message):
        user_id = message.from_user.id
        if user_id in self.suspicious_users:
            return True
        if message.text == "/start":
            return True
        if message.content_type == "sticker":
            return True
        return False
