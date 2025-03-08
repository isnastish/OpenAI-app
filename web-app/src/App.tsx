import React, { useState } from 'react';
import { BrowserRouter, Routes, Route, useNavigate } from 'react-router';

const LoginView: React.FC = () => {
    const navigate = useNavigate();

    const onSubmitLoginData = (): void => {
        navigate('/ai');
    };

    return (
        <div className="container">
            <div className="row justify-content-center mt-5">
                <div className="col-md-5">
                    <h3 className="text-center">
                        <p className="font-monospace">Login</p>
                    </h3>
                    <div className="input-group flex-nowrap">
                        <span className="input-group-text" id="addon-wrapping">
                            email
                        </span>
                        <input
                            type="text"
                            autoFocus={true}
                            required
                            className="form-control"
                            placeholder="admin@gmail.com"
                            aria-label="Username"
                            aria-describedby="addon-wrapping"
                        />
                    </div>
                    <hr></hr>
                    <div className="col-12">
                        <button
                            type="submit"
                            className="btn btn-outline-primary"
                            onClick={onSubmitLoginData}
                        >
                            Sign in
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

const AiView: React.FC = () => {
    const [aiQuestion, setAiQuestion] = useState<string>('');
    const askAi = async () => {
        throw new Error(aiQuestion);
    };

    return (
        <div className="container">
            <div className="row justify-content-center mt-5">
                <div className="col-md-7">
                    <h3 className="text-center">
                        <p className="font-monospace">Ask AI</p>
                    </h3>
                    <div className="mb-3">
                        <textarea
                            className="form-control rounded-4"
                            id="floatingTextarea2"
                            placeholder="Ask any question"
                            rows={5}
                            autoFocus={true}
                            onChange={(e) => setAiQuestion(e.target.value)}
                        ></textarea>
                    </div>
                    <button
                        type="button"
                        className="btn btn-outline-primary btn-lg"
                        onClick={askAi}
                    >
                        Submit
                    </button>
                </div>
            </div>
        </div>
    );
};

function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/login" element={<LoginView />}></Route>
                <Route path="/" element={<LoginView />}></Route>
                <Route path="/ai" element={<AiView />}></Route>
            </Routes>
        </BrowserRouter>
    );
}

export default App;
