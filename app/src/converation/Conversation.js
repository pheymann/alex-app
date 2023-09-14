import { useEffect, useRef, useState } from "react";
import { useParams } from "react-router-dom";
import BasicPage from "../BasicPage";
import ArtContextPromptField from "./ArtContextPrompField";
import UserMessageField from "./UserMessageField";
import QuestionPromptField from "./QuestionPromptField";
import './Conversation.css';
import { LoadingPromptField } from "./PromptField";
import AssistantResponseField from "./AssistantResponseField"
import { logError, pushLogMessage } from "../logger";

export default function Conversation({ awsContext }) {
  const pathParams = useParams();
  const conversationId = pathParams.id;
  const isNewConversation = !conversationId || conversationId === 'new';

  const [conversation, setConversation] = useState(null);
  const [loading, setLoading] = useState(true);

  const logEntriesRef = useRef([]);

  useEffect(() => {
    const token = awsContext.token;

    if (isNewConversation) {
      setConversation({
        messages: [{
        role: 'prompt-art-context',
        }],
      });
      setLoading(false);
    }
    else {
      fetch(`/api/conversation/${conversationId}`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      })
        .then(response => response.text())
        .then(rawData => {
          pushLogMessage(logEntriesRef, { level: 'debug', message: rawData });

          const json = JSON.parse(rawData);

          setConversation({
            ...json,
            messages: [
              {
                role: 'user',
                text: `Tell me something about ${json.metadata.artContext}`,
              },
              ...json.messages,
              {
                role: 'prompt-user-question',
              },
            ],
          });
        })
        .catch(error => {
          logError({ token, error, logEntriesRef: logEntriesRef});
          alert('Error getting conversation:\n' + error);
        })
        .finally(() => {
          setLoading(false);
        });
      }
    },
    [isNewConversation, conversationId, awsContext.token]
  );

  if (loading) {
    return(
      <BasicPage awsContext={awsContext}>
        <div className="container container-limited-width d-flex justify-content-center">
          <div className="spinner-border" role="status">
            <span className="visually-hidden">Loading...</span>
          </div>
        </div>
      </BasicPage>
    );
  }

  const containerizedFields = ['user', 'assistant']

  return (
    <BasicPage awsContext={ awsContext }>
      <div className='container container-limited-width'>
        <div>
          {
            conversation.messages
              .filter(message => containerizedFields.includes(message.role))
              .map((message, index) => {
                const key = `${message.speechClipUuid}_${index}`;

                switch (message.role) {
                  case 'user':
                    return <UserMessageField key={ key } message={ message } />

                  case 'assistant':
                    return <AssistantResponseField key={ key } index={ index } message={ message } />

                  default:
                    console.error(`unknown message role: ${message.role}`);
                    alert(`unknown message role: ${message.role}`);
                    return <div key={ key }></div>
                }
              })
          }
        </div>
      </div>

      {/* outside container elements */}
      <div>
        {
          conversation.messages
            .filter(message => !containerizedFields.includes(message.role))
            .map((message, index) => {
              const key = `${message.speechClipUuid}_${index}`;

              switch (message.role) {
                case 'prompt-art-context':
                  return <ArtContextPromptField key={ key }
                                                conversation={ conversation }
                                                setConversation={ setConversation }
                                                awsContext={ awsContext } />

                case 'prompt-user-question':
                  return <QuestionPromptField key={ key }
                                              conversation={ conversation }
                                              setConversation={ setConversation }
                                              awsContext={ awsContext } />

                case 'loading':
                  return <LoadingPromptField key={ key } />

                default:
                  logError({ awsContext, error: `unknown message role: ${message.role}`, logEntriesRef: logEntriesRef});
                  return <div key={ key }></div>
              }
            })
        }
      </div>
    </BasicPage>
  );
}
