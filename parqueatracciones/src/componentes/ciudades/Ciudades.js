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
            .then(response => response.json())
            .then(ciudadesJson => this.setState( {
                ciudades: ciudadesJson.data,
                isFetch: false
            }))
    }
    /**
     * Render que muestra la información de los visitantes 
     * @returns : Renderizado de los visitantes
     */
    render () {
        
        const { ciudades, isFetch } = this.state

        if (isFetch) {
            return 'Cargando...'
        }
        return (
          <div className ="container">
            <DataTable value={ciudades}>
                <Column field="nombre" header="nombre"></Column>
                <Column field="temperatura" header="temperatura"></Column>
            </DataTable>
          </div>
        );
    }
    
}
export default Ciudades;