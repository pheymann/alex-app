import { useRef, useState } from "react";
import { PromptField } from "./PromptField";
import { logError, pushLogMessage } from "../logger";

export default function QuestionPromptField({
  conversation,
  setConversation,
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

    awsFetch.call(`/api/conversation/${conversation.id}/continue`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        question: question,
      }),
    })
      .then(rawData => {
        pushLogMessage(logEntriesRef, { level: 'debug', message: rawData });

        const json = JSON.parse(rawData);

        // removing loading message
        continuedConversation.messages.pop();

        const responseConversation = {
          ...continuedConversation,
          messages: [
            ...continuedConversation.messages,
            {
              role: json.role,
              text: json.text,
              speechClipUuid: json.speechClipUuid,
              speechClipUrl: json.speechClipUrl,
            },
            {
              role: 'prompt-user-question',
            }
          ],
        };

        setConversation(responseConversation);
        setQuestion('');
      })
      .catch(error => {
        logError({ awsFetch, error, logEntriesRef: logEntriesRef});
        alert('Error continuing conversation:\n' + error);
      });
  };

  return (
    <PromptField  value={ question }
                  onChangeValue={ setQuestion }
                  onSubmit={ () => handleQuestion() }
                  placeholder='Do you have any questions?'
                  maxLength={ 500 } >
    </PromptField>
  );
}
