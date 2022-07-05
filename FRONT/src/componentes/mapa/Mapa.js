import React, {useState, useEffect} from 'react'
import { DataTable } from 'primereact/datatable'
import { Column } from 'primereact/column';

class Mapa extends React.Component {
    /**
     * Constructor de la clase Mapa
     * @param {*} props 
     */
    constructor(props) {
        super(props)
        this.state = {
            mapa: [],
            isFetch: true,
        }
    }
    /**
     * Función componentDidMount que sirve para cargar el mapa
     */
    componentDidMount() {
        fetch("http://192.168.43.50:8082/mapa")
            .then(response => response.json())
            .then(mapaJson => this.setState( {
                mapa: mapaJson.data,
                isFetch: false
            }))
            .catch(error => console.log(error))
    }
    /**
     * Render de la clase Mapa
     * @returns Mapa con todas las atracciones
     */
    render () {
        
        const { mapa, isFetch } = this.state

        if (isFetch) {
            return <div>El mapa del parque no está disponible por el momento</div>
        }
        return (
          <div className ="container">
            
            <DataTable value={mapa}>
                <Column field="fila" header="Fila"></Column>
                <Column field="infoParque" header="InfoParque"></Column>
            </DataTable>

          </div>
            
        );
    }
}

export default Mapa;
