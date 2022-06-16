import logo from './logo.svg';

import './App.css';

import "primereact/resources/themes/lara-light-indigo/theme.css";  //theme
import "primereact/resources/primereact.min.css";                  //core css
import "primeicons/primeicons.css";
                            //icons

import Menu from './componentes/menu/Menu';
import Footer from './componentes/footer/Footer'
import Mapa from './componentes/mapa/Mapa'

function App() {
  function Welcome() {
    return <h1> HOLA, Wilmer</h1>;
  }
  return (
    <div className="App">
      <Menu />
     
      <Mapa />
      <Footer />
    </div>
        
  );
}

export default App;