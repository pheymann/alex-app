import {useState, useEffect, useRef} from 'react';
import { Link } from 'react-router-dom';
import './Home.css';
import './BasicStyling.css'
import NewConversationButton from './NewConversationButton';
import BasicPage from './BasicPage';
import { logError, pushLogMessage } from './logger';
import { Errors, codeToError, errorAlertMessage } from './ErrorAlert';
import { Translation } from './i18n';

export default function Home({ awsFetch, language, setLanguage, signOut }) {
  const [conversations, setConversations] = useState([]);
  const [error, setError] = useState(null);

  const logEntriesRef = useRef([]);

  const i18n = Translation.get(language);

  useEffect(() => {
    awsFetch.call(`/api/conversation/list`, {
      method: 'GET',
    })
      .then(rawData => {
        pushLogMessage(logEntriesRef, { level: 'debug', message: rawData });

        const json = JSON.parse(rawData);
        setConversations(json);
      })
      .catch(error => {
        logError({ awsFetch, error, logEntriesRef: logEntriesRef});
        setError(Errors.ConversationListingError);
      });

      // handle errors triggered by other views
      const urlSearchParams = new URLSearchParams(window.location.search);
      const params = Object.fromEntries(urlSearchParams.entries());

      params.errorCode && setError(codeToError(params.errorCode));
  }, [awsFetch]);

  const errorMessage = errorAlertMessage(error, i18n);

  return (
    <BasicPage  awsFetch={ awsFetch }
                language={ language }
                setLanguage={ setLanguage }
                signOut={ signOut }
    >
      <div className='container container-limited-width'>
        <div className='row'>
          <div className='col text-center'>
            <NewConversationButton  className='home-new-conversation-button'
                                    i18n={ i18n }
            />
          </div>
        </div>

          { error &&
            <div className='row'>
              <div className='col text-center alert alert-warning'>
                { errorMessage }
              </div>
            </div>
          }

          { conversations &&
            conversations.map((conversation, index) => {
              const key = `${conversation.id}_${index}`;

              return (
                <div key={key} className='row'>
                  <div className='col'>
                    <Link className='conversation-link' to={`/conversation/${conversation.id}`}>
                      { conversation.metadata.artContext }
                    </Link>
                  </div>
                </div>
              );
            })
          }
      </div>
    </BasicPage>
  );
}
