import { Routes, Route } from 'react-router-dom';
import Home from './Home';
import ArtConversation from './ArtConversation';

export default function App() {
  return (
    <Routes>
      <Route exact path="/" element={<Home />} />
      <Route path="/conversation/:id" element={<ArtConversation />} />
    </Routes>
  );
}
