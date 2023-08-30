import { Link } from "react-router-dom";
import "./NewConversationButton.css";

export default function NewConversationButton(props) {
  const compositeClasses = `${props.className || ''} new-conversation-button`;

  return (
    <Link className={compositeClasses} to={'/conversation/new'}>
      New conversation
    </Link>
  );
}
