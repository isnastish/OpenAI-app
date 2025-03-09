import { useState } from "react";

const AiView: React.FC = () => {
    const [aiQuestion, setAiQuestion] = useState<string>('');
    const askAi = async (): Promise<void> => {
        if (!aiQuestion) {
            return;
        }
    };

    return (
        <div className="container">
            <div className="row justify-content-center mt-5">
                <div className="col-md-7">
                    <h3 className="text-center">
                        <p className="font-monospace">Ask AI</p>
                    </h3>
                    <div className="mb-3">
                        <textarea
                            className="form-control rounded-4"
                            id="floatingTextarea2"
                            placeholder="Type any question"
                            rows={5}
                            autoFocus={true}
                            onChange={(e) => setAiQuestion(e.target.value)}
                        ></textarea>
                    </div>
                    <div className="text-end">
                        <button
                            type="button"
                            className="btn btn-outline-primary btn-lg"
                            onClick={askAi}
                        >
                            Submit
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default AiView; 