import React from 'react';

function App() {
    return (
        <div className="container">
          <div className='row justify-content-center mt-5'>
            <div className='col-md-5'>
                <h3 className="text-center">Login</h3>
                <div className="input-group flex-nowrap">
                  <span className="input-group-text" id="addon-wrapping">email</span>
                  <input type="text" autoFocus={true} required className="form-control" placeholder="admin@gmail.com" aria-label="Username" aria-describedby="addon-wrapping"/>
                </div>
                <hr></hr>
                <div className="col-12">
                  <button type="submit" className="btn btn-outline-primary">Sign in</button>
                </div>
            </div>
          </div>
        </div>
    );
}

export default App;
