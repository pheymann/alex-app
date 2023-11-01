
import { Language } from './language';

const GermanTranslation = {
  header: {
    home: 'Alex',
    signOut: 'Abmelden',
  },
  signIn: {
    emailPlaceholder: 'Email',
    passwordPlaceholder: 'Passwort',
    button: 'Anmelden',
    newPassword: 'Neues Passwort',
    confirmPassword: 'Passwort bestätigen',
    changePasswordButton: 'Ändern',
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
    loading: 'Ich denke nach ... das kann bis zu 30s dauern.',
  },
  error: {
    unknown: 'Ooops! Etwas ist schief gelaufen.',
    conversationListing: 'Etwas ist schief gelaufen, als wir versucht haben, deine Unterhaltungen zu laden.',
    startingConversation: 'Etwas ist schief gelaufen, als wir versucht haben, eine neue Unterhaltung zu starten.',
    question: 'Etwas ist schief gelaufen, als wir versucht haben, deine Frage zu beantworten.',
    signIn: 'Entweder ist deine Email oder das Passwort falsch.',
    completeSignUp: 'Etwas ist schief gelaufen, als wir versucht haben, deine Anmeldung abzuschließen.'
  },
};

const EnglishTranslation = {
  header: {
    home: 'Alex',
    signOut: 'Sign Out',
  },
  signIn: {
    emailPlaceholder: 'Email',
    passwordPlaceholder: 'Password',
    button: 'Sign In',
    newPassword: 'New Password',
    confirmPassword: 'Confirm Password',
    changePasswordButton: 'Change',
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
    loading: 'Thinking ... that may take up to 30s.',
  },
  error: {
    unknown: 'Ooops! Something went wrong.',
    conversationListing: 'Something went wrong while we tried to load your conversations.',
    startingConversation: 'Something went wrong while we tried to start a new conversation.',
    question: 'Something went wrong while we tried to answer your question.',
    signIn: 'Either your email or password is invalid.',
    completeSignUp: 'Something went wrong while we tried to complete your sign up.'
  },
};

export const Translation = initTranslations();

function initTranslations() {
  const map = new Map();
  map.set(Language.German, GermanTranslation);
  map.set(Language.English, EnglishTranslation);

  return map;
}
