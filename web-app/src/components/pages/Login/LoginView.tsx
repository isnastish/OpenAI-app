import React, { useState } from 'react';
import { useNavigate } from 'react-router';

interface LoginData {
    email: string;
    emailError: string;
    setEmail: (email: string) => void;

    password: string;
    passwordError: string;
    setPassword: (password: string) => void;

    // handleLogin: (event) => void;
    authError: string;

    accountExists: boolean;
    setAccountExists: (accountExists: boolean) => void;

    clearAll: () => void;
}

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

export default LoginView;
