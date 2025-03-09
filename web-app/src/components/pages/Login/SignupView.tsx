import { useState } from "react";
import { useNavigate } from "react-router";

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

export default SignupView;