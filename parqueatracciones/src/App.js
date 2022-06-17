import React  from 'react';
import './App.css';

import "primereact/resources/themes/lara-light-indigo/theme.css";  //theme
import "primereact/resources/primereact.min.css";                  //core css
import "primeicons/primeicons.css";
                            //icons

import Menu from './componentes/menu/Menu';
import Footer from './componentes/footer/Footer'
import Mapa from './componentes/mapa/Mapa'
import Visitante from './componentes/visitantes/Visitantes'
import Ciudad from './componentes/ciudades/Ciudades'

function App() {
  return (
    <div className="App">
      <Visitante />
      <Ciudad />
      <Mapa />
    </div>
        
  );
}

export default App;