import React from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import UserComponent from './components/UserComponent';

function App() {
  return (
    <Router>
      <div className="App">
        <Switch>
          <Route path="/users" component={UserComponent} />
        </Switch>
      </div>
    </Router>
  );
}

export default App;