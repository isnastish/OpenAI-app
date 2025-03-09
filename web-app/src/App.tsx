import AiView from './components/pages/Ai/AiView';

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
                                autoFocus={true}
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
    const [emailError, setEmailError] = useState('');
    const [password, setPassword] = useState('');
    const [passwordError, setPasswordError] = useState('');

    const [loginError, setLoginError] = useState('');

    const navigate = useNavigate();
    const onSubmittedLoginData = async () => {
        let hasError: boolean = false;

        if (!email) {
            setEmailError('email cannot be empty');
            hasError = true;
        }
        if (!password || password.length < 8 || password.length > 128) {
            setPasswordError(
                'password length should be greater than 8 and less than 128'
            );
            hasError = true;
        }

        if (hasError) {
            return;
        }

        //
        // TODO: Try using axios instead.
        //

        try {
            const resp = await fetch('/api/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ email: email, password: password }),
                credentials: 'include',
            });

            if (resp.status === 200) {
                navigate('/ai');
                return;
            }

            if (resp.status === 401 || resp.status === 500) {
                const error = await resp.text();
                return;
            }

            throw new Error(`HTTP error, status ${resp.status}`);
        } catch (err) {
            // TODO: handle this code path properly
            setLoginError('unhandled error');
        }
    };

    return (
        <div className="container">
            <div className="row justify-content-center mt-5">
                <div className="col-md-5">
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
                            onChange={(e) => setEmail(e.target.value)}
                        />
                    </div>
                    <p className="fw-lighter text-danger">{emailError}</p>
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
                            onChange={(e) => setPassword(e.target.value)}
                        />
                    </div>
                    <p className="fw-lighter text-danger">{passwordError}</p>
                    <div className="text-end">
                        <button
                            type="submit"
                            className="btn btn-outline-primary"
                            onClick={onSubmittedLoginData}
                        >
                            Login
                        </button>
                    </div>
                    <hr />
                    <div className="d-flex justify-content-between">
                        <span>Don&apos;t have account?</span>
                        <button
                            className="btn btn-outline-danger"
                            onClick={() => {
                                navigate('/signup');
                            }}
                        >
                            Sign up
                        </button>
                    </div>
                    <p className="fw-lighter text-danger">{loginError}</p>
                </div>
            </div>
        </div>
    );
};

const SignupView: React.FC = () => {
    const [username, setUsername] = useState('');
    const [usernameError, setUsernameError] = useState('');
    const [email, setEmail] = useState('');
    const [emailError, setEmailError] = useState('');
    const [password, setPassword] = useState('');
    const [passwordError, setPasswordError] = useState('');
    const [password2, setPassword2] = useState(''); // confirmation password
    const [password2Error, setPassword2Error] = useState('');

    const navigate = useNavigate();

    const onSubmittedSignupData = async (): Promise<void> => {
        let hasError: boolean = false;

        if (!username) {
            setUsernameError('username cannot be empty');
            hasError = true;
        }
        if (!email) {
            setEmailError('email cannot be empty');
            hasError = true;
        }
        if (!password || password.length < 8 || password.length > 128) {
            setPasswordError('password should be in a range [8, 128)');
            hasError = true;
        }
        if ((!password && !password2) || password2 !== password) {
            setPassword2Error(`password doesn't match`);
            hasError = true;
        }

        if (hasError) {
            return;
        }

        navigate('/ai');
    };

    return (
        <div className="container">
            <div className="row justify-content-center mt-5">
                <div className="col-md-5">
                    <h3 className="text-center">Sign up</h3>
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
                    <p className="fw-lighter text-danger">{emailError}</p>
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
                    <p className="fw-lighter text-danger">{usernameError}</p>
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
                    <p className="fw-lighter text-danger">{passwordError}</p>
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
                    <p className="fw-lighter text-danger">{password2Error}</p>
                    <div className="text-end">
                        <button
                            type="submit"
                            className="btn btn-outline-primary"
                            onClick={onSubmittedSignupData}
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
