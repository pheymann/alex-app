import React, { useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Auth } from "aws-amplify";
import BasicPage from "./BasicPage";
import { Translation } from "./i18n";
import { Errors, errorAlertMessage } from "./ErrorAlert";
import './SignIn.css';
import './BasicStyling.css'

export default function SignIn({ awsFetch, language, setLanguage, signOut }) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);
  const buttonRef = useRef(null);

  const navigate = useNavigate();

  const i18n = Translation.get(language);

  const handleSignIn = async () => {
    setLoading(true);
    try {
      await Auth.signIn(email, password);
      navigate('/');
    } catch (error) {
      setError(Errors.SignInError);
    }
    setLoading(false);
  };

  useEffect(() => {
    const handleKeyPress = (event) => {
      if (event.key === 'Enter') {
        buttonRef.current.click();
      }
    };

    document.addEventListener('keydown', handleKeyPress);

    return () => {
      document.removeEventListener('keydown', handleKeyPress);
    };
  }, []);

  const errorMessage = errorAlertMessage(error, i18n);

  return (
    <BasicPage  awsFetch={ awsFetch }
                language={ language }
                setLanguage={ setLanguage }
                signOut={ signOut }
    >
      <div className="container container-limited-width sign-in">
        { error &&
          <div className='row'>
            <div className='col text-center alert alert-warning'>
              { errorMessage }
            </div>
          </div>
        }

        <div className='row'>
          <div className='col text-center'>
            <input  type="text"
                    className="sign-in-input"
                    value={ email }
                    placeholder={ i18n.signIn.emailPlaceholder }
                    autoFocus={ true }
                    onChange={ (e) => setEmail(e.target.value) }
            />
          </div>
        </div>
        <div className='row'>
          <div className='col text-center'>
            <input  type="password"
                    className="sign-in-input"
                    value={ password }
                    placeholder={ i18n.signIn.passwordPlaceholder }
                    onChange={ (e) => setPassword(e.target.value) }
            />
          </div>
        </div>
        <div className='row'>
          <div className='col d-flex justify-content-center'>
            { loading &&
              <div className='sign-in-button d-flex justify-content-center'>
                <div className="spinner-border" role="status">
                  <span className="visually-hidden">Loading...</span>
                </div>
              </div>
            }

            { !loading &&
              <button ref={ buttonRef }
                      className='sign-in-button'
                      onClick={ () => handleSignIn() }
              >
                { i18n.signIn.button }
              </button>
            }
          </div>
        </div>
      </div>
    </BasicPage>
  );
}
