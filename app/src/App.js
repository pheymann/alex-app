import { Routes, Route } from 'react-router-dom';
import Home from './Home';
import ArtConversation from './ArtConversation';

export default function App() {
  return (
    <Routes>
      <Route exact path="/" element={<Home />} />
      <Route path="/conversation/new" element={<ArtConversation />} />
      <Route
        path="/conversation/:id"
        render={(props) => <ArtConversation conversationId={props.match.params.id}/>} />
    </Routes>
  );
}
