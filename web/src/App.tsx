import './App.css';
import LoginPresenter from './components/pages/LoginPage/LoginPresenter';
import OpenaiView from './components/pages/OpenaiPage/OpenaiView';

import { BrowserRouter, Routes, Route } from 'react-router-dom';

function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/openai" element={<OpenaiView />} />
                <Route path="/login" element={<LoginPresenter />} />
                <Route path="/" element={<LoginPresenter />} />
            </Routes>
        </BrowserRouter>
    );
}

export default App;
