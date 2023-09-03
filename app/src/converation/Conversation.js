import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import BasicPage from "../BasicPage";
import ArtContextPromptField from "./ArtContextPrompField";
import UserMessageField from "./UserMessageField";
import QuestionPromptField from "./QuestionPromptField";
import './Conversation.css';
import { LoadingPromptField } from "./PromptField";
import AssistantResponseField from "./AssistantResponseField";

export default function Conversation({ awsContext }) {
  const pathParams = useParams();
  const conversationId = pathParams.id;
  const isNewConversation = !conversationId || conversationId === 'new';

  const [conversation, setConversation] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
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
          'Authorization': `Bearer ${awsContext.token}`,
        },
      })
        .then(response => response.json())
        .then(data => {
          setConversation({
            ...data,
            messages: [
              {
                role: 'user',
                text: `Tell me something about ${data.metadata.artContext}`,
              },
              ...data.messages,
              {
                role: 'prompt-user-question',
              },
            ],
          });
        })
        .catch(error => {
          console.log(error);
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

  return (
    <BasicPage awsContext={ awsContext }>
      <div className='container container-limited-width'>
        <div>
          {
            conversation.messages.map((message, index) => {
              const key = `${message.speechClipUuid}_${index}`;

              switch (message.role) {
                case 'user':
                  return <UserMessageField key={ key } message={ message } />

                case 'assistant':
                  return <AssistantResponseField key={ key } index={ index } message={ message } />

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
                  console.error(`unknown message role: ${message.role}`);
                  return <div key={ key }></div>
              }
            })
          }
        </div>
      </div>
    </BasicPage>
  );
}




