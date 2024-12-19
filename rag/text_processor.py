import re

class TextProcessor:
    def __init__(self, input_file, output_file):
        self.input_file = input_file
        self.output_file = output_file
        self.ques_ans = {}
        self.content = ""

    def process_dataset(self):
        with open(self.input_file, 'r') as f:
            data = f.read()

        sections = re.split(r'(?<=\w):\s+', data)
        for i in range(0, len(sections), 2):
            if i + 1 < len(sections):
                question = sections[i].strip()
                answer = sections[i + 1].strip()
                self.ques_ans[question] = answer

        for key, value in self.ques_ans.items():
            self.content += f"{key}\n{value}\n\n"

        with open(self.output_file, 'w') as f:
            f.write(self.content)

        return self.content
