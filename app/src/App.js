import { Routes, Route, useNavigate } from 'react-router-dom';
import { Amplify, Auth } from 'aws-amplify';
import awsExports from './aws-exports';
import Home from './Home';
import Conversation from './converation/Conversation';
import { useEffect, useState } from 'react';
import './App.css';
import { LanguageLocalStorageKey, decodeLanguage } from './language';
import SignIn from './SignIn';

Amplify.configure({
  Auth: {
      region: awsExports.REGION,
      userPoolId: awsExports.USER_POOL_ID,
      userPoolWebClientId: awsExports.USER_POOL_APP_CLIENT_ID
  }
})

export function App({ validateSession, buildAwsFetch, defaultLanguage }) {
  const [loading, setLoading] = useState(true);
  const [language, setLanguage] = useState(defaultLanguage);
  const [subDomain, setSubDomain] = useState(null);

  const navigate = useNavigate();

  useEffect(() => {
    const checkAuth = async () => {
      try {
        const subDomain = /:\/\/([^/.]+)\..+/.exec(window.location.href);
        if (subDomain !== null && subDomain.length > 0) {
          console.log(subDomain[1]);
          setSubDomain(subDomain[1]);
        }

        // check that we are signed in
        await validateSession();
        setLoading(false);
      } catch (error) {
        setLoading(false);
        navigate('/signin');
      }
    };

    const localLanguage = localStorage.getItem(LanguageLocalStorageKey);
    if (localLanguage !== null) {
      setLanguage(decodeLanguage(localLanguage));
    }

    checkAuth();
  }, [navigate, validateSession, buildAwsFetch, language]);

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
      <Route path='/signin' element={
        <SignIn subDomain={ subDomain }
                awsFetch={ buildAwsFetch(language) }
                language={ language }
                setLanguage={ setLanguage }
                signOut={ () => Auth.signOut() } />
      } />
      <Route exact path='/' element={
        <Home awsFetch={ buildAwsFetch(language) }
              language={ language }
              setLanguage={ setLanguage }
              signOut={ () => Auth.signOut() } />
      } />
      <Route path='/conversation/:id' element={
        <Conversation awsFetch={ buildAwsFetch(language) }
                      language={ language }
                      setLanguage={ setLanguage }
                      signOut={ () => Auth.signOut() } />
      } />
    </Routes>
  );
}
