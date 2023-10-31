import React, { useCallback, useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Auth } from "aws-amplify";
import BasicPage from "./BasicPage";
import { Translation } from "./i18n";
import { Errors, errorAlertMessage } from "./ErrorAlert";
import './SignIn.css';
import './BasicStyling.css'

const promotions = {
  "hamburger-kunsthalle": `W&bnQJ.X^hq3{wC"H'Pu4>`,
}

export default function SignIn({ subDomain, awsFetch, language, setLanguage, signOut }) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);
  const [newPassword, setNewPassword] = useState(false);
  const [user, setUser] = useState(null);

  const buttonRef = useRef(null);

  const navigate = useNavigate();

  const i18n = Translation.get(language);

  const handleSignIn = useCallback((directEmail, directPassword) => {
    setLoading(true);

    const signIn = async () => {
      try {
        const finalEmail = directEmail || email;
        const finalPassword = directPassword || password;

        const user = await Auth.signIn(finalEmail, finalPassword);

        if (user.challengeName === 'NEW_PASSWORD_REQUIRED') {
          setUser(user);
          setNewPassword(true);
        } else {
          navigate('/');
        }
      } catch (error) {
        if (error.message === 'Pending sign-in attempt already in progress') {
          return;
        }
        setError(Errors.SignInError);
      }
      setLoading(false);
    };
    signIn();
  }, [email, navigate, password]);

  useEffect(() => {
    // promotion sign in
    if (subDomain !== undefined && promotions[subDomain] !== undefined) {
      const email = `${subDomain}@sprichmitalex.de`;
      setEmail(email);
      setPassword(promotions[subDomain]);

      handleSignIn(email, promotions[subDomain]);
    }

    const handleKeyPress = (event) => {
      if (event.key === 'Enter') {
        buttonRef.current.click();
      }
    };

    document.addEventListener('keydown', handleKeyPress);

    return () => {
      document.removeEventListener('keydown', handleKeyPress);
    };
  }, [subDomain, handleSignIn]);

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

        { newPassword &&
          <ChangePassword i18n={ i18n } user={ user } setError={ setError } />
        }

        { !newPassword &&
          <div>
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
        }
      </div>
    </BasicPage>
  );
}

function ChangePassword({ i18n, user, setError }) {
  const [loading, setLoading] = useState(false);
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const buttonRef = useRef(null);
  const navigate = useNavigate();

  const handleChangePassword = async () => {
    setLoading(true);
    try {
      if (newPassword !== confirmPassword) {
        setError(Errors.CompleteSignUpError);
        setLoading(false);
        return;
      }

      await Auth.completeNewPassword(user, newPassword);

      navigate('/');
    } catch (error) {
      setError(Errors.CompleteSignUpError);
    }
    setLoading(false);
  };

  return (
    <div>
      <div className='row'>
        <div className='col text-center'>
          <input  type="password"
                  className="sign-in-input"
                  value={ newPassword }
                  placeholder={ i18n.signIn.newPassword }
                  autoFocus={ true }
                  onChange={ (e) => setNewPassword(e.target.value) }
          />
        </div>
      </div>
      <div className='row'>
        <div className='col text-center'>
          <input  type="password"
                  className="sign-in-input"
                  value={ confirmPassword }
                  placeholder={ i18n.signIn.confirmPassword }
                  onChange={ (e) => setConfirmPassword(e.target.value) }
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
                    onClick={ () => handleChangePassword() }
            >
              { i18n.signIn.changePasswordButton }
            </button>
          }
        </div>
      </div>
    </div>
  );
}
