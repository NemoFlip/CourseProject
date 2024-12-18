from deep_translator import GoogleTranslator

class Translator:
    """Класс для перевода текста с использованием GoogleTranslator."""
    @staticmethod
    def translate(text, src_lang="auto", dest_lang="en"):
        """Перевод текста с одного языка на другой."""
        return GoogleTranslator(source=src_lang, target=dest_lang).translate(text)
