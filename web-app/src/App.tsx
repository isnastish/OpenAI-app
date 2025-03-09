import AiView from './components/pages/Ai/AiView';
import LoginView from './components/pages/Login/LoginView';
import SignupView from './components/pages/Login/SignupView';

import { BrowserRouter, Routes, Route, useNavigate } from 'react-router';

function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/signup" element={<SignupView />} />
                <Route path="/" element={<SignupView />} />
                <Route path="/ai" element={<AiView />} />
                <Route path="/login" element={<LoginView />} />
            </Routes>
        </BrowserRouter>
    );
}

export default App;
