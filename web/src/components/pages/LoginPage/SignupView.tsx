import React, { Fragment, FormEvent } from 'react';

import FormInput from '../../common/inputs/FormInput';
import './LoginPresenter.css';

interface SignUpData {
    firstName: string;
    setFirstName: (firstName: string) => void;

    lastName: string;
    setLastName: (lastName: string) => void;

    email: string;
    setEmail: (email: string) => void;

    password: string;
    setPassword: (password: string) => void;

    confirmedPassword: string;
    setConfirmedPassword: (password: string) => void;

    handleSignUp: (event: FormEvent) => void;

    clearAll: () => void;

    accountExists: boolean;

    setAccountExists: (accountExists: boolean) => void;
}

const SignUpView: React.FC<SignUpData> = ({
    firstName,
    setFirstName,
    lastName,
    setLastName,
    email,
    setEmail,
    password,
    setPassword,
    confirmedPassword,
    setConfirmedPassword,
    handleSignUp,
    clearAll,
    accountExists,
    setAccountExists,
}) => {
    return (
        <Fragment>
            <div className="login-presenter-class">
                <form id="signup-form" onSubmit={(e) => handleSignUp(e)}>
                    <h2>Sign Up</h2>
                    <FormInput
                        isAutoFocus={true}
                        labelText="First name"
                        value={firstName}
                        onChangeHandler={(e) => setFirstName(e.target.value)}
                    />
                    <FormInput
                        labelText="Last name"
                        value={lastName}
                        onChangeHandler={(e) => setLastName(e.target.value)}
                    />
                    <FormInput
                        labelText="Email"
                        value={email}
                        onChangeHandler={(e) => setEmail(e.target.value)}
                    />
                    <FormInput
                        labelText="Password"
                        value={password}
                        onChangeHandler={(e) => setPassword(e.target.value)}
                    />
                    <FormInput
                        labelText="Confirm password"
                        value={confirmedPassword}
                        onChangeHandler={(e) =>
                            setConfirmedPassword(e.target.value)
                        }
                    />
                </form>
                <p className="form-actions">
                    <button
                        className="submit-button"
                        form="signup-form"
                        onClick={clearAll}
                    >
                        Create account
                    </button>
                </p>
            </div>
            <div className="login-presenter-class">
                <label>Already have an account? </label>
                <button
                    className="submit-button"
                    onClick={() => setAccountExists(!accountExists)}
                >
                    Sign in
                </button>
            </div>
        </Fragment>
    );
};

export { SignUpView, SignUpData }; 
