import { FormEvent } from 'react';

interface LoginData {
    email: string;
    setEmail: (email: string) => void;
    emailError: string;

    password: string;
    setPassword: (password: string) => void;
    passwordError: string;

    handleLogin: (event: FormEvent) => void;
    authError: string;

    accountExists: boolean;
    setAccountExists: (accountExists: boolean) => void;

    clearAll: () => void;
}

export default LoginData;
