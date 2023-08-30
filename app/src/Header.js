import { Link, useNavigate } from "react-router-dom";
import "./Header.css";

export default function Header({ awsContext }) {
  const navigate = useNavigate();

  return (
    <header className="App-header">
      <div className="container">
        <div className="row">
          <div className="col-6">
            <Link to="/">Home</Link>
          </div>
          <div className="col-6 text-end">
            <button onClick={() => {
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
