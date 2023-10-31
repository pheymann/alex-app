
import { Language } from './language';

const GermanTranslation = {
  header: {
    home: 'Home',
    signOut: 'Abmelden',
  },
  signIn: {
    emailPlaceholder: 'Email',
    passwordPlaceholder: 'Passwort',
    button: 'Anmelden',
  },
  home: {
    exploreArtButton: 'Kunst entdecken',
  },
  conversation: {
    initialField: 'Erzähl mir etwas über:',
    artContextPrompt: {
      title: 'Erzähl mir etwas über: ',
      placeholder: 'Die Mona Lisa von Leonardo da Vinci',
      field: 'Erzähl mir etwas über:'
    },
    questionPrompt: {
      placeholder: 'Hast du eine Frage?'
    }
  },
  prompt: {
    send: 'Senden',
  },
  error: {
    unknown: 'Ooops! Etwas ist schief gelaufen.',
    conversationListing: 'Etwas ist schief gelaufen, als wir versucht haben, deine Unterhaltungen zu laden.',
    startingConversation: 'Etwas ist schief gelaufen, als wir versucht haben, eine neue Unterhaltung zu starten.',
    question: 'Etwas ist schief gelaufen, als wir versucht haben, deine Frage zu beantworten.',
    signIn: 'Entweder ist deine Email oder das Passwort falsch.',
  },
};

const EnglishTranslation = {
  header: {
    home: 'Home',
    signOut: 'Sign Out',
  },
  signIn: {
    emailPlaceholder: 'Email',
    passwordPlaceholder: 'Password',
    button: 'Sign In',
  },
  home: {
    exploreArtButton: 'Explore Art',
  },
  conversation: {
    initialField: 'Tell me something about',
    artContextPrompt: {
      title: 'Tell me something about: ',
      placeholder: 'The Mona Lisa by Leonardo da Vinci',
      field: 'Tell me something about'
    },
    questionPrompt: {
      placeholder: 'Do you have any questions?'
    }
  },
  prompt: {
    send: 'Send',
  },
  error: {
    unknown: 'Ooops! Something went wrong.',
    conversationListing: 'Something went wrong while we tried to load your conversations.',
    startingConversation: 'Something went wrong while we tried to start a new conversation.',
    question: 'Something went wrong while we tried to answer your question.',
    signIn: 'Either your email or password is invalid.',
  },
};

export const Translation = initTranslations();

function initTranslations() {
  const map = new Map();
  map.set(Language.German, GermanTranslation);
  map.set(Language.English, EnglishTranslation);

  return map;
}
