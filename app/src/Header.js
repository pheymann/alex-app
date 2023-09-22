import { Link, useNavigate } from "react-router-dom";
import "./Header.css";
import { logError } from "./logger";
import { useRef } from "react";

export default function Header({ awsFetch, signOut }) {
  const navigate = useNavigate();
  const logEntriesRef = useRef([]);

  return (
    <header>
      <div className="container header-container-limited-width">
        <div className="row">
          <div className="col-6">
            <Link className="btn" to="/">Home</Link>
          </div>
          <div className="col-6 text-end">
            <button className='app-header-logout-button'
                    onClick={() => {
              signOut()
                .then(() => navigate('/login'))
                .catch(err => {
                  logError({ awsFetch, error: err, logEntriesRef: logEntriesRef });
                  alert('Error signing out:\n' + err);
                });
              }}
            >
              Logout
            </button>
          </div>
        </div>
      </div>
    </header>
  );
}
