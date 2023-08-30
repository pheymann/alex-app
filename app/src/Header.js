import { Link, useNavigate } from "react-router-dom";
import "./Header.css";

export default function Header({ awsContext }) {
  const navigate = useNavigate();

  return (
    <header>
      <div className="container header-container-limited-width">
        <div className="row">
          <div className="col-6">
            <Link className="btn" to="/">Home</Link>
          </div>
          <div className="col-6 text-end">
            <button className='btn app-header-logout-button'
                    onClick={() => {
              awsContext.signOut()
                .then(() => navigate('/login'))
                .catch(err => console.log(err));
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
