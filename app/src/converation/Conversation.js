import { useEffect, useRef, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import BasicPage from "../BasicPage";
import ArtContextPromptField from "./ArtContextPrompField";
import UserMessageField from "./UserMessageField";
import QuestionPromptField from "./QuestionPromptField";
import './Conversation.css';
import { LoadingPromptField } from "./PromptField";
import AssistantResponseField from "./AssistantResponseField"
import { logError, pushLogMessage } from "../logger";
import ErrorField from "./ErrorField";
import { Translation } from "../i18n";

export default function Conversation({ awsFetch, language, setLanguage, signOut }) {
  const pathParams = useParams();
  const conversationId = pathParams.id;
  const isNewConversation = !conversationId || conversationId === 'new';

  const [conversation, setConversation] = useState(null);
  const [loading, setLoading] = useState(true);

  const logEntriesRef = useRef([]);
  const navigate = useNavigate();

  const i18n = Translation.get(language);

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
      awsFetch.call(`/api/conversation/${conversationId}`, {
        method: 'GET',
      })
        .then(
          rawData => {
            pushLogMessage(logEntriesRef, { level: 'debug', message: rawData });

            const json = JSON.parse(rawData);

            setConversation({
              ...json,
              messages: [
                {
                  role: 'user',
                  text: `${i18n.conversation.initialField} ${json.metadata.artContext}`,
                },
                ...json.messages,
                {
                  role: 'prompt-user-question',
                },
              ],
            });
          },
          error =>  {
            // make sure logging has minimal side effect
            logError({ awsFetch, error, logEntriesRef: logEntriesRef});
            return Promise.reject(error);
          },
        )
        .catch(_ => {
          navigate('/');
        })
        .finally(() => {
          setLoading(false);
        });
      }
    },
    [isNewConversation, conversationId, awsFetch, navigate, i18n]
  );

  if (loading) {
    return(
      <BasicPage  awsFetch={ awsFetch }
                  language={ language }
                  setLanguage={ setLanguage }
                  signOut={ signOut }
      >
        <div className="container container-limited-width d-flex justify-content-center">
          <div className="spinner-border" role="status">
            <span className="visually-hidden">Loading...</span>
          </div>
        </div>
      </BasicPage>
    );
  }

  const containerizedFields = ['user', 'assistant', 'error']

  return (
    <BasicPage  awsFetch={ awsFetch }
                language={ language }
                setLanguage={ setLanguage }
                signOut={ signOut }
    >
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

                  case 'error':
                    return <ErrorField key={ key } errorCode={ message.errorCode } i18n={ i18n }/>

                  default:
                    pushLogMessage(logEntriesRef, { level: 'error', message: `unknown message role: ${message.role}` });
                    navigate('/?errorCode=0');
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
                                                i18n= { i18n }
                                                awsFetch={ awsFetch } />

                case 'prompt-user-question':
                  return <QuestionPromptField key={ key }
                                              conversation={ conversation }
                                              setConversation={ setConversation }
                                              i18n= { i18n }
                                              awsFetch={ awsFetch } />

                case 'loading':
                  return <LoadingPromptField key={ key } />

                default:
                  logError({ awsFetch, error: `unknown message role: ${message.role}`, logEntriesRef: logEntriesRef});
                  return <div key={ key }></div>
              }
            })
        }
      </div>
    </BasicPage>
  );
}
