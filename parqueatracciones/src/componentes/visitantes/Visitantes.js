import React from 'react'
import { DataTable } from 'primereact/datatable'
import { Column } from 'primereact/column';

class Visitantes extends React.Component {
    /**
     * Constructor de la clase Visitantes
     * @param {} props 
     */
    constructor(props) {
        super(props)
        this.state = {
            visitantes: [],
            isFetch: true,
        }
    }
    /**
     * ComponentDidMount que carga la información de los visitantes
     */
    componentDidMount() {
        fetch("http://localhost:8082/visitantes")
            .then(response => response.json())
            .then(visitantesJson => this.setState( {
                visitantes: visitantesJson.data,
                isFetch: false
            }))
    }
    /**
     * Render que muestra la información de los visitantes 
     * @returns : Renderizado de los visitantes
     */
    render () {
        
        const { visitantes, isFetch } = this.state

        if (isFetch) {
            return 'Cargando...'
        }
        return (
          <div className ="container">
            <DataTable value={visitantes}>
                <Column field="id" header="ID"></Column>
                <Column field="nombre" header="Nombre"></Column>
                <Column field="posicionx" header="Posicion x"></Column>
                <Column field="posiciony" header="Posicion y"></Column>
                <Column field="destinox" header="Destino x"></Column>
                <Column field="destinoy" header="Destino y"></Column>
                <Column field="idEnParque" header="ID en el parque"></Column>
                <Column field="ultimoEvento" header="Log"></Column>
            </DataTable>
          </div>
        );
    }
    
}
export default Visitantes;