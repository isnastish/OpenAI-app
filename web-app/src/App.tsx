import React, { useState } from 'react';
import { BrowserRouter, Routes, Route, useNavigate } from 'react-router';

const ErrorView: React.FC = () => {
    return (
        <div className="container-fluid vh-100">
            <div className="row justify-content-center align-items-center h-100">
                <div className="col-md-4">
                    <form>
                        <div className="mb-3">
                            <label htmlFor="email" className="form-label">
                                Email address
                            </label>
                            <input
                                type="email"
                                className="form-control"
                                id="email"
                                placeholder="Enter your email"
                            />
                        </div>
                        <div className="mb-3">
                            <label htmlFor="password" className="form-label">
                                Password
                            </label>
                            <input
                                type="password"
                                className="form-control"
                                id="password"
                                placeholder="Enter your password"
                            />
                        </div>
                        <div className="mb-3">
                            <button
                                type="submit"
                                className="btn btn-primary w-100"
                            >
                                Submit
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
};

const LoginView: React.FC = () => {
    const [email, setEmail] = useState('');

    const navigate = useNavigate();
    const onSubmittedLoginData = async () => {
        navigate('/ai');
    };

    return (
        <div className="container">
            <div className="row justify-content-center mt-5">
                <div className="col-md-5 font-monospace">
                    <h3 className="text-center">
                        <p className="font-monospace">Login</p>
                    </h3>
                    <div className="input-group mb-3">
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
                    <div className="input-group mb-3">
                        <span className="input-group-text" id="addon-wrapping">
                            password
                        </span>
                        <input
                            type="text"
                            required
                            className="form-control"
                            placeholder="********"
                            aria-describedby="addon-wrapping" /* TODO: Figure out why do we need this.*/
                        />
                    </div>
                    <div className="text-end">
                        <button
                            type="submit"
                            className="btn btn-outline-primary"
                            onClick={onSubmittedLoginData}
                        >
                            Login
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

const SignupView: React.FC = () => {
    const [username, setUsername] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [password2, setPassword2] = useState(''); // confirmation password

    const navigate = useNavigate();

    const onSubmitLoginData = async (): Promise<void> => {
        navigate('/ai');
    };

    return (
        <div className="container">
            <div className="row justify-content-center mt-5">
                <div className="col-md-5 font-monospace">
                    <h3 className="text-center">
                        <p className="font-monospace">Sign up</p>
                    </h3>
                    <div className="input-group mb-3">
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
                    <div className="input-group mb-3">
                        <span className="input-group-text" id="addon-wrapping">
                            name
                        </span>
                        <input
                            type="text"
                            required
                            className="form-control"
                            placeholder="Ivan Ivanov"
                            aria-describedby="addon-wrapping"
                        />
                    </div>
                    <div className="input-group mb-3">
                        <span className="input-group-text" id="addon-wrapping">
                            password
                        </span>
                        <input
                            type="text"
                            required
                            placeholder="********"
                            className="form-control"
                            aria-describedby="addon-wrapping"
                        />
                    </div>
                    <div className="input-group mb-3">
                        <span className="input-group-text" id="addon-wrapping">
                            confirm
                        </span>
                        <input
                            type="text"
                            required
                            placeholder="********"
                            className="form-control"
                            aria-describedby="addon-wrapping"
                        />
                    </div>
                    <div className="text-end">
                        <button
                            type="submit"
                            className="btn btn-outline-primary"
                            onClick={onSubmitLoginData}
                        >
                            Sign up
                        </button>
                    </div>
                    <hr></hr>
                    <div className="d-flex justify-content-between">
                        <span>Already have account?</span>
                        <button
                            className="btn btn-outline-danger"
                            onClick={() => {
                                navigate('/login');
                            }}
                        >
                            Login
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

const AiView: React.FC = () => {
    const [aiQuestion, setAiQuestion] = useState<string>('');
    const askAi = async (): Promise<void> => {
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
                            placeholder="Type any question"
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
                <Route path="/signup" element={<SignupView />} />
                <Route path="/" element={<SignupView />} />
                <Route path="/ai" element={<AiView />} />
                <Route path="/login" element={<LoginView />} />
            </Routes>
        </BrowserRouter>
    );
}

export default App;
