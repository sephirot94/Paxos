import './App.css';
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";

import Wallet from "./components/wallet.component";
import React from 'react';

class App extends React.Component {

  render() {
    return (
      <Router>
        <div className="App">
          <div className="auth-wrapper">
            <div className="auth-inner">
              <Switch>
                <Route exact path='/' component={Wallet} />
              </Switch>
            </div>
          </div>
        </div>
      </Router>
    )
  };
}

export default App;
