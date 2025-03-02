import './FormatInput.css';

import React, { ChangeEvent } from 'react';

interface FormInputData {
    isAutoFocus?: boolean;

    labelText: string;

    value: string;

    onChangeHandler: (event: ChangeEvent<HTMLInputElement>) => void;
}

const FormInput: React.FC<FormInputData> = ({
    isAutoFocus,
    labelText,
    value,
    onChangeHandler,
}) => {
    const inputId = labelText.replace(' ', '-').toLowerCase();
    console.log(inputId);

    return (
        <div className="control no-margin">
            <label htmlFor={inputId}>{labelText}</label>
            <input
                id={inputId}
                type="text"
                autoFocus={isAutoFocus ? true : false}
                required
                value={value}
                onChange={onChangeHandler}
            ></input>
        </div>
    );
};

export default FormInput;
