import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import reportWebVitals from './reportWebVitals';
import {
  HashRouter as Router, Route, Switch 
} from 'react-router-dom';

//Aqui hemos importado boostrap 4
import '../node_modules/bootstrap/dist/css/bootstrap.min.css'; // Archivo CSS de Bootstrap 4 
import '../node_modules/bootstrap/dist/js/bootstrap.min.js'; // Archivo Javascript de Bootstrap 4 

//Importamos los componentes necesarios
import Menu from './componentes/menu/Menu';
import Footer from './componentes/footer/Footer';

const root = ReactDOM.createRoot(document.getElementById('root'));

root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);

/*
root.render(
  <Router>
    <div>
        <Switch>
          {/* Páginas *//*}
          /*
          <Route exact path='/' component={Menu} />
        </Switch>
    </div>
  </Router>,
  document.getElementById('root')
);
*/
// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
