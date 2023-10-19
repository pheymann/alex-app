
export const Errors = Object.freeze({
  UnknownError:              Symbol(0),
  ConversationListingError:   Symbol(1),
  StartingConversationError:  Symbol(2),
  QuestionError:              Symbol(3),
});

export function errorToCode(error) {
  switch (error) {
    case Errors.ConversationListingError:
      return 1;

    case Errors.StartingConversationError:
      return 2;

    case Errors.QuestionError:
      return 3;

    default:
      return 0;
  }
}

export function codeToError(errorCode) {
  switch (errorCode) {
    case 1:
    case "1":
      return Errors.ConversationListingError;

    case 2:
    case "2":
      return Errors.StartingConversationError;

    case 3:
    case "3":
      return Errors.QuestionError;

    default:
      return Errors.UnknownError;
  }
}

export function errorAlertMessage(error, i18n) {
  let message = i18n.error.unknown;

  switch (error) {
    case Errors.ConversationListingError:
      return message = i18n.error.conversationListing;

    case Errors.StartingConversationError:
      return message = i18n.error.startingConversation;

    case Errors.QuestionError:
      return message = i18n.error.question;

    default:
  }

  return message;
}
