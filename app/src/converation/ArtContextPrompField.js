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

    const startProcessing = async () => {
      try {
        const createResponse = await awsFetch.callResponse(`/api/conversation/create/art`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            artContext: artContext,
          }),
        })

        const createdJson = await createResponse.json();

        // remove loading message
        const responseConversation = {
          ...createdJson,
          messages: [
            conversation.messages[0],
          ],
        };

        if (createResponse.ok) {
          const pollingInterval = setInterval(async () => {
            try {
              const pollResponse = await awsFetch.callResponse(`/api/conversation/${responseConversation.id}/poll`,  {
                method: 'GET',
              });

              if (pollResponse.status === 200) {
                clearInterval(pollingInterval);

                const message = await pollResponse.json();
                pushLogMessage(logEntriesRef, { level: 'debug', message: message });

                responseConversation.messages = [
                  ...responseConversation.messages,
                  message,
                  {
                    role: 'prompt-user-question',
                  },
                ];

                setConversation(responseConversation);
              }
            } catch (error) {
              clearInterval(pollingInterval);
              logError({ awsFetch, error, logEntriesRef: logEntriesRef});
              navigate('/?errorCode=' + errorToCode(Errors.StartingConversationError));
            }
          }, 1000);
        }
      } catch (error) {
          logError({ awsFetch, error, logEntriesRef: logEntriesRef});
          navigate('/?errorCode=' + errorToCode(Errors.StartingConversationError));
      }
    };

    startProcessing();
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
