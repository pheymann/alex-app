import { useRef, useState } from "react";
import { PromptField } from "./PromptField";
import { logError, pushLogMessage } from "../logger";

export default function ArtContextPromptField({
  setConversation,
  awsFetch,
}) {
  const [artContext, setArtContext] = useState('');

  const logEntriesRef = useRef([]);

  const handleStartConversation = () => {
    if (artContext === '') {
      // TODO: show error message
      setConversation({
        messages: [{
          role: 'prompt-art-context',
        }],
      });
      console.error('missing art context');
      return;
    }

    const conversation = {
      messages: [
        {
          role: 'user',
          text: `Tell me something about ${artContext}`,
        },
        {
          role: 'loading',
        }
      ],
    };
    setConversation(conversation);

    awsFetch.call(`/api/conversation/create/art`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        artContext: artContext,
      }),
    })
      .then(rawData => {
        pushLogMessage(logEntriesRef, { level: 'debug', message: rawData });

        const json = JSON.parse(rawData);
        const responseConversation = {
          ...json,
          messages: [
            conversation.messages[0],
            ...json.messages,
            {
              role: 'prompt-user-question',
            },
          ],
        };
        setConversation(responseConversation);
      })
      .catch(error => {
        logError({ awsFetch, error, logEntriesRef: logEntriesRef});
        alert('Error starting conversation:\n' + error);
      });
  };

  return (
    <PromptField  value={ artContext }
                  onChangeValue={ setArtContext }
                  onSubmit={ () => handleStartConversation() }
                  placeholder='The Mona Lisa by Leonardo da Vinci'
                  maxLength={ 150 }>
      <div className="row">
        <div className='col'>
          <p>
            Tell me something about:
          </p>
        </div>
      </div>
    </PromptField>
  );
}
