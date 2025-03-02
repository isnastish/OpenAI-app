import React, { Fragment } from 'react';

import LoginData from '../../../models/LoginData';
import FormInput from '../../common/inputs/FormInput';

const LoginView: React.FC<LoginData> = ({
    email,
    setEmail,
    emailError,
    password,
    setPassword,
    passwordError,
    handleLogin,
    authError,
    accountExists,
    setAccountExists,
    clearAll,
}) => {
    return (
        <Fragment>
            <div className="login-presenter-class">
                <form id="login-form" onSubmit={(e) => handleLogin(e)}>
                    <h2>Login</h2>
                    <div className="my-class">
                        <FormInput
                            isAutoFocus={true}
                            labelText="Email"
                            value={email}
                            onChangeHandler={(e) => setEmail(e.target.value)}
                        />
                        <p className="errorText">{emailError}</p>
                        <FormInput
                            labelText="Password"
                            value={password}
                            onChangeHandler={(e) => setPassword(e.target.value)}
                        />
                        <p className="errorText">{passwordError}</p>
                    </div>
                </form>
                <p>
                    <button
                        onClick={clearAll}
                        className="submit-button"
                        form="login-form"
                    >
                        Sign In
                    </button>
                </p>
                <p className="errorText">{authError}</p>
            </div>
            <div className="login-presenter-class">
                <label>Don&#39;t have an account? </label>
                <button
                    className="submit-button"
                    onClick={() => setAccountExists(!accountExists)}
                >
                    Sing up
                </button>
            </div>
        </Fragment>
    );
};

export default LoginView;
