import { FormEvent } from 'react';

interface SignUpData {
    firstName: string;
    setFirstName: (firstName: string) => void;

    lastName: string;
    setLastName: (lastName: string) => void;

    email: string;
    setEmail: (email: string) => void;
    // TODO: Speify email error as well.

    password: string;
    setPassword: (password: string) => void;

    confirmedPassword: string;
    setConfirmedPassword: (password: string) => void;

    handleSignUp: (event: FormEvent) => void;

    clearAll: () => void;

    accountExists: boolean;

    setAccountExists: (accountExists: boolean) => void;
}

export default SignUpData;
