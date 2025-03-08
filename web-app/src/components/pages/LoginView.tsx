import React, { Fragment, FormEvent } from 'react';

interface LoginData {
    email: string;
    emailError: string;
    setEmail: (email: string) => void;

    password: string;
    passwordError: string;
    setPassword: (password: string) => void;

    handleLogin: (event: FormEvent) => void;
    authError: string;

    accountExists: boolean;
    setAccountExists: (accountExists: boolean) => void;

    clearAll: () => void;
}

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
        // <div className="login-presenter-class">
        //     <form id="login-form" onSubmit={(e) => handleLogin(e)}>
        //         <h2>Login</h2>
        //         <div className="my-class">
        //             <FormInput
        //                 isAutoFocus={true}
        //                 labelText="Email"
        //                 value={email}
        //                 onChangeHandler={(e) => setEmail(e.target.value)}
        //             />
        //             <p className="errorText">{emailError}</p>
        //             <FormInput
        //                 labelText="Password"
        //                 value={password}
        //                 onChangeHandler={(e) => setPassword(e.target.value)}
        //             />
        //             <p className="errorText">{passwordError}</p>
        //         </div>
        //     </form>
        //     <p>
        //         <button
        //             onClick={clearAll}
        //             className="submit-button"
        //             form="login-form"
        //         >
        //             Sign In
        //         </button>
        //     </p>
        //     <p className="errorText">{authError}</p>
        // </div>
        // <div className="login-presenter-class">
        //     <label>Don&#39;t have an account? </label>
        //     <button
        //         className="submit-button"
        //         onClick={() => setAccountExists(!accountExists)}
        //     >
        //         Sing up
        //     </button>
        // </div>
        <div></div>
    );
};

export { LoginView, LoginData };
