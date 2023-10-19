import { Link, useNavigate } from "react-router-dom";
import "./Header.css";
import { logError } from "./logger";
import { useRef } from "react";
import { LanguageLocalStorageKey, encodeLanguage, prettyPrintLanguage, switchToLanguageFrom } from "./language";
import { Translation } from "./i18n";

export default function Header({ awsFetch, language, setLanguage, signOut }) {
  const navigate = useNavigate();
  const logEntriesRef = useRef([]);

  const switchToLang = switchToLanguageFrom(language);
  const i18n = Translation.get(language);

  return (
    <header>
      <div className="container header-container-limited-width">
        <div className="row">
          <div className="col-6">
            <Link className="btn" to="/">
              { i18n.header.home }
            </Link>
          </div>
          <div className="col-2">
          </div>
          <div className="col-2 text-end">
            <button className="app-header-button"
                    onClick={() => {
                      setLanguage(switchToLang);
                      localStorage.setItem(LanguageLocalStorageKey, encodeLanguage(switchToLang));
                    }}
            >
              { prettyPrintLanguage(switchToLang) }
            </button>
          </div>
          <div className="col-2 ">
            <button className='app-header-button'
                    onClick={() => {
              signOut()
                .then(() => navigate('/login'))
                .catch(err => {
                  logError({ awsFetch, error: err, logEntriesRef: logEntriesRef });
                  alert('Error signing out:\n' + err);
                });
              }}
            >
              { i18n.header.signOut }
            </button>
          </div>
        </div>
      </div>
    </header>
  );
}
