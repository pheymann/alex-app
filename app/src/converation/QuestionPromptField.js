import { useRef, useState } from "react";
import { PromptField } from "./PromptField";
import { logError, pushLogMessage } from "../logger";
import { Errors } from "../ErrorAlert";

export default function QuestionPromptField({
  conversation,
  setConversation,
  i18n,
  awsFetch,
}) {
  const [question, setQuestion] = useState('');

  const logEntriesRef = useRef([]);

  const handleQuestion = () => {
    if (question === '') {
      // TODO: show error message
      console.error('missing user question');
      return;
    }

    pushLogMessage(logEntriesRef, { level: 'debug', message: `question: ${question}` });

    // remove user prompt field
    conversation.messages.pop();

    const continuedConversation = {
      ...conversation,
      messages: [
        ...conversation.messages,
        {
          role: 'user',
          text: question,
        },
        {
          role: 'loading',
        }
      ],
    };
    setConversation(continuedConversation);

    const startProcessing = async () => {
      try {
        const continueResponse = await awsFetch.callResponse(`/api/conversation/${conversation.id}/continue`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            question: question,
          }),
        })

        if (continueResponse.ok) {
          await awsFetch.poll(
            `/api/conversation/${conversation.id}/poll`,
            {
              method: 'GET',
            },
            {
              handleSuccess: async (pollResponse, clearInterval) => {
                if (pollResponse.status === 200) {
                  clearInterval();

                  const message = await pollResponse.json();
                  pushLogMessage(logEntriesRef, { level: 'debug', message: message });

                  // removing loading message
                  continuedConversation.messages.pop();

                  const responseConversation = {
                    ...continuedConversation,
                    messages: [
                      ...continuedConversation.messages,
                      {
                        role: message.role,
                        text: message.text,
                        speechClipUuid: message.speechClipUuid,
                        speechClipUrl: message.speechClipUrl,
                      },
                      {
                        role: 'prompt-user-question',
                      }
                    ],
                  };

                  setConversation(responseConversation);
                  setQuestion('');
                }
              },
              handleError: (error) => {
                logError({ awsFetch, error, logEntriesRef: logEntriesRef});

                // removing loading message
                continuedConversation.messages.pop();

                const responseConversation = {
                  ...continuedConversation,
                  messages: [
                    ...continuedConversation.messages,
                    {
                      role: 'error',
                      errorCode: Errors.QuestionError,
                    },
                    {
                      role: 'prompt-user-question',
                    }
                  ],
                };

                setConversation(responseConversation);
                setQuestion('');
              },
            },
          );
        }
      } catch (error) {
        logError({ awsFetch, error, logEntriesRef: logEntriesRef});

        // removing loading message
        continuedConversation.messages.pop();

        const responseConversation = {
          ...continuedConversation,
          messages: [
            ...continuedConversation.messages,
            {
              role: 'error',
              errorCode: Errors.QuestionError,
            },
            {
              role: 'prompt-user-question',
            }
          ],
        };

        setConversation(responseConversation);
        setQuestion('');
    }
    };

    startProcessing();
  }

  return (
    <PromptField  value={ question }
                  onChangeValue={ setQuestion }
                  onSubmit={ () => handleQuestion() }
                  placeholder={ i18n.conversation.questionPrompt.placeholder }
                  maxLength={ 500 }
                  i18n={ i18n }
    >
    </PromptField>
  );
}
