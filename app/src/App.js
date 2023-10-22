import { Routes, Route, useNavigate } from 'react-router-dom';
import { Amplify, Auth } from 'aws-amplify';
import awsExports from './aws-exports';
import Home from './Home';
import Login from './Login';
import Conversation from './converation/Conversation';
import { useEffect, useState } from 'react';
import './App.css';
import { LanguageLocalStorageKey, decodeLanguage } from './language';

Amplify.configure({
  Auth: {
      region: awsExports.REGION,
      userPoolId: awsExports.USER_POOL_ID,
      userPoolWebClientId: awsExports.USER_POOL_APP_CLIENT_ID
  }
})

export function App({ loadAwsCtx , buildAwsFetch, defaultLanguage }) {
  const [loading, setLoading] = useState(true);
  const [awsContext, setAwsContext] = useState(null);
  const [awsFetch, setAwsFetch] = useState(null);
  const [language, setLanguage] = useState(defaultLanguage);

  const navigate = useNavigate();

  const urlSearchParams = new URLSearchParams(window.location.search);
  const params = Object.fromEntries(urlSearchParams.entries());
  const forceReload = params.forceReload;

  useEffect(() => {
    const checkAuth = async () => {
      try {
        const awsContext = await loadAwsCtx();

        setAwsContext(awsContext);
        setAwsFetch(buildAwsFetch(awsContext, language));
        setLoading(false);
      } catch (error) {
        console.log(error);
        setLoading(false);
        navigate('/login');
      }
    };

    checkAuth();

    const localLanguage = localStorage.getItem(LanguageLocalStorageKey);
    if (localLanguage !== null) {
      setLanguage(decodeLanguage(localLanguage));
    }
  }, [navigate, loadAwsCtx, buildAwsFetch, language, forceReload]);

  if (loading) {
    return(
      <div className="container">
        <div className="spinner-border" role="status">
          <span className="visually-hidden">Loading...</span>
        </div>
      </div>
    );
  }

  return (
    <Routes>
      <Route path='/login' element={<Login />} />
      <Route exact path='/' element={
        <Home awsFetch={ awsFetch }
              language={ language }
              setLanguage={ setLanguage }
              signOut={ () => awsContext.signOut() } />
      } />
      <Route path='/conversation/:id' element={
        <Conversation awsFetch={ awsFetch }
                      language={ language }
                      setLanguage={ setLanguage }
                      signOut={ () => awsContext.signOut() } />
      } />
    </Routes>
  );
}

export async function defaultLoadAwsCtx() {
  const awsSession = await Auth.currentSession();

  return {
    awsSession,
    token: awsSession.getIdToken().getJwtToken(),
    signOut: () => Auth.signOut(),
  };
}
