import React, {useState, useEffect} from 'react'
import axios from 'axios'
import { DataTable } from 'primereact/datatable'
import { Column } from 'primereact/column';

class Mapa extends React.Component {
    
    constructor(props) {
        super(props)
        this.state = {
            mapa: [],
            isFetch: true,
        }
    }

    componentDidMount() {
        fetch("http://localhost:8082/mapa")
            .then(response => response.json())
            .then(mapaJson => this.setState( {
                mapa: mapaJson.infoParque,
                isFetch: false
            }))
    }

    render () {
        
        const { mapa, isFetch } = this.state

        if (isFetch) {
            return 'Cargando...'
        }
        return (
          <div className ="container">
            
            <DataTable value={mapa}>
                <Column field="fila" header="Fila"></Column>
                <Column field="infoParque" header="InfoParque"></Column>
            </DataTable>

          </div>
            
           /* <div className="container">

                <ol>
                    <li>this.mapa[0].infoParque</li>
                </ol>

            </div>*/
            
        );
    }
}

export default Mapa;

/*
function useDatos() {
        const [weather, setWeather] = useState([])
        useEffect(() => {
            
        return weather

        this.state.clima.map((climas) => (
                    <li> {climas} </li>)
    }
*/