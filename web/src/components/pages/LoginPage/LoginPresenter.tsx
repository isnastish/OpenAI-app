import { FormEvent, useState } from 'react';
import { useNavigate } from 'react-router-dom';

import LoginView from './LoginView';
import SignUpView from './SignupView';

const LoginPresenter: React.FC = () => {
    const [firstName, setFirstName] = useState<string>('');
    const [lastName, setLastName] = useState<string>('');

    const [email, setEmail] = useState<string>('');
    const [emailError, setEmailError] = useState<string>('');

    const [password, setPassword] = useState<string>('');
    const [passwordError, setPasswordError] = useState<string>('');

    const [confirmedPassword, setConfirmedPassword] = useState<string>('');
    const [password1Error, setPassword1Error] = useState<string>('');

    const [authError, setAuthError] = useState<string>('');

    const [accountExists, setAccountExists] = useState<boolean>(false);

    const navigate = useNavigate();

    const clearAll = () => {
        console.log(`clear all input data and errors`);
    };

    const handleLogin = async (event: FormEvent) => {
        event.preventDefault();
        let hasError = false;

        if (!email.trim().length) {
            setEmailError('Email address cannot be empty');
            hasError = true;
        }

        if (password.length <= 8 || password.length > 64) {
            setPasswordError(
                'Password should be more than 8 and less than 64 characters'
            );
        }

        if (hasError) {
            console.log('An error occured');
            return;
        }

        try {
            const resp = await fetch('/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ email: email, password: password }),
                credentials: 'include',
            });

            if (!resp.ok) {
                if (resp.status === 401) {
                    setAuthError(await resp.text());
                } else {
                    // NOTE: Setting auth error for testing purposes now.
                    setAuthError(await resp.text());
                    // throw new Error(await resp.text());
                }

                return;
            }

            navigate('/opeani');
        } catch (err) {
            // TODO: Dispaly an error message, that something went wrong.
            // We probably should create a separate component for that.
        }
    };

    const handleSignUp = async (event: FormEvent) => {
        event.preventDefault();

        let hasError = false;

        if (password !== confirmedPassword) {
            setPassword1Error("Passwords don't match");
            hasError = true;
        }

        if (!hasError && !email.trim().length) {
            // email not set
            hasError = true;
        }

        if (!hasError && !firstName.trim().length) {
            // first name is not specified
            hasError = true;
        }

        if (!hasError && !lastName.trim().length) {
            // last name is not specified
            hasError = true;
        }

        if (!hasError) {
            const requestBody = {
                firstname: firstName,
                lastname: lastName,
                email: email,
                password: password,
            };

            try {
                // const resp = await fetch('/signup', {
                //     method: 'POST',
                //     headers: {
                //         'Content-Type': 'application/json',
                //     },
                //     credentials: 'include',
                //     body: JSON.stringify(requestBody),
                // });

                // if (resp.status !== 200) {
                //     const errorMsg = await resp.text();
                //     throw new Error(errorMsg);
                // } else {
                //     // NOTE: An article about useHistory & useNavigate:
                //     // https://medium.com/@arshguleria1612/upgrade-your-react-navigation-replace-usehistory-with-usenavigate-for-efficient-routing-1708eb7ad672
                //     // NOTE: Should navigate to /home page instead.
                //     // setAccountExists(true);
                //     navigate('/openai');
                // }
                navigate('/openai');
            } catch (err) {
                console.error(err);
            }
        } else {
            console.log('TODO');
        }
    };

    if (accountExists) {
        return (
            <LoginView
                email={email}
                setEmail={setEmail}
                emailError={emailError}
                password={password}
                setPassword={setPassword}
                passwordError={passwordError}
                handleLogin={handleLogin}
                authError={authError}
                accountExists={accountExists}
                setAccountExists={setAccountExists}
                clearAll={clearAll}
            />
        );
    } else {
        return (
            <SignUpView
                firstName={firstName}
                setFirstName={setFirstName}
                lastName={lastName}
                setLastName={setLastName}
                email={email}
                setEmail={setEmail}
                password={password}
                setPassword={setPassword}
                confirmedPassword={confirmedPassword}
                setConfirmedPassword={setConfirmedPassword}
                handleSignUp={handleSignUp}
                clearAll={clearAll}
                accountExists={accountExists}
                setAccountExists={setAccountExists}
            />
        );
    }
};

export default LoginPresenter;
