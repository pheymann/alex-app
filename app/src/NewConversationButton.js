import { Link } from "react-router-dom";
import "./NewConversationButton.css";

export default function NewConversationButton() {
  return (
    <Link className='new-conversation-button' to={'/conversation/new'}>
      New conversation
    </Link>
  );
}
