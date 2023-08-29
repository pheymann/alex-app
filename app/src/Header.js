import { Link } from "react-router-dom";
import "./Header.css";

export default function Header({ signOut }) {
  return (
    <header className="App-header">
      <div className="container">
        <div className="row">
          <div className="col-6">
            <Link to="/">Home</Link>
          </div>
          <div className="col-6 text-end">
            <Link to="/login"
                  onClick={() => {
                    signOut()
                      .catch(err => console.log(err));
                  }}>
              Logout
            </Link>
          </div>
        </div>
      </div>
    </header>
  );
}
