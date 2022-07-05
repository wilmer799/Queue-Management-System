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
        fetch("http://192.168.43.50:8082/visitantes")
            .then(response => response.json())
            .then(visitantesJson => this.setState( {
                visitantes: visitantesJson.data,
                isFetch: false
            }))
            .catch(error => console.log(error))
    }
    /**
     * Render que muestra la información de los visitantes 
     * @returns : Renderizado de los visitantes
     */
    render () {
        
        const { visitantes, isFetch } = this.state

        if (isFetch) {
            return <div>El estado de los visitantes no está disponible por el momento</div>
        }
        return (
          <div className ="container">
            <DataTable value={visitantes}>
                <Column field="id" header="ID"></Column>
                <Column field="nombre" header="Nombre"></Column>
                <Column field="posicionx" header="Posición x"></Column>
                <Column field="posiciony" header="Posición y"></Column>
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