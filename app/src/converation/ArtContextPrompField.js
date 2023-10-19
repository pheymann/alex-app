import { useRef, useState } from "react";
import { PromptField } from "./PromptField";
import { logError, pushLogMessage } from "../logger";
import { useNavigate } from "react-router-dom";
import { Errors, errorToCode } from "../ErrorAlert";

export default function ArtContextPromptField({
  setConversation,
  i18n,
  awsFetch,
}) {
  const [artContext, setArtContext] = useState('');
  const navigate = useNavigate();

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
          text: `${i18n.conversation.artContextPrompt.field} ${artContext}`,
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
        navigate('/?errorCode=' + errorToCode(Errors.StartingConversationError));
      });
  };

  return (
    <PromptField  value={ artContext }
                  onChangeValue={ setArtContext }
                  onSubmit={ () => handleStartConversation() }
                  placeholder={ i18n.conversation.artContextPrompt.placeholder }
                  maxLength={ 150 }
                  i18n={ i18n }
    >
      <div className="row">
        <div className='col'>
          <p>
            { i18n.conversation.artContextPrompt.title }
          </p>
        </div>
      </div>
    </PromptField>
  );
}
