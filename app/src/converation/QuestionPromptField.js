import { useState } from "react";
import { PromptField } from "./PromptField";

export default function QuestionPromptField({
  conversation,
  setConversation,
  awsContext,
}) {
  const [question, setQuestion] = useState('');

  const handleQuestion = () => {
    if (question === '') {
      // TODO: show error message
      console.error('missing user question');
      return;
    }

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

    fetch(`/api/conversation/${conversation.id}/continue`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${awsContext.token}`,
      },
      body: JSON.stringify({
        question: question,
      }),
    })
      .then(response => response.json())
      .then(data => {
        // removing loading message
        continuedConversation.messages.pop();

        const responseConversation = {
          ...continuedConversation,
          messages: [
            ...continuedConversation.messages,
            {
              role: data.role,
              text: data.text,
              speechClipUuid: data.speechClipUuid,
              speechClipUrl: data.speechClipUrl,
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
        console.log(error);
      });
  };

  return (
    <PromptField  value={ question }
                  onChangeValue={ setQuestion }
                  onSubmit={ () => handleQuestion() }
                  placeholder='What is on your mind?'>
    </PromptField>
  );
}
