
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

export function errorAlertMessage(error) {
  let message = "Ooops! Something went wrong.";

  switch (error) {
    case Errors.ConversationListingError:
      return message = "Something went wrong while we tried to show you your conversations.";

    case Errors.StartingConversationError:
      return message = "Something went wrong while we tried to start a new conversation.";

    case Errors.QuestionError:
      return message = "Something went wrong while we tried to answer your question.";

    default:
  }

  return message;
}
