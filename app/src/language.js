
export const Language = Object.freeze({
  German: Symbol(0),
  English: Symbol(1),
});

export const LanguageLocalStorageKey = "alex-user-language";

export function switchToLanguageFrom(language) {
  switch (language) {
    case Language.German:
      return Language.English;

    case Language.English:
      return Language.German;

    default:
      return Language.German;
  }
}

export function prettyPrintLanguage(language) {
  switch (language) {
    case Language.German:
      return "Deutsch";

    case Language.English:
      return "English";

    default:
      return "Unknown";
  }
}

export function encodeLanguage(language) {
  switch (language) {
    case Language.German:
      return "German";

    case Language.English:
      return "English";

    default:
      return "Unknown";
  }
}

export function decodeLanguage(language) {
  switch (language) {
    case "German":
      return Language.German;

    case "English":
      return Language.English;

    default:
      return Language.German;
  }
}
