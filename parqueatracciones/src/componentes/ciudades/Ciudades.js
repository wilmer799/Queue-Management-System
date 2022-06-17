import React from 'react'
import { DataTable } from 'primereact/datatable'
import { Column } from 'primereact/column';

class Ciudades extends React.Component {
    /**
     * Constructor de la clase Visitantes
     * @param {} props 
     */
    constructor(props) {
        super(props)
        this.state = {
            ciudades: [],
            isFetch: true,
        }
    }
    /**
     * ComponentDidMount que carga la información de los visitantes
     */
    componentDidMount() {
        fetch("http://localhost:8082/ciudades")
            .then(response => {
                if (!response.ok) throw Error(response.status);
                response.json();
            })
            .then(ciudadesJson => this.setState( {
                ciudades: ciudadesJson.data,
                isFetch: false
            }))
            .catch(error => console.log(error));
    }
    /**
     * Render que muestra la información de los visitantes 
     * @returns : Renderizado de los visitantes
     */
    render () {
        
        const { ciudades, isFetch } = this.state

        if (isFetch) {
            return (
                <div>Información de las ciudades no disponible</div>
            )
        }
        return (
          <div className ="container">
            <DataTable value={ciudades}>
                <Column field="cuadrante" header="Cuadrante"></Column>
                <Column field="nombre" header="Nombre"></Column>
                <Column field="temperatura" header="Temperatura"></Column>
            </DataTable>
          </div>
        );
    }
    
}
export default Ciudades;