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

  const showsArtContextPrompt = conversation.messages.some(message => message.role === 'prompt-art-context');
  const showsUserQuestionPrompt = conversation.messages.some(message => message.role === 'prompt-user-question');

  return (
    <BasicPage awsContext={ awsContext }>
      <div className='container container-limited-width'>
        <div>
          {
            conversation.messages
              .filter(message => message.role !== 'prompt-art-context' && message.role !== 'prompt-user-question')
              .map((message, index) => {
                const key = `${message.speechClipUuid}_${index}`;

                switch (message.role) {
                  case 'user':
                    return <UserMessageField key={ key } message={ message } />

                  case 'assistant':
                    return <AssistantResponseField key={ key } index={ index } message={ message } />

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

      { showsArtContextPrompt &&
        <ArtContextPromptField  conversation={ conversation }
                                setConversation={ setConversation }
                                awsContext={ awsContext } />
      }

      { showsUserQuestionPrompt &&
        <QuestionPromptField  conversation={ conversation }
                              setConversation={ setConversation }
                              awsContext={ awsContext } />
      }
    </BasicPage>
  );
}




