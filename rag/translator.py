from deep_translator import GoogleTranslator

class Translator:
    @staticmethod
    def translate(text, src_lang="auto", dest_lang="en"):
        return GoogleTranslator(source=src_lang, target=dest_lang).translate(text)
